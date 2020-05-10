package worker_test

import (
	"github.com/iwataka/mybot/worker"

	"context"
	"fmt"
)

type MyWorker struct {
	name string
}

func NewMyWorker(name string) *MyWorker {
	return &MyWorker{name}
}

func (w *MyWorker) Start(ctx context.Context, outChan chan<- interface{}) error {
	<-ctx.Done()
	return nil
}

func (w *MyWorker) Name() string {
	return w.name
}

func Example() {
	w := NewMyWorker("foo")
	wm := worker.NewWorkerManager(w, 0)
	defer wm.Close()

	// Start worker
	wm.Send(worker.StartSignal)
	fmt.Printf("Worker Status: %s\n", wm.Receive())

	// Stop worker
	wm.Send(worker.StopSignal)
	fmt.Printf("Worker Status: %s\n", wm.Receive())

	// Output: Worker Status: Started
	// Worker Status: Stopped
}
