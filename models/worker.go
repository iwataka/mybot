package models

// Worker is worker which has its own operation and provides APIs to start/stop
// it.
type Worker interface {
	// Start starts this worker blockingly.
	// To stop this worker, close the given channel.
	Start(<-chan interface{}) error
	// Name returns a name of this Worker, to distinguish this from others.
	Name() string
}
