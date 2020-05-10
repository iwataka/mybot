package data_test

import (
	"github.com/iwataka/mybot/data"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestSlackAction_Add(t *testing.T) {
	a1 := data.SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a1.Pin = true
	a1.Reactions = []string{":smile:"}
	a2 := data.SlackAction{
		Channels: []string{"foo", "fizz"},
	}
	a2.Reactions = []string{":smile:", ":plane:"}

	result := a1.Add(a2)

	require.True(t, result.Pin)
	require.False(t, result.Star)
	require.Len(t, result.Reactions, 2)
	require.Len(t, result.Channels, 3)
}

func TestSlackAction_Sub(t *testing.T) {
	a1 := data.SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a2 := data.SlackAction{
		Channels: []string{"foo", "fizz"},
	}

	result := a1.Sub(a2)

	require.Len(t, result.Channels, 1)
}

func TestSlackAction_IsEmpty(t *testing.T) {
	a := data.NewSlackAction()
	require.True(t, a.IsEmpty())

	a.Channels = []string{"foo"}
	require.False(t, a.IsEmpty())
}
