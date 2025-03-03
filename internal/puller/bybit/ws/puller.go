package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilia-tsyplenkov/klines-stat/config"
	"github.com/ilia-tsyplenkov/klines-stat/internal/models"
	response "github.com/ilia-tsyplenkov/klines-stat/internal/responses/bybit"
	log "github.com/sirupsen/logrus"
	bybitWS "github.com/wuhewuhe/bybit.go.api"
)

type RecentTradePuller struct {
	ctx          context.Context
	exchangeCfg  config.Exchange
	requestURL   string
	ws           *bybitWS.WebSocket
	pairQueries  map[string]map[string]chan *models.RecentTrade
	msgQuery     chan string
	storageQuery chan *models.RecentTrade
}

func New(
	ctx context.Context,
	exchangeCfg config.Exchange,
	pairQueries map[string]map[string]chan *models.RecentTrade,
	storageQuery chan *models.RecentTrade,

) (*RecentTradePuller, error) {
	puller := &RecentTradePuller{
		ctx:          ctx,
		exchangeCfg:  exchangeCfg,
		pairQueries:  pairQueries,
		msgQuery:     make(chan string, 1_000),
		storageQuery: storageQuery,
	}
	ws := bybitWS.NewBybitPublicWebSocket(exchangeCfg.WSApiUrl, func(message string) error {
		puller.msgQuery <- message
		return nil
	})
	puller.ws = ws

	return puller, nil
}

// Creates web socket connection,
// subscribes to all topics in the config and
// launches internal handler of ws data
func (p *RecentTradePuller) Start() {
	go p.handler()
	p.ws.Connect()
	defer p.ws.Disconnect()
	l := log.WithField("action", "ws puller start")

	for pair := range p.exchangeCfg.Tickers {
		if _, err := p.ws.SendSubscription([]string{fmt.Sprintf("publicTrade.%s", pair)}); err != nil {
			l.Errorf("failed to subscribe to %q", pair)
		}
	}

	<-p.ctx.Done()
}

// Receives and pareses messages,
// sends each message to the kline builders of this pair
// E.g BTCUDSDT msg goes to 1m, 15m, 60m, 1d kline builders
// plus sends it to the saver worker
func (p *RecentTradePuller) handler() {

	for {
		select {
		case <-p.ctx.Done():
			return
		case msg := <-p.msgQuery:
			rt := &response.RTResponse{}
			err := json.Unmarshal([]byte(msg), rt)
			if err != nil {
				log.Errorf("rt msg handler: msg: %q: %+v", msg, err)
			}

			for _, rt := range rt.Data {
				// send to each candle builder
				for _, ch := range p.pairQueries[rt.Pair] {
					select {
					case ch <- rt:
						{
						}
					default:
						{
						}
					}
				}

				// send to storage worker
				select {
				case p.storageQuery <- rt:
					{
					}
				default:
					{
					}

				}

			}
		}
	}

}
