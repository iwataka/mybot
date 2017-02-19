package mybot

import (
	"testing"
)

func TestStringsOp(t *testing.T) {
	s1 := []string{"foo", "bar"}
	s2 := []string{"foo", "else"}

	addResult := StringsOp(s1, s2, true)
	if len(addResult) != 3 {
		t.Fatalf("%d expected but %d found", 3, len(addResult))
	}

	subResult := StringsOp(s1, s2, false)
	if len(subResult) != 1 {
		t.Fatalf("%d expected but %d found", 1, len(subResult))
	}
}

func TestBoolOP(t *testing.T) {
	result1 := BoolOp(true, false, true)
	if !result1 {
		t.Fatalf("%v expected but %v found", true, result1)
	}

	result2 := BoolOp(true, true, false)
	if result2 {
		t.Fatalf("%v expected but %v found", false, result2)
	}
}

func TestStringsContains(t *testing.T) {
	ss := []string{"foo", "bar"}
	str := "foo"
	if !StringsContains(ss, str) {
		t.Fatalf("%v does not contain %s", ss, str)
	}

	str = "else"
	if StringsContains(ss, str) {
		t.Fatalf("%v contains %s", ss, str)
	}
}

func TestRandString(t *testing.T) {
	if str := RandString(0); len(str) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(str))
	}
	if str := RandString(1); len(str) != 1 {
		t.Fatalf("%d expected but %d found", 1, len(str))
	}
	if str := RandString(10); len(str) != 10 {
		t.Fatalf("%d expected but %d found", 10, len(str))
	}
}
