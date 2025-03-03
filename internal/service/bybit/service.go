package bybit

import (
	st "github.com/ilia-tsyplenkov/klines-stat/internal/storage"
)

type Service struct {
	storage st.Storager
}

func New() *Service {
	return nil
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) startWSPullers() {}

func (s *Service) startRESTPullers() {}

func (s *Service) startKLineBuilders() {}

func (s *Service) startStorageWorkers() {}
