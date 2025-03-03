package service

type Servicer interface {
	// start the service for the specified exchange
	Start() error
}
