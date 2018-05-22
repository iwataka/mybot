package main

import (
	"io/ioutil"
	"testing"

	"github.com/iwataka/mybot/data"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	f, err := ioutil.TempFile("", "mybot")
	assert.NoError(t, err)

	path := f.Name()
	c, err := data.NewFileCache(path)
	assert.NoError(t, err)

	err = c.Save()
	assert.NoError(t, err)

	c, err = data.NewFileCache(path)
	assert.NoError(t, err)
}
