package service

type Service struct {
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
