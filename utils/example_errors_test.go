package utils

import (
	"github.com/BurntSushi/toml"

	"fmt"
)

func ExampleTomlUndecodedKeysError() {
	err := TomlUndecodedKeysError{[]toml.Key{[]string{"foo"}}, "foo.toml"}
	fmt.Println(err.Error())
	// Output: [foo] undecoded in foo.toml
}
