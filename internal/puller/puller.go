package puller

import (
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
	"github.com/ilia-tsyplenkov/klines-stat/internal/puller/bybit/rest"
	"github.com/ilia-tsyplenkov/klines-stat/internal/puller/bybit/ws"
)

type RestPuller interface {
	Pull() (*models.Kline, error)
}

type WsPuller interface {
	Start()
}

var _ RestPuller = (*rest.KLinePuller)(nil)
var _ WsPuller = (*ws.RecentTradePuller)(nil)
