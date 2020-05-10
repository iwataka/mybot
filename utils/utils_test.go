package utils_test

import (
	"github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/require"

	"os"
	"path/filepath"
	"testing"
)

func Test_DecodeFile_TakesJsonFile(t *testing.T) {
	arr := []int{}
	require.NoError(t, utils.DecodeFile("testdata/array.json", &arr))
	require.NoError(t, utils.DecodeFile("testdata/empty.json", nil))
	require.Error(t, utils.DecodeFile("testdata/invalid.json", &arr))
	require.NoError(t, os.Chmod("testdata/unreadable.json", 0222))
	require.NoError(t, utils.DecodeFile("testdata/unreadable.json", nil))
	require.NoError(t, os.Chmod("testdata/unreadable.json", 0644))
}

func Test_DecodeFile_TakesTomlFile(t *testing.T) {
	foo := struct {
		Arr []int
	}{
		[]int{},
	}
	require.NoError(t, utils.DecodeFile("testdata/array.toml", &foo))
	require.NoError(t, utils.DecodeFile("testdata/empty.toml", nil))
	require.Error(t, utils.DecodeFile("testdata/invalid.toml", &foo))
	err := utils.DecodeFile("testdata/array.toml", &struct{}{})
	require.Error(t, err)
	e, ok := err.(*utils.TomlUndecodedKeysError)
	require.Equal(t, e.File, "testdata/array.toml")
	require.True(t, ok)
}

func Test_DecodeFile_TakesYamlFile(t *testing.T) {
	foo := struct {
		Parent struct {
			Child1 int `yaml:"child1"`
			Child2 int `yaml:"child2"`
		} `yaml:"parent"`
	}{}
	require.NoError(t, utils.DecodeFile("testdata/map.yaml", &foo))
}

func Test_EncodeFile_TakesJsonFile(t *testing.T) {
	tmp := "tmp"
	require.NoError(t, os.Mkdir(tmp, os.FileMode(0777)))
	defer os.RemoveAll(tmp)

	arr := []int{1, 2, 3}
	arrayJson := filepath.Join(tmp, "array.json")
	require.NoError(t, utils.EncodeFile(arrayJson, arr))

	unwritableJson := filepath.Join(tmp, "unwritable.json")
	_, err := os.Create(unwritableJson)
	require.NoError(t, err)
	require.NoError(t, os.Chmod(unwritableJson, 0555))
	require.Error(t, utils.EncodeFile(unwritableJson, arr))

	ch := make(chan bool)
	require.Error(t, utils.EncodeFile(arrayJson, ch))
}

func Test_EncodeFile_TakesTomlFile(t *testing.T) {
	tmp := "tmp"
	require.NoError(t, os.Mkdir(tmp, os.FileMode(0777)))
	defer os.RemoveAll(tmp)

	m := map[string]int{"foo": 0}
	mapToml := filepath.Join(tmp, "map.toml")
	require.NoError(t, utils.EncodeFile(mapToml, m))

	arr := []int{1, 2, 3}
	arrayToml := filepath.Join(tmp, "array.toml")
	require.Error(t, utils.EncodeFile(arrayToml, arr))
}

func Test_EncodeFile_TakesYamlFile(t *testing.T) {
	tmp := "tmp"
	require.NoError(t, os.Mkdir(tmp, os.FileMode(0777)))
	defer os.RemoveAll(tmp)

	m := map[string]int{"foo": 0}
	mapYaml := filepath.Join(tmp, "map.yaml")
	require.NoError(t, utils.EncodeFile(mapYaml, m))
}

func Test_CalcStringSlices(t *testing.T) {
	s1 := []string{"foo", "bar"}
	s2 := []string{"foo", "else"}

	addResult := utils.CalcStringSlices(s1, s2, true)
	require.Len(t, addResult, 3)

	subResult := utils.CalcStringSlices(s1, s2, false)
	require.Len(t, subResult, 1)
}

func Test_CalcBools(t *testing.T) {
	require.True(t, utils.CalcBools(true, false, true))
	require.False(t, utils.CalcBools(true, true, false))
}

func Test_CheckStringCotnained(t *testing.T) {
	ss := []string{"foo", "bar"}
	str := "foo"
	require.True(t, utils.CheckStringContained(ss, str))
	str = "else"
	require.False(t, utils.CheckStringContained(ss, str))
}

func Test_GenerateRandString(t *testing.T) {
	require.Len(t, utils.GenerateRandString(0), 0)
	require.Len(t, utils.GenerateRandString(1), 1)
	require.Len(t, utils.GenerateRandString(10), 10)
}

func Test_ExitIfError_ErrorNotFound(t *testing.T) {
	utils.ExitIfError(nil)
}

func Test_TruePtr(t *testing.T) {
	require.True(t, *utils.TruePtr())
}

func Test_FalsePtr(t *testing.T) {
	require.False(t, *utils.FalsePtr())
}

func Test_IntPtr(t *testing.T) {
	n := 100
	require.Equal(t, n, *utils.IntPtr(n))
}

func Test_Float64Ptr(t *testing.T) {
	var f float64 = 1.1
	require.Equal(t, f, *utils.Float64Ptr(f))
}

func Test_UniqStrSlice(t *testing.T) {
	s1 := []string{"foo", "foo"}
	require.Equal(t, []string{"foo"}, utils.UniqStrSlice(s1))
	s2 := []string{"foo", "bar"}
	require.Equal(t, s2, utils.UniqStrSlice(s2))
}
