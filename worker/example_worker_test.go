package worker_test

import (
	"github.com/iwataka/mybot/worker"

	"fmt"
)

func ExampleWorkerSignal_String() {
	fmt.Println(worker.StartSignal)
	fmt.Println(worker.RestartSignal)
	fmt.Println(worker.StopSignal)
	fmt.Println(worker.WorkerSignal(-1))
	// Output: Start
	// Restart
	// Stop
	//
}

func ExampleWorkerStatus_String() {
	fmt.Println(worker.StatusFinished)
	fmt.Println(worker.StatusStarted)
	fmt.Println(worker.StatusStopped)
	fmt.Println(worker.WorkerStatus(-1))
	// Output: Finished
	// Started
	// Stopped
	//
}
