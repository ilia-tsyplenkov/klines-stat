package bybit

import (
	"context"
	"strconv"
	"time"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
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

	tick := time.After(time.Duration(b.kline.UtcBegin+b.timeframe-time.Now().UTC().Unix()) * time.Millisecond)
	for {
		select {
		case <-b.ctx.Done():
			return
		case <-tick:
			b.storageCh <- b.kline
			b.kline = &models.Kline{
				Pair:      b.kline.Pair,
				TimeFrame: b.kline.TimeFrame,
				UtcBegin:  b.kline.UtcEnd,
				UtcEnd:    b.kline.UtcEnd + b.timeframe,
			}
			tick = time.After(time.Duration(b.timeframe) * time.Millisecond)
		case rt := <-b.rtCh:
			if b.kline.UtcBegin < rt.Timestamp || b.kline.UtcEnd <= rt.Timestamp {
				continue
			}
			price, _ := strconv.ParseFloat(rt.Price, 64)
			b.kline.C = price

			if b.kline.O == 0.0 {
				b.kline.O = price
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
