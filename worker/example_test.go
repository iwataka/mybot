package worker

import (
	"context"
	"fmt"
)

type MyWorker struct {
	name string
}

func NewMyWorker(name string) *MyWorker {
	return &MyWorker{name}
}

func (w *MyWorker) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (w *MyWorker) Name() string {
	return w.name
}

func Example() {
	w := NewMyWorker("foo")
	wm := NewWorkerManager(w, 0)
	defer wm.Close()

	// Start worker
	wm.Send(StartSignal)
	fmt.Printf("Worker Status: %s\n", wm.Receive())

	// Stop worker
	wm.Send(StopSignal)
	fmt.Printf("Worker Status: %s\n", wm.Receive())

	// Output: Worker Status: Started
	// Worker Status: Stopped
}
