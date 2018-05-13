package worker_test

import (
	"fmt"

	. "github.com/iwataka/mybot/worker"
)

func ExampleWorkerSignal_String() {
	fmt.Println(NewWorkerSignal(StartSignal))
	fmt.Println(NewWorkerSignal(RestartSignal))
	fmt.Println(NewWorkerSignal(StopSignal))
	fmt.Println(NewWorkerSignal(KillSignal))
	fmt.Println(NewWorkerSignal(PingSignal))
	// Output: Start
	// Restart
	// Stop
	// Kill
	// Ping
}
