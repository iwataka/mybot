package utils_test

import (
	"github.com/iwataka/mybot/utils"

	"fmt"
)

func ExampleCalcStringSlices() {
	ss1 := []string{"foo", "bar"}
	ss2 := []string{"foo", "other"}
	fmt.Println(len(utils.CalcStringSlices(ss1, ss2, true)))
	fmt.Println(len(utils.CalcStringSlices(ss1, ss2, false)))
	// Output: 3
	// 1
}

func ExampleCalcBools() {
	fmt.Println(utils.CalcBools(true, true, true))
	fmt.Println(utils.CalcBools(false, false, true))
	fmt.Println(utils.CalcBools(true, true, false))
	fmt.Println(utils.CalcBools(false, false, false))
	// Output: true
	// false
	// false
	// false
}

func ExampleCheckStringContained() {
	ss := []string{"foo", "bar"}
	fmt.Println(utils.CheckStringContained(ss, "foo"))
	fmt.Println(utils.CheckStringContained(ss, "other"))
	// Output: true
	// false
}
