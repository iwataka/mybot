package tmpl_test

import (
	"fmt"

	. "github.com/iwataka/mybot/tmpl"
)

func ExampleNewMap() {
	m := NewMap("key1", "val", "key2", 1)
	fmt.Println(m["key1"])
	fmt.Println(m["key2"])
	// Output: val
	// 1
}
