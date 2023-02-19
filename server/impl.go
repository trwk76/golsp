package server

type Implementation interface {
	Initialized()
	Shutdown()
}
