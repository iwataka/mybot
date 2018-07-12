package models

type WorkerMessageHandler interface {
	Handle(msg interface{}) error
}

// Worker is worker which has its own operation and provides APIs to start/stop
// it.
type Worker interface {
	Start() error
	Stop() error
	// Name returns a name of this Worker, to distinguish this from others.
	Name() string
}
