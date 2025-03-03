package bybit

import (
	"context"
	"strconv"
	"time"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
	log "github.com/sirupsen/logrus"
)

const gap float64 = 0.0000000000001

type KlineBuilder struct {
	ctx       context.Context
	cfg       config.Exchange
	kline     *models.Kline
	timeframe int64
	rtCh      chan *models.RecentTrade
	storageCh chan *models.Kline
}

func New(
	ctx context.Context,
	cfg config.Exchange,
	kline *models.Kline,
	// timefrage int64,
	rtCh chan *models.RecentTrade,
	storageCh chan *models.Kline,
) *KlineBuilder {
	return &KlineBuilder{
		ctx:       ctx,
		cfg:       cfg,
		kline:     kline,
		timeframe: cfg.Timeframes[kline.TimeFrame],
		rtCh:      rtCh,
		storageCh: storageCh,
	}
}

func (b *KlineBuilder) Start() {

	l := log.WithFields(log.Fields{
		"action":    "builder",
		"pair":      b.kline.Pair,
		"timeframe": b.kline.TimeFrame,
	})

	l.Infof("start with kline: %+v", *b.kline)

	nextTick := b.kline.UtcBegin + b.timeframe - time.Now().UTC().Unix()*1000
	ticker := time.NewTicker(time.Duration(nextTick) * time.Millisecond)

	defer ticker.Stop()
	for {
		select {
		case <-b.ctx.Done():
			return
		case <-ticker.C:
			l.Info("tick")
			b.storageCh <- b.kline

			b.kline = &models.Kline{
				Pair:      b.kline.Pair,
				TimeFrame: b.kline.TimeFrame,
				UtcBegin:  b.kline.UtcEnd,
				UtcEnd:    b.kline.UtcEnd + b.timeframe,
			}

			ticker.Reset(time.Duration(b.timeframe) * time.Millisecond)
		default:
		fillKlineLoop:
			for {
				select {
				default:
					break fillKlineLoop
				case rt := <-b.rtCh:
					if b.kline.UtcBegin > rt.Timestamp || b.kline.UtcEnd < rt.Timestamp {
						continue
					}
					price, err := strconv.ParseFloat(rt.Price, 64)

					if err != nil {
						l.Error(err)
					}
					b.kline.C = price

					if b.kline.O == 0.0 {
						b.kline.O = price
						b.kline.L = price
						b.kline.H = price
					}
					if b.kline.L > price {
						b.kline.L = price
					}
					if b.kline.H < price {
						b.kline.H = price
					}

				}
			}
		}
	}
}
