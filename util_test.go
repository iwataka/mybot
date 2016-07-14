package main

import (
	"testing"
)

func TestFormatUrl(t *testing.T) {
	var expected string
	var actual string
	var err error

	expected = "http://foo.com/bar/dest"
	actual, err = formatUrl("http://foo.com", "bar/dest")
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Fatalf("Expecting to have %s but got: %s", expected, actual)
	}

	expected = "http://foo.com/dest"
	actual, err = formatUrl("http://foo.com/bar", "/dest")
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Fatalf("Expecting to have %s but got: %s", expected, actual)
	}

	expected = "http://foo.com/bar/dest"
	actual, err = formatUrl("http://foo.com", "http://foo.com/bar/dest")
	if err != nil {
		t.Fatal(err)
	}
	if expected != actual {
		t.Fatalf("Expecting to have %s but got: %s", expected, actual)
	}
}
