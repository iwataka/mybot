package worker

import (
	"fmt"
	"log"
	"time"
)

type MyWorker struct {
	name string
	ch   chan bool
}

func NewMyWorker(name string) *MyWorker {
	return &MyWorker{name, make(chan bool)}
}

func (w *MyWorker) Start() error {
	for range w.ch {
		return nil
	}
	return nil
}

func (w *MyWorker) Stop() error {
	w.ch <- false
	return nil
}

func (w *MyWorker) Name() string {
	return w.name
}

func Example() {
	w := NewMyWorker("foo")
	inChan, outChan := ActivateWorker(w, time.Minute)

	// ch is a channel to wait until the below goroutine processing
	// finishes (not used in actual codes)
	ch := make(chan bool)
	// Goroutine for capturing outputs
	go func() {
		for msg := range outChan {
			switch m := msg.(type) {
			case WorkerStatus:
				fmt.Printf("Worker %s\n", m)
			case error:
				log.Printf("%+v\n", m)
			}
			ch <- true
		}
	}()

	inChan <- NewWorkerSignal(StartSignal)
	<-ch
	inChan <- NewWorkerSignal(StopSignal)
	<-ch
	// Output: Worker Started
	// Worker Stopped
}
