package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	builder "github.com/ilia-tsyplenkov/klines-stat/internal/builder/bybit"
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
	"github.com/ilia-tsyplenkov/klines-stat/internal/puller/bybit/rest"
	"github.com/ilia-tsyplenkov/klines-stat/internal/puller/bybit/ws"
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

	klineBuilderQueries := make(map[string]map[string]chan *models.RecentTrade)
	for pair := range cfg.Exchange["bybit"].Tickers {
		channels := make(map[string]chan *models.RecentTrade)
		for timeframe := range cfg.Exchange["bybit"].Timeframes {
			ch := make(chan *models.RecentTrade, 100)
			channels[timeframe] = ch
		}
		klineBuilderQueries[pair] = channels
	}

	bybitCfg := cfg.Exchange["bybit"]

	storage := pg.New(
		ctx,
		conn,
		bybitCfg,
		klineStorageCh,
		rtStorageCh,
	)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		storage.KLinesSaver()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		storage.RecentTradesSaver()
	}()

	rtPuller, err := ws.New(
		ctx,
		bybitCfg,
		klineBuilderQueries,
		rtStorageCh,
	)
	if err != nil {
		panic(err)
	}
	go rtPuller.Start()

	for pair := range bybitCfg.Tickers {
		for tf := range bybitCfg.Timeframes {
			puller := rest.New(
				ctx,
				klineStorageCh,
				cfg.Exchange["bybit"],
				pair,
				tf,
				bybitCfg.StartSince,
			)
			wg.Add(1)
			go func(puller *rest.KLinePuller) {
				defer wg.Done()
				kline, err := puller.Pull()
				if err != nil {
					log.Errorf("rest puller failed: pair: %s tf: %s: %v", pair, tf, err)
				}

				builder.New(ctx, bybitCfg, kline, klineBuilderQueries[pair][tf], klineStorageCh).Start()
			}(puller)

		}
	}
	wg.Wait()
	// fmt.Println("all done")
	// cancel()
	// time.Sleep(1 * time.Second)
}
