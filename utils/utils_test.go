package utils_test

import (
	. "github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"os"
	"path/filepath"
	"testing"
)

func TestDecodeFile_TakesJsonFile(t *testing.T) {
	arr := []int{}
	assert.NoError(t, DecodeFile("testdata/array.json", &arr))
	assert.NoError(t, DecodeFile("testdata/empty.json", nil))
	assert.Error(t, DecodeFile("testdata/invalid.json", &arr))
	require.NoError(t, os.Chmod("testdata/unreadable.json", 0222))
	assert.NoError(t, DecodeFile("testdata/unreadable.json", nil))
	require.NoError(t, os.Chmod("testdata/unreadable.json", 0644))
}

func TestDecodeFile_TakesTomlFile(t *testing.T) {
	foo := struct {
		Arr []int
	}{
		[]int{},
	}
	assert.NoError(t, DecodeFile("testdata/array.toml", &foo))
	assert.NoError(t, DecodeFile("testdata/empty.toml", nil))
	assert.Error(t, DecodeFile("testdata/invalid.toml", &foo))
	err := DecodeFile("testdata/array.toml", &struct{}{})
	assert.Error(t, err)
	_, ok := err.(*TomlUndecodedKeysError)
	assert.True(t, ok)
}

func TestEncodeFile_TakesJsonFile(t *testing.T) {
	tmp := "tmp"
	require.NoError(t, os.Mkdir(tmp, os.FileMode(0777)))
	defer os.RemoveAll(tmp)

	arr := []int{1, 2, 3}
	arrayJson := filepath.Join(tmp, "array.json")
	assert.NoError(t, EncodeFile(arrayJson, arr))

	unwritableJson := filepath.Join(tmp, "unwritable.json")
	_, err := os.Create(unwritableJson)
	require.NoError(t, err)
	require.NoError(t, os.Chmod(unwritableJson, 0555))
	assert.Error(t, EncodeFile(unwritableJson, arr))

	ch := make(chan bool)
	assert.Error(t, EncodeFile(arrayJson, ch))
}

func TestEncodeFile_TakesTomlFile(t *testing.T) {
	tmp := "tmp"
	require.NoError(t, os.Mkdir(tmp, os.FileMode(0777)))
	defer os.RemoveAll(tmp)

	m := map[string]int{"foo": 0}
	mapToml := filepath.Join(tmp, "map.toml")
	assert.NoError(t, EncodeFile(mapToml, m))

	arr := []int{1, 2, 3}
	arrayToml := filepath.Join(tmp, "array.toml")
	assert.Error(t, EncodeFile(arrayToml, arr))
}

func TestCalcStringSlices(t *testing.T) {
	s1 := []string{"foo", "bar"}
	s2 := []string{"foo", "else"}

	addResult := CalcStringSlices(s1, s2, true)
	assert.Len(t, addResult, 3)

	subResult := CalcStringSlices(s1, s2, false)
	assert.Len(t, subResult, 1)
}

func TestCalcBools(t *testing.T) {
	assert.True(t, CalcBools(true, false, true))
	assert.False(t, CalcBools(true, true, false))
}

func TestCheckStringCotnained(t *testing.T) {
	ss := []string{"foo", "bar"}
	str := "foo"
	assert.True(t, CheckStringContained(ss, str))
	str = "else"
	assert.False(t, CheckStringContained(ss, str))
}

func TestGenerateRandString(t *testing.T) {
	assert.Len(t, GenerateRandString(0), 0)
	assert.Len(t, GenerateRandString(1), 1)
	assert.Len(t, GenerateRandString(10), 10)
}

func TestExitIfError_ErrorNotFound(t *testing.T) {
	ExitIfError(nil)
}

func TestTruePtr(t *testing.T) {
	require.True(t, *TruePtr())
}

func TestFalsePtr(t *testing.T) {
	require.False(t, *FalsePtr())
}

func TestIntPtr(t *testing.T) {
	n := 100
	require.Equal(t, n, *IntPtr(n))
}

func TestFloat64Ptr(t *testing.T) {
	var f float64 = 1.1
	require.Equal(t, f, *Float64Ptr(f))
}
