package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	bbService "github.com/ilia-tsyplenkov/klines-stat/internal/service/bybit"
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

	// graceful shutdown
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
		log.Fatal(err)
	}

	srv := bbService.New(
		ctx,
		"bybit",
		cfg.Exchange["bybit"],
		conn,
	)

	if err := srv.Start(); err != nil {
		log.Error(err)
	}
}
