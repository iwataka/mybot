package tmpl

import (
	"github.com/stretchr/testify/require"

	"fmt"
	"testing"
)

func Test_GetBoolSelectboxValue(t *testing.T) {
	val := make(map[string][]string)
	name := "foo"
	val[name] = []string{"true", "undefined"}

	var result *bool
	trueVal := true

	result = GetBoolSelectboxValue(val, 0, name)
	require.Equal(t, &trueVal, result)

	result = GetBoolSelectboxValue(val, 1, name)
	require.Nil(t, result)

	result = GetBoolSelectboxValue(val, 2, name)
	require.Nil(t, result)
}

func Test_GetListTextboxValue(t *testing.T) {
	val := make(map[string][]string)
	name := "foo"
	item1 := "fizz"
	item2 := "buzz"
	val[name] = []string{fmt.Sprintf("%s, %s ", item1, item2), ""}

	var result []string

	result = GetListTextboxValue(val, 0, name)
	require.Equal(t, []string{item1, item2}, result)

	result = GetListTextboxValue(val, 1, name)
	require.Equal(t, []string{}, result)

	result = GetListTextboxValue(val, 2, name)
	require.Equal(t, []string{}, result)
}

func Test_GetFloat64Ptr(t *testing.T) {
	val := make(map[string][]string)
	name := "foo"
	val[name] = []string{"1.23", "", "foo"}

	var result *float64
	var err error
	fval := 1.23

	result, err = GetFloat64Ptr(val, 0, name)
	require.NoError(t, err)
	require.Equal(t, &fval, result)

	result, err = GetFloat64Ptr(val, 1, name)
	require.NoError(t, err)
	require.Nil(t, result)

	result, err = GetFloat64Ptr(val, 2, name)
	require.Error(t, err)
	require.Nil(t, result)

	result, err = GetFloat64Ptr(val, 3, name)
	require.NoError(t, err)
	require.Nil(t, result)
}

func Test_GetIntPtr(t *testing.T) {
	val := make(map[string][]string)
	name := "foo"
	val[name] = []string{"1", "", "foo"}

	var result *int
	var err error
	ival := 1

	result, err = GetIntPtr(val, 0, name)
	require.NoError(t, err)
	require.Equal(t, &ival, result)

	result, err = GetIntPtr(val, 1, name)
	require.NoError(t, err)
	require.Nil(t, result)

	result, err = GetIntPtr(val, 2, name)
	require.Error(t, err)
	require.Nil(t, result)

	result, err = GetIntPtr(val, 3, name)
	require.NoError(t, err)
	require.Nil(t, result)
}

func Test_GetString(t *testing.T) {
	val := make(map[string][]string)
	value := []string{"foo"}
	key := "key"
	val[key] = value

	var result string
	def := "default"

	result = GetString(val, key, 0, def)
	require.Equal(t, value[0], result)

	result = GetString(val, "other_key", 0, def)
	require.Equal(t, def, result)

	result = GetString(val, key, 1, def)
	require.Equal(t, def, result)
}

func Test_NewMap(t *testing.T) {
	key1 := "key1"
	key2 := "key2"
	val1 := 1
	val2 := []string{"foo"}
	m := NewMap(key1, val1, key2, val2)

	require.Equal(t, val1, m[key1])
	require.Equal(t, val2, m[key2])
	require.Panics(t, func() { NewMap(key1, key2, val1) })
}
