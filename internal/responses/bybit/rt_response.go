package bybit

import (
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
)

type RTResponse struct {
	Topic string                `json:"topic"`
	Ts    int64                 `json:"ts"`
	Type  string                `json:"type"`
	Data  []*models.RecentTrade `json:"data"`
}
