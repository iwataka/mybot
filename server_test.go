package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/mocks"
	"github.com/markbates/goth"
	"github.com/sclevine/agouti"
)

type TestAuthenticator struct{}

func (a *TestAuthenticator) SetProvider(req *http.Request, name string) {
}

func (a *TestAuthenticator) InitProvider(host, name string) {
}

func (a *TestAuthenticator) CompleteUserAuth(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error) {
	return goth.User{Name: "foo", NickName: "bar", UserID: "1234"}, nil
}

func (a *TestAuthenticator) Logout(provider string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func TestGetConfig(t *testing.T) {
	tmpAuth := auth
	defer func() { auth = tmpAuth }()
	auth = &TestAuthenticator{}

	s := httptest.NewServer(http.HandlerFunc(getConfig))
	defer s.Close()

	tmpCfg := config
	config = mybot.NewTestFileConfig("lib/testdata/config.template.toml", t)
	defer func() { config = tmpCfg }()

	testGet(t, s.URL, "Get /config")
}

func TestGetLog(t *testing.T) {
	tmpAuth := auth
	defer func() { auth = tmpAuth }()
	auth = &TestAuthenticator{}

	s := httptest.NewServer(http.HandlerFunc(getLog))
	defer s.Close()

	testGet(t, s.URL, "Get /log")
}

func TestGetStatus(t *testing.T) {
	tmpAuth := auth
	defer func() { auth = tmpAuth }()
	auth = &TestAuthenticator{}

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
	twitterApp = &mybot.FileTwitterOAuthApp{}
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
	tmpAuth := auth
	defer func() { auth = tmpAuth }()
	auth = &TestAuthenticator{}

	tmpCfg := config
	c := mybot.NewTestFileConfig("lib/testdata/config.template.toml", t)
	config = mybot.NewTestFileConfig("lib/testdata/config.template.toml", t)
	defer func() { config = tmpCfg }()

	tmpStatus := status
	status = mybot.NewStatus()
	defer func() { status = tmpStatus }()

	wg := new(sync.WaitGroup)
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == methodPost {
			config = mybot.NewTestFileConfig("", t)
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
	testPostConfigSingleDelete(t, s.URL, page, wg, c)
	testPostConfigDoubleDelete(t, s.URL, page, wg, c)
	testPostConfigError(t, s.URL, page, wg, c)
	testPostConfigTagsInput(t, s.URL, page, wg, c)
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

	if !reflect.DeepEqual(c.GetProperties(), config.GetProperties()) {
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

	if len(config.GetTwitterTimelines()) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.GetTwitterTimelines()))
	}
	if len(config.GetTwitterFavorites()) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.GetTwitterFavorites()))
	}
	if len(config.GetTwitterSearches()) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.GetTwitterSearches()))
	}
	if len(config.GetSlackMessages()) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.GetSlackMessages()))
	}
	if len(config.GetIncomingWebhooks()) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(config.GetIncomingWebhooks()))
	}

	config = c
}

func testPostConfigSingleDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c *mybot.FileConfig,
) {
	if err := page.Navigate(url); err != nil {
		t.Fatal(err)
	}

	if err := page.AllByButton("Delete").At(0).Click(); err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	if err := page.FindByID("overwrite").Submit(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()

	if len(config.GetTwitterTimelines()) != len(c.GetTwitterTimelines())-1 {
		t.Fatalf("%s's length is not %d", config.GetTwitterTimelines(), len(c.GetTwitterTimelines())-1)
	}

	config = c
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

	if !reflect.DeepEqual(c.GetProperties(), config.GetProperties()) {
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

	if !reflect.DeepEqual(c.GetProperties(), config.GetProperties()) {
		t.Fatalf("%v expected but %v found", c, config)
	}
}

func testPostConfigTagsInput(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c *mybot.FileConfig,
) {
	if err := page.Navigate(url); err != nil {
		t.Fatal(err)
	}

	name := "twitter.timelines.screen_names"
	keys := "foo,bar"
	if err := page.AllByName(name).SendKeys(keys); err == nil {
		t.Fatal("Tagsinput data-role elements must be uneditable currently")
	}
}

func TestPostIncoming(t *testing.T) {
	slackAPI = mybot.NewSlackAPI("", nil, nil)

	tmpCfg := config
	config = mybot.NewTestFileConfig("lib/testdata/config.template.toml", t)
	defer func() { config = tmpCfg }()

	s := httptest.NewServer(http.HandlerFunc(hooksHandler))
	defer s.Close()

	dest := config.GetIncomingWebhooks()[0].Endpoint
	buf := new(bytes.Buffer)
	buf.WriteString(`{"text": "foo"}`)
	testPost(t, s.URL+dest, "application/json", buf, fmt.Sprintf("%s %s", "POST", dest))
}

func TestPostConfigTimelineAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(config.GetTwitterTimelines()) },
		addTimelineConfig,
		"message",
	)
}

func TestPostConfigFavoriteAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(config.GetTwitterFavorites()) },
		addFavoriteConfig,
		"message",
	)
}

func TestPostConfigSearchAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(config.GetTwitterSearches()) },
		addSearchConfig,
		"message",
	)
}

func TestPostConfigMessageAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(config.GetSlackMessages()) },
		addMessageConfig,
		"message",
	)
}

func TestPostConfigIncomingAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(config.GetIncomingWebhooks()) },
		addIncomingConfig,
		"incoming",
	)
}

func testPostConfigAdd(
	t *testing.T,
	length func() int,
	add func(),
	name string,
) {
	tmpCfg := config
	config = mybot.NewTestFileConfig("lib/testdata/config.template.toml", t)
	defer func() { config = tmpCfg }()

	prev := length()
	add()
	cur := length()
	if cur != (prev + 1) {
		t.Fatalf("Failed to add %s", name)
	}

	_, err := configPage("", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetIndex(t *testing.T) {
	tmpAuth := auth
	defer func() { auth = tmpAuth }()
	auth = &TestAuthenticator{}

	ctrl := gomock.NewController(t)
	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	user := anaconda.User{}
	user.Name = "foo"
	twitterAPIMock.EXPECT().GetSelf(gomock.Any()).Return(user, nil)
	listResult := anaconda.CollectionListResult{}
	twitterAPIMock.EXPECT().GetCollectionListByUserId(gomock.Any(), gomock.Any()).Return(listResult, nil)
	twitterAPI = &mybot.TwitterAPI{API: twitterAPIMock, Cache: nil, Config: nil}

	tmpCache := cache
	defer func() { cache = tmpCache }()
	cache = mybot.NewTestFileCache("", t)
	img := mybot.ImageCacheData{}
	cache.SetImage(img)

	s := httptest.NewServer(http.HandlerFunc(getIndex))
	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Status code: %s (%d)", res.Status, res.StatusCode)
	}
}
