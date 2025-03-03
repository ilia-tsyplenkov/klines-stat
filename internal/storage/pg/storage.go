package pg

import (
	"context"
	"time"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
	queries "github.com/ilia-tsyplenkov/klines-stat/internal/storage/queries/pg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Storage struct {
	ctx        context.Context
	cfg        config.Exchange
	conn       *pgxpool.Pool
	klineQuery chan *models.Kline
	rtQuery    chan *models.RecentTrade
}

func New(
	ctx context.Context,
	conn *pgxpool.Pool,
	cfg config.Exchange,
	klineCh chan *models.Kline,
	rtCh chan *models.RecentTrade,
) *Storage {
	return &Storage{
		ctx:        ctx,
		conn:       conn,
		cfg:        cfg,
		klineQuery: klineCh,
		rtQuery:    rtCh,
	}
}

func (s *Storage) KLinesSaver() {

	batch := &pgx.Batch{}
	l := log.WithField("action", "kline saver")

	for tick := time.Tick(100 * time.Millisecond); ; {
		select {
		case <-s.ctx.Done():
			return
		case <-tick:
			// write the batch
			if batch.Len() == 0 {
				continue
			}

			br := s.conn.SendBatch(s.ctx, batch)
			_, err := br.Exec()
			if err != nil {
				l.Errorf("writing batch: %v", err)
			} else {
				l.Infof("saving klines: %d", batch.Len())
			}

			br.Close()

			batch = &pgx.Batch{}
		case kl := <-s.klineQuery:
			// put to the batch
			batch.Queue(
				queries.InsertKLineQuery,
				s.cfg.Tickers[kl.Pair],
				kl.TimeFrame,
				kl.O,
				kl.H,
				kl.L,
				kl.C,
				kl.UtcBegin,
				kl.UtcEnd,
			)
		}
	}
}

func (s *Storage) RecentTradesSaver() {

	l := log.WithField("action", "rt saver")

	batch := &pgx.Batch{}

	for tick := time.Tick(100 * time.Millisecond); ; {
		select {
		case <-s.ctx.Done():
			return
		case <-tick:
			// write the batch
			if batch.Len() == 0 {
				continue
			}
			br := s.conn.SendBatch(s.ctx, batch)
			_, err := br.Exec()
			if err != nil {
				l.Errorf("writing batch: %v", err)
			}
			br.Close()
			batch = &pgx.Batch{}

		default:
		fillBatchLoop:
			for {
				select {
				default:
					break fillBatchLoop

				case rt := <-s.rtQuery:
					// put to the batch
					batch.Queue(
						queries.InsertRecentTradeQuery,
						rt.Tid,
						s.cfg.Tickers[rt.Pair],
						rt.Price,
						rt.Amount,
						rt.Side,
						rt.Timestamp,
					)
				}
			}

		}
	}
}
