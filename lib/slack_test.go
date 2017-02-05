package mybot

import (
	"testing"
)

func TestSlackActionAdd(t *testing.T) {
	a1 := &SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a2 := &SlackAction{
		Channels: []string{"foo", "fizz"},
	}
	a1.Add(a2)
	if len(a1.Channels) != 3 {
		t.Fatalf("%v expected but %v found", 3, len(a1.Channels))
	}
}

func TestSlackActionSub(t *testing.T) {
	a1 := &SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a2 := &SlackAction{
		Channels: []string{"foo", "fizz"},
	}
	a1.Sub(a2)
	if len(a1.Channels) != 1 {
		t.Fatalf("%v expected but %v found", 1, len(a1.Channels))
	}
}
