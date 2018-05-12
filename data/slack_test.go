package data_test

import (
	. "github.com/iwataka/mybot/data"

	"testing"
)

func TestSlackActionAdd(t *testing.T) {
	a1 := SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a1.Pin = true
	a1.Reactions = []string{":smile:"}
	a2 := SlackAction{
		Channels: []string{"foo", "fizz"},
	}
	a2.Reactions = []string{":smile:", ":plane:"}
	result := a1.Add(a2)

	if !result.Pin {
		t.Fatalf("%v expected but %v found", true, result.Pin)
	}
	if result.Star {
		t.Fatalf("%v expected but %v found", false, result.Star)
	}
	if len(result.Reactions) != 2 {
		t.Fatalf("%v expected but %v found", 2, len(result.Reactions))
	}
	if len(result.Channels) != 3 {
		t.Fatalf("%v expected but %v found", 3, len(result.Channels))
	}
}

func TestSlackActionSub(t *testing.T) {
	a1 := SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a2 := SlackAction{
		Channels: []string{"foo", "fizz"},
	}
	result := a1.Sub(a2)
	if len(result.Channels) != 1 {
		t.Fatalf("%v expected but %v found", 1, len(result.Channels))
	}
}

func TestSlackActionIsEmpty(t *testing.T) {
	a := NewSlackAction()
	if !a.IsEmpty() {
		t.Fatalf("%v should be empty", a)
	}
	a.Channels = []string{"foo"}
	if a.IsEmpty() {
		t.Fatalf("%v should not be empty", a)
	}
}
