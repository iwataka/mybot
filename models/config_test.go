package models_test

import (
	"testing"

	"github.com/iwataka/mybot/models"
	"github.com/stretchr/testify/require"
)

func Test_NewTimelineProperties(t *testing.T) {
	props := models.NewTimelineProperties()
	require.True(t, props.ExcludeReplies)
	require.False(t, props.IncludeRts)
}

func Test_NewSearchProperties(t *testing.T) {
	props := models.NewSearchProperties()
	require.Equal(t, "mixed", props.ResultType)
}
