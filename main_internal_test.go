package main

import (
	"os"
	"path/filepath"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/mocks"
	"github.com/stretchr/testify/require"
)

func Test_argValueWithMkdir(t *testing.T) {
	parent := "parent"
	child := "child"
	path := filepath.Join(parent, child)
	key := "key"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context := mocks.NewMockContext(ctrl)
	context.EXPECT().String(key).Return(path)

	var info os.FileInfo
	info, _ = os.Stat(parent)
	require.Nil(t, info)

	val, err := argValueWithMkdir(context, key)
	require.NoError(t, err)
	require.Equal(t, path, val)
	defer os.RemoveAll(parent)
	info, _ = os.Stat(path)
	require.NotNil(t, info)
	require.True(t, info.IsDir())
}

func Test_argValueWithMkdir_DirAlreadyExists(t *testing.T) {
	path := "README.md"
	key := "key"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	context := mocks.NewMockContext(ctrl)
	context.EXPECT().String(key).Return(path)

	_, err := argValueWithMkdir(context, key)
	require.Error(t, err)
}
