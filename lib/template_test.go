package mybot

import (
	"testing"
)

func TestGetBoolSelectboxValue(t *testing.T) {
	val := make(map[string][]string)
	name := "foo"
	val[name] = []string{"true", ""}

	first := GetBoolSelectboxValue(val, 0, name)
	if first == nil || *first != true {
		t.Fatalf("%v expected but %v found", true, first)
	}

	second := GetBoolSelectboxValue(val, 1, name)
	if second != nil {
		t.Fatalf("%v expected but %v found", nil, second)
	}
}

func TestGetListTextboxValue(t *testing.T) {
	val := make(map[string][]string)
	name := "foo"
	val[name] = []string{"fizz, buzz ", ""}

	first := GetListTextboxValue(val, 0, name)
	if first[0] != "fizz" || first[1] != "buzz" || len(first) != 2 {
		t.Fatalf("%v expected but %v found", []string{"fizz", "buzz"}, first)
	}

	second := GetListTextboxValue(val, 1, name)
	if len(second) != 0 {
		t.Fatalf("%v expected but %v found", []string{}, second)
	}
}

func TestGetFloat64Ptr(t *testing.T) {
	val := make(map[string][]string)
	name := "foo"
	val[name] = []string{"1.23", ""}

	first, err := GetFloat64Ptr(val, 0, name)
	if err != nil {
		t.Fatal(err)
	}
	if first == nil || *first != 1.23 {
		t.Fatalf("%s expected but %v found", "*1.23", first)
	}

	second, err := GetFloat64Ptr(val, 1, name)
	if err != nil {
		t.Fatal(err)
	}
	if second != nil {
		t.Fatalf("%v expected but %v found", nil, second)
	}
}

func TestGetIntPtr(t *testing.T) {
	val := make(map[string][]string)
	name := "foo"
	val[name] = []string{"1", ""}

	first, err := GetIntPtr(val, 0, name)
	if err != nil {
		t.Fatal(err)
	}
	if first == nil || *first != 1 {
		t.Fatalf("%s expected but %v found", "*1", first)
	}

	second, err := GetIntPtr(val, 1, name)
	if err != nil {
		t.Fatal(err)
	}
	if second != nil {
		t.Fatalf("%v expected but %v found", nil, second)
	}
}
