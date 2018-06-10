package worker_test

import (
	"fmt"

	"github.com/iwataka/mybot/worker"
)

func ExampleWorkerSignal_String() {
	fmt.Println(worker.NewWorkerSignal(worker.StartSignal))
	fmt.Println(worker.NewWorkerSignal(worker.RestartSignal))
	fmt.Println(worker.NewWorkerSignal(worker.StopSignal))
	fmt.Println(worker.NewWorkerSignal(worker.KillSignal))
	fmt.Println(worker.NewWorkerSignal(worker.PingSignal))
	fmt.Println(worker.NewWorkerSignal(-1))
	// Output: Start
	// Restart
	// Stop
	// Kill
	// Ping
	//
}

func ExampleWorkerStatus_String() {
	fmt.Println(worker.StatusActive)
	fmt.Println(worker.StatusFinished)
	fmt.Println(worker.StatusInactive)
	fmt.Println(worker.StatusKilled)
	fmt.Println(worker.StatusStarted)
	fmt.Println(worker.StatusStopped)
	fmt.Println(worker.WorkerStatus(-1))
	// Output: Active
	// Finished
	// Inactive
	// Killed
	// Started
	// Stopped
	//
}
