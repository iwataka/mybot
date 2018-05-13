package worker_test

import (
	"fmt"
	"time"

	"github.com/iwataka/mybot/worker"
)

type MyWorker struct {
	name string
	ch   chan bool
}

func (w *MyWorker) Start() error {
	fmt.Println("Started")
	for {
		select {
		case msg := <-w.ch:
			if msg {
				fmt.Println("Received")
			} else {
				return nil
			}
		}
	}
}

func (w *MyWorker) Stop() {
	fmt.Println("Stopped")
	w.ch <- false
}

func (w *MyWorker) Name() string {
	return w.name
}

func Example() {
	ch := make(chan bool)
	w := &MyWorker{"foo", ch}
	inChan, outChan := worker.ActivateWorker(w, time.Minute)
	inChan <- worker.NewWorkerSignal(worker.StartSignal)
	<-outChan
	ch <- true
	inChan <- worker.NewWorkerSignal(worker.StopSignal)
	<-outChan
	// Output: Started
	// Received
	// Stopped
}
