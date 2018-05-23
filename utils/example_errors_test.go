package utils_test

import (
	"github.com/BurntSushi/toml"
	. "github.com/iwataka/mybot/utils"

	"fmt"
)

func ExampleStreamInterruptedError() {
	err := NewStreamInterruptedError()
	fmt.Println(err.Error())
	// Output: Interrupted
}

func ExampleTomlUndecodedKeysError() {
	err := TomlUndecodedKeysError{[]toml.Key{[]string{"foo"}}, "foo.toml"}
	fmt.Println(err.Error())
	// Output: [foo] undecoded in foo.toml
}
