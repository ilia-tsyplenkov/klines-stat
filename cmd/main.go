package main

import (
	"github.com/ilia-tsyplenkov/klines-stat/config"
	log "github.com/sirupsen/logrus"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	log.Infof("config: %+v", cfg)
}
