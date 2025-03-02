package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
	"github.com/ilia-tsyplenkov/klines-stat/internal/responses/bybit"
	log "github.com/sirupsen/logrus"
)

type KLinePuller struct {
	ctx          context.Context
	storageQueue chan *models.Kline
	exchageCfg   config.Exchange
	timeframe    string
	requestURL   string
	startTs      int64
}

func New(
	ctx context.Context,
	storageQueue chan *models.Kline,
	exchageCfg config.Exchange,
	pair string,
	timeframe string,
	startTs int64,
) *KLinePuller {

	bbPair := exchageCfg.Tickers[pair]
	requestUrl := fmt.Sprintf("%s?category=%s&symbol=%s&interval=%s", exchageCfg.RestApiURL, exchageCfg.Category, bbPair, timeframe)
	log.Infof("kline puller new[%s](%s): timeframe: %s: requestUrl: %q", bbPair, pair, timeframe, requestUrl)

	return &KLinePuller{
		ctx:          ctx,
		storageQueue: storageQueue,
		exchageCfg:   exchageCfg,
		timeframe:    timeframe,
		requestURL:   requestUrl,
		startTs:      startTs,
	}
}

func (p *KLinePuller) Pull() (*models.Kline, error) {
	var kline *models.Kline
	for i := 1; ; i++ {
		klineURL := fmt.Sprintf("%s&start=%d", p.requestURL, p.startTs)
		// log.Infof("klineURL: %s", klineURL)
		klinesResp, err := p.getKLines(klineURL)
		if err != nil {
			panic(err)
		}
		log.Infof("%d:%s[%d](%s)\n", i, klinesResp.Result.Symbol, len(klinesResp.Result.List), p.timeframe)

		for i := len(klinesResp.Result.List) - 1; i >= 0; i-- {
			kl := klinesResp.Result.List[i]
			openPrice, _ := strconv.ParseFloat(kl[1], 64)
			highPrice, _ := strconv.ParseFloat(kl[2], 64)
			lowPrice, _ := strconv.ParseFloat(kl[3], 64)
			closePrice, _ := strconv.ParseFloat(kl[4], 64)
			utcBegin, _ := strconv.ParseInt(kl[0], 10, 64)
			kline = &models.Kline{
				Pair:      klinesResp.Result.Symbol,
				TimeFrame: p.timeframe,
				O:         openPrice,
				H:         highPrice,
				L:         lowPrice,
				C:         closePrice,
				UtcBegin:  utcBegin,
				UtcEnd:    utcBegin + p.exchageCfg.Timeframes[p.timeframe],
			}

			p.storageQueue <- kline
		}
		if !klinesResp.IsLast {
			p.startTs = klinesResp.NextTS
		} else {
			// log.Info("is last record")
			break
		}
	}
	return kline, nil
}

func (p *KLinePuller) getKLines(url string) (*bybit.KLineResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	klines := &bybit.KLineResponse{}
	if err = json.NewDecoder(resp.Body).Decode(klines); err != nil {
		return nil, err
	}

	ts, _ := strconv.Atoi(klines.Result.List[0][0])

	// log.Infof("ts: %d tf: %d klines.Time: %d\n", ts, p.exchageCfg.Timeframes[p.timeframe], klines.Time)
	klines.IsLast = int64(ts)+p.exchageCfg.Timeframes[p.timeframe] > klines.Time

	klines.NextTS = int64(ts) + p.exchageCfg.Timeframes[p.timeframe]

	return klines, nil
}

func (p *KLinePuller) StartKLineBuilder(last *models.Kline) error {

	return nil

}
