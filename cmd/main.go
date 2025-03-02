package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
	"github.com/ilia-tsyplenkov/klines-stat/internal/puller/bybit/rest"
	"github.com/ilia-tsyplenkov/klines-stat/internal/storage/pg"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	log.Infof("config: %+v", cfg)

	pgxCfg, err := pgxpool.ParseConfig(cfg.DB)
	if err != nil {
		log.Fatalf("parse connection string: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	signals := make(chan os.Signal, 1)
	go func() {
		signal.Notify(signals,
			os.Interrupt,
			syscall.SIGTERM,
		)
		sig := <-signals
		log.Infof("got signal: %v", sig)
		log.Info("shutting down....")
		cancel()
	}()

	conn, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	defer conn.Close()
	if err != nil {
		log.Fatalf("failed connect to db: %v", err)
	}

	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}

	klineStorageCh := make(chan *models.Kline, 100)
	rtStorageCh := make(chan *models.RecentTrade, 100)

	storage := pg.New(
		ctx,
		conn,
		klineStorageCh,
		rtStorageCh,
	)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		storage.KLinesSaver()
	}()

	for pair := range cfg.Exchange["bybit"].Tickers {
		for tf := range cfg.Exchange["bybit"].Timeframes {
			puller := rest.New(
				ctx,
				klineStorageCh,
				cfg.Exchange["bybit"],
				pair,
				tf,
				1738368000000,
			)
			wg.Add(1)
			go func(puller *rest.KLinePuller) {
				defer wg.Done()
				puller.Pull()
			}(puller)

		}
	}
	wg.Wait()
	// fmt.Println("all done")
	// cancel()
	// time.Sleep(1 * time.Second)
}
