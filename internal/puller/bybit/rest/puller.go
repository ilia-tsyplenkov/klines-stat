package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
	"github.com/ilia-tsyplenkov/klines-stat/internal/responses/bybit"
)

type KLinePuller struct {
	ctx          context.Context
	storageQueue chan *models.Kline
	exchageCfg   config.Exchange
	timeframe    string
	requestURL   string
}

func New() *KLinePuller {
	return nil
}

func (p *KLinePuller) Pull() error {
	return nil
}

func (p *KLinePuller) getKLines(url string) (*bybit.KLineResponse, error) {
	resp, err := http.Get(p.requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	klines := &bybit.KLineResponse{}
	if err = json.NewDecoder(resp.Body).Decode(klines); err != nil {
		return nil, err
	}

	ts, _ := strconv.Atoi(klines.Result.List[0][0])

	klines.IsLast = int64(ts)+p.exchageCfg.Timeframes[p.timeframe] > klines.Time
	klines.NextTS = int64(ts) + p.exchageCfg.Timeframes[p.timeframe]

	return klines, nil
}

func (p *KLinePuller) StartKLineBuilder(last *models.Kline) error {

	return nil

}
