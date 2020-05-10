package data_test

import (
	"github.com/iwataka/mybot/data"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestTwitterAction(t *testing.T) {
	a1 := data.TwitterAction{
		Collections: []string{"col1", "col2"},
	}
	a1.Retweet = true
	a2 := data.TwitterAction{
		Collections: []string{"col1", "col3"},
	}
	a2.Retweet = true
	a2.Favorite = true

	result1 := a1.Add(a2)

	require.True(t, result1.Retweet)
	require.True(t, result1.Favorite)
	require.Len(t, result1.Collections, 3)

	result2 := a1.Sub(a2)
	require.False(t, result2.Retweet)
	require.False(t, result2.Favorite)
	require.Len(t, result2.Collections, 1)
}
