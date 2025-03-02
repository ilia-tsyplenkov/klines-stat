package storage

import (
	"github.com/ilia-tsyplenkov/klines-stat/internal/storage/pg"
)

type Storager interface {
	KLinesSaver()
	RecentTradesSaver()
}

var _ Storager = (*pg.Storage)(nil)
