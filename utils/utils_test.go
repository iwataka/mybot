package utils_test

import (
	. "github.com/iwataka/mybot/utils"

	"testing"
)

func TestCalcStringSlices(t *testing.T) {
	s1 := []string{"foo", "bar"}
	s2 := []string{"foo", "else"}

	addResult := CalcStringSlices(s1, s2, true)
	if len(addResult) != 3 {
		t.Fatalf("%d expected but %d found", 3, len(addResult))
	}

	subResult := CalcStringSlices(s1, s2, false)
	if len(subResult) != 1 {
		t.Fatalf("%d expected but %d found", 1, len(subResult))
	}
}

func TestCalcBools(t *testing.T) {
	result1 := CalcBools(true, false, true)
	if !result1 {
		t.Fatalf("%v expected but %v found", true, result1)
	}

	result2 := CalcBools(true, true, false)
	if result2 {
		t.Fatalf("%v expected but %v found", false, result2)
	}
}

func TestCheckStringCotnained(t *testing.T) {
	ss := []string{"foo", "bar"}
	str := "foo"
	if !CheckStringContained(ss, str) {
		t.Fatalf("%v does not contain %s", ss, str)
	}

	str = "else"
	if CheckStringContained(ss, str) {
		t.Fatalf("%v contains %s", ss, str)
	}
}

func TestGenerateRandString(t *testing.T) {
	if str := GenerateRandString(0); len(str) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(str))
	}
	if str := GenerateRandString(1); len(str) != 1 {
		t.Fatalf("%d expected but %d found", 1, len(str))
	}
	if str := GenerateRandString(10); len(str) != 10 {
		t.Fatalf("%d expected but %d found", 10, len(str))
	}
}
