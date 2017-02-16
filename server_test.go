package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iwataka/mybot/lib"
)

type LoggerMock struct {
}

func (l *LoggerMock) Println(v ...interface{}) error {
	return nil
}

func (l *LoggerMock) HandleError(err error) {
	return
}

func (l *LoggerMock) ReadString() string {
	return "foo"
}

func TestGetConfig(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(getConfig))
	defer s.Close()

	tmpCfg := config
	c, err := mybot.NewFileConfig("lib/test_assets/config.template.toml")
	if err != nil {
		t.Fatal(err)
	}
	config = c
	defer func() { config = tmpCfg }()

	testGet(t, s.URL, "Get /config")
}

func TestGetLog(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(getLog))
	defer s.Close()

	tmpLogger := logger
	logger = &LoggerMock{}
	defer func() { logger = tmpLogger }()

	testGet(t, s.URL, "Get /log")
}

func TestGetStatus(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(getStatus))
	defer s.Close()

	tmpStatus := status
	status = mybot.NewStatus()
	defer func() { status = tmpStatus }()

	testGet(t, s.URL, "Get /status")
}

func TestGetSetupTwitter(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(getSetupTwitter))
	defer s.Close()

	tmpTwitterApp := twitterApp
	twitterApp = &mybot.OAuthApp{}
	defer func() { twitterApp = tmpTwitterApp }()

	testGet(t, s.URL, "Get /setup/twitter")
}

func testGet(t *testing.T, url string, msg string) {
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	assertHTTPResponse(t, res, msg)
}

func testPost(t *testing.T, url string, bodyType string, body io.Reader, msg string) {
	res, err := http.Post(url, bodyType, body)
	if err != nil {
		t.Fatal(err)
	}
	assertHTTPResponse(t, res, msg)
}

func assertHTTPResponse(t *testing.T, res *http.Response, msg string) {
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Error code %d: %s", res.StatusCode, msg)
	}
}
