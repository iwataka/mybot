package worker

import (
	"fmt"
)

func ExampleWorkerSignal_String() {
	fmt.Println(StartSignal)
	fmt.Println(RestartSignal)
	fmt.Println(StopSignal)
	fmt.Println(WorkerSignal(-1))
	// Output: Start
	// Restart
	// Stop
	//
}

func ExampleWorkerStatus_String() {
	fmt.Println(StatusFinished)
	fmt.Println(StatusStarted)
	fmt.Println(StatusStopped)
	fmt.Println(WorkerStatus(-1))
	// Output: Finished
	// Started
	// Stopped
	//
}
