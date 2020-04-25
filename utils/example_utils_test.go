package utils

import (
	"fmt"
)

func ExampleCalcStringSlices() {
	ss1 := []string{"foo", "bar"}
	ss2 := []string{"foo", "other"}
	fmt.Println(len(CalcStringSlices(ss1, ss2, true)))
	fmt.Println(len(CalcStringSlices(ss1, ss2, false)))
	// Output: 3
	// 1
}

func ExampleCalcBools() {
	fmt.Println(CalcBools(true, true, true))
	fmt.Println(CalcBools(false, false, true))
	fmt.Println(CalcBools(true, true, false))
	fmt.Println(CalcBools(false, false, false))
	// Output: true
	// false
	// false
	// false
}

func ExampleCheckStringContained() {
	ss := []string{"foo", "bar"}
	fmt.Println(CheckStringContained(ss, "foo"))
	fmt.Println(CheckStringContained(ss, "other"))
	// Output: true
	// false
}
