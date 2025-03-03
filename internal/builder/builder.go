package builder

import (
	"github.com/ilia-tsyplenkov/klines-stat/internal/builder/bybit"
)

type Builder interface {
	// start kline builder
	Start()
}

var _ Builder = (*bybit.KlineBuilder)(nil)
