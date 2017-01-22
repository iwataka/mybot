package mybot

import (
	"testing"
)

func TestGetIntPtr(t *testing.T) {
	val := make(map[string][]string)
	val["foo"] = []string{"1"}
	index := 0
	name := "foo"
	result, err := GetIntPtr(val, index, name)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || *result != 1 {
		t.Fatalf("%s expected but %v found", "*1", result)
	}
}
