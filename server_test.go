package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/iwataka/mybot/lib"
	"github.com/sclevine/agouti"
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
	config = newFileConfig("lib/testdata/config.template.toml", t)
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

func TestPostConfig(t *testing.T) {
	tmpCfg := config
	c := newFileConfig("lib/testdata/config.template.toml", t)
	config = newFileConfig("lib/testdata/config.template.toml", t)
	defer func() { config = tmpCfg }()

	wg := new(sync.WaitGroup)
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == methodPost {
			*config = *newFileConfig("", t)
			postConfig(w, r)
			wg.Done()
		} else if r.Method == methodGet {
			getConfig(w, r)
		}
	}

	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	driver := agouti.PhantomJS()
	if err := driver.Start(); err != nil {
		t.Fatal(err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		t.Fatal(err)
	}

	testPostConfig(t, s.URL, page, wg, c)
	testPostConfigDelete(t, s.URL, page, wg, c)
	testPostConfigDoubleDelete(t, s.URL, page, wg, c)
	testPostConfigError(t, s.URL, page, wg, c)
}

func testPostConfig(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c *mybot.FileConfig,
) {
	if err := page.Navigate(url); err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	if err := page.FindByID("overwrite").Submit(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()

	c.File = config.File
	if !reflect.DeepEqual(c, config) {
		t.Fatalf("%v expected but %v found", c, config)
	}
}

func testPostConfigDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c *mybot.FileConfig,
) {
	if err := page.Navigate(url); err != nil {
		t.Fatal(err)
	}

	if err := page.AllByButton("Delete").Click(); err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	if err := page.FindByID("overwrite").Submit(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()

	if len(config.Twitter.Timelines) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.Twitter.Timelines))
	}
	if len(config.Twitter.Favorites) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.Twitter.Favorites))
	}
	if len(config.Twitter.Searches) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.Twitter.Searches))
	}
	if len(config.Twitter.APIs) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.Twitter.APIs))
	}
	if len(config.Slack.Messages) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.Slack.Messages))
	}

	*config = *c
}

func testPostConfigDoubleDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c *mybot.FileConfig,
) {
	if err := page.Navigate(url); err != nil {
		t.Fatal(err)
	}

	if err := page.AllByButton("Delete").DoubleClick(); err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	if err := page.FindByID("overwrite").Submit(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()

	c.File = config.File
	if !reflect.DeepEqual(c, config) {
		t.Fatalf("%v expected but %v found", c, config)
	}
}

func testPostConfigError(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c *mybot.FileConfig,
) {
	if err := page.Navigate(url); err != nil {
		t.Fatal(err)
	}

	if err := page.AllByName("twitter.timelines.count").Fill("foo"); err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	if err := page.FindByID("overwrite").Submit(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()

	c.File = config.File
	if !reflect.DeepEqual(c, config) {
		t.Fatalf("%v expected but %v found", c, config)
	}
}

func newFileConfig(path string, t *testing.T) *mybot.FileConfig {
	c, err := mybot.NewFileConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
