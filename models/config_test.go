package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewTimelineProperties(t *testing.T) {
	props := NewTimelineProperties()
	require.True(t, props.ExcludeReplies)
	require.False(t, props.IncludeRts)
}

func Test_NewSearchProperties(t *testing.T) {
	props := NewSearchProperties()
	require.Equal(t, "mixed", props.ResultType)
}
