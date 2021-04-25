package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/worker"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
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

func Test_cliFlagsAreUnique(t *testing.T) {
	configDir := "~/.config"
	cacheDir := "~/.cache"
	allFlags := []cli.Flag{}
	allFlags = append(allFlags, getCommonFlags(configDir, cacheDir)...)
	allFlags = append(allFlags, getServeSpecificFlags(configDir)...)
	allFlags = append(allFlags, getValidateSpecificFlags()...)
	for i, f1 := range allFlags {
		for j := i + 1; j < len(allFlags); j++ {
			f2 := allFlags[j]
			require.NotEqual(t, f2.GetName(), f1.GetName())
		}
	}
}

func Test_userSpecificData_statuses(t *testing.T) {
	ctrl := gomock.NewController(t)
	w := mocks.NewMockWorker(ctrl)
	wm := worker.NewWorkerManager(w, 0)
	d := userSpecificData{
		workerMgrs: map[int]*worker.WorkerManager{
			twitterUserRoutineKey: wm,
		},
	}
	statuses := d.statuses()
	require.Len(t, statuses, 4)
	for _, s := range statuses {
		require.False(t, s)
	}
}

func Test_userSpecificData_statuses_withEmptyData(t *testing.T) {
	d := userSpecificData{}
	statuses := d.statuses()
	require.Len(t, statuses, 4)
	for _, s := range statuses {
		require.False(t, s)
	}
}

func Test_userSpecificData_delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	w := mocks.NewMockWorker(ctrl)
	wm := worker.NewWorkerManager(w, 0)
	cacheFile := filepath.Join(os.TempDir(), "cache.json")
	cache, err := data.NewFileCache(cacheFile)
	require.NoError(t, err)
	d := userSpecificData{
		workerMgrs: map[int]*worker.WorkerManager{
			twitterUserRoutineKey: wm,
		},
		cache: cache,
	}
	require.NoError(t, d.cache.Save())
	require.NoError(t, d.delete())
	_, err = os.Stat(cacheFile)
	require.Error(t, err)
	require.Empty(t, d.workerMgrs)
}

func Test_userSpecificData_delete_withError(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := mocks.NewMockCache(ctrl)
	c.EXPECT().Delete().Return(errors.New("error"))
	d := userSpecificData{cache: c}
	require.Error(t, d.delete())
}

func Test_newFileCache(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	_, err = newFileCache(dir, "userID")
	require.NoError(t, err)
	rmdir(t, dir)
}

func Test_newFileConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	_, err = newFileConfig(dir, "userID")
	require.NoError(t, err)
	rmdir(t, dir)
}

func Test_newFileOAuthCreds(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	_, err = newFileOAuthCreds(dir, "userID")
	require.NoError(t, err)
	rmdir(t, dir)
}

func rmdir(t *testing.T, dir string) {
	info, err := os.Stat(dir)
	require.NoError(t, err)
	require.True(t, info.IsDir())
	require.NoError(t, os.Remove(dir))
}
