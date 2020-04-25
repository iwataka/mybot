package worker

import (
	"fmt"
)

func ExampleWorkerSignal_String() {
	fmt.Println(NewWorkerSignal(StartSignal))
	fmt.Println(NewWorkerSignal(RestartSignal))
	fmt.Println(NewWorkerSignal(StopSignal))
	fmt.Println(NewWorkerSignal(KillSignal))
	fmt.Println(NewWorkerSignal(PingSignal))
	fmt.Println(NewWorkerSignal(-1))
	// Output: Start
	// Restart
	// Stop
	// Kill
	// Ping
	//
}

func ExampleWorkerStatus_String() {
	fmt.Println(StatusActive)
	fmt.Println(StatusFinished)
	fmt.Println(StatusInactive)
	fmt.Println(StatusKilled)
	fmt.Println(StatusStarted)
	fmt.Println(StatusStopped)
	fmt.Println(WorkerStatus(-1))
	// Output: Active
	// Finished
	// Inactive
	// Killed
	// Started
	// Stopped
	//
}
