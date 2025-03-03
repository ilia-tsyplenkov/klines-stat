package bybit

import (
	"context"
	"sync"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	builder "github.com/ilia-tsyplenkov/klines-stat/internal/builder/bybit"
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
	"github.com/ilia-tsyplenkov/klines-stat/internal/puller/bybit/rest"
	"github.com/ilia-tsyplenkov/klines-stat/internal/puller/bybit/ws"
	"github.com/ilia-tsyplenkov/klines-stat/internal/storage/pg"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	ctx          context.Context
	cfg          config.Exchange
	conn         *pgxpool.Pool
	exchangeName string

	wg sync.WaitGroup
}

func New(
	ctx context.Context,
	exchangeName string,
	cfg config.Exchange,
	conn *pgxpool.Pool,
) *Service {
	return &Service{
		ctx:          ctx,
		cfg:          cfg,
		conn:         conn,
		exchangeName: exchangeName,
	}
}

func (s *Service) Start() error {

	klineStorageCh := make(chan *models.Kline, 100)
	rtStorageCh := make(chan *models.RecentTrade, 100)

	klineBuilderQueries := make(map[string]map[string]chan *models.RecentTrade)
	for pair := range s.cfg.Tickers {
		channels := make(map[string]chan *models.RecentTrade)
		for timeframe := range s.cfg.Timeframes {
			ch := make(chan *models.RecentTrade, 100)
			channels[timeframe] = ch
		}
		klineBuilderQueries[pair] = channels
	}

	storage := pg.New(
		s.ctx,
		s.conn,
		s.cfg,
		klineStorageCh,
		rtStorageCh,
	)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		storage.KLinesSaver()
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		storage.RecentTradesSaver()
	}()

	rtPuller, err := ws.New(
		s.ctx,
		s.cfg,
		klineBuilderQueries,
		rtStorageCh,
	)
	if err != nil {
		panic(err)
	}
	go rtPuller.Start()

	for pair := range s.cfg.Tickers {
		for tf := range s.cfg.Timeframes {
			puller := rest.New(
				s.ctx,
				klineStorageCh,
				s.cfg,
				pair,
				tf,
				s.cfg.StartSince,
			)
			s.wg.Add(1)
			go func(puller *rest.KLinePuller) {
				defer s.wg.Done()
				kline, err := puller.Pull()
				if err != nil {
					log.Errorf("rest puller failed: pair: %s tf: %s: %v", pair, tf, err)
				}

				builder.New(s.ctx, s.cfg, kline, klineBuilderQueries[pair][tf], klineStorageCh).Start()
			}(puller)

		}
	}
	s.wg.Wait()
	return nil
}
