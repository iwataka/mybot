package worker_test

import (
	"fmt"
	"log"
	"time"

	"github.com/iwataka/mybot/worker"
)

type MyWorker struct {
	name string
	ch   chan bool
}

func NewMyWorker(name string) *MyWorker {
	return &MyWorker{name, make(chan bool)}
}

func (w *MyWorker) Start() error {
	for {
		select {
		case <-w.ch:
			return nil
		}
	}
}

func (w *MyWorker) Stop() {
	w.ch <- false
}

func (w *MyWorker) Name() string {
	return w.name
}

func Example() {
	w := NewMyWorker("foo")
	inChan, outChan := worker.ActivateWorker(w, time.Minute)

	// ch is a channel to wait until the below goroutine processing
	// finishes (not used in actual codes)
	ch := make(chan bool)
	// Goroutine for capturing outputs
	go func() {
		for msg := range outChan {
			switch m := msg.(type) {
			case worker.WorkerStatus:
				fmt.Printf("Worker %s\n", m)
			case error:
				log.Printf("%+v\n", m)
			}
			ch <- true
		}
	}()

	inChan <- worker.NewWorkerSignal(worker.StartSignal)
	<-ch
	inChan <- worker.NewWorkerSignal(worker.StopSignal)
	<-ch
	// Output: Worker Started
	// Worker Stopped
}
