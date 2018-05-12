package utils_test

import (
	. "github.com/iwataka/mybot/utils"

	"fmt"
)

func ExampleStreamInterruptedError() {
	err := NewStreamInterruptedError()
	fmt.Println(err.Error())
	// Output: Interrupted
}
