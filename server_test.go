package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/worker"
	"github.com/markbates/goth"
	"github.com/sclevine/agouti"
)

const (
	serverTestTwitterUserID = "123456"
)

var (
	serverTestUserSpecificData *userSpecificData
	serverTestTwitterUser      = goth.User{Name: "foo", NickName: "bar", UserID: serverTestTwitterUserID}
)

type TestAuthenticator struct{}

func (a *TestAuthenticator) SetProvider(req *http.Request, name string) {
}

func (a *TestAuthenticator) InitProvider(host, name, callback string) {
}

func (a *TestAuthenticator) CompleteUserAuth(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error) {
	return serverTestTwitterUser, nil
}

func (a *TestAuthenticator) Logout(provider string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func init() {
	err := initServer()
	if err != nil {
		panic(err)
	}

	serverTestUserSpecificData = &userSpecificData{}
	serverTestUserSpecificData.config, err = mybot.NewFileConfig("lib/testdata/config.template.toml")
	serverTestUserSpecificData.statuses = map[int]*bool{}
	serverTestUserSpecificData.workerChans = map[int]chan *worker.WorkerSignal{}
	serverTestUserSpecificData.slackAPI = mybot.NewSlackAPI("", serverTestUserSpecificData.config, nil)
	initStatuses(serverTestUserSpecificData.statuses)
	if err != nil {
		panic(err)
	}
	userSpecificDataMap[twitterUserIDPrefix+serverTestTwitterUserID] = serverTestUserSpecificData
}

func TestGetConfig(t *testing.T) {
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = &TestAuthenticator{}

	s := httptest.NewServer(http.HandlerFunc(configHandler))
	defer s.Close()

	testGet(t, s.URL, "Get /config")
}

func TestGetLog(t *testing.T) {
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = &TestAuthenticator{}

	s := httptest.NewServer(http.HandlerFunc(logHandler))
	defer s.Close()

	testGet(t, s.URL, "Get /log")
}

func TestGetStatus(t *testing.T) {
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = &TestAuthenticator{}

	s := httptest.NewServer(http.HandlerFunc(statusHandler))
	defer s.Close()

	testGet(t, s.URL, "Get /status")
}

func TestGetSetupTwitter(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(getSetup))
	defer s.Close()

	tmpTwitterApp := twitterApp
	var err error
	twitterApp, err = mybot.NewFileTwitterOAuthApp("")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { twitterApp = tmpTwitterApp }()

	tmpSlackApp := slackApp
	slackApp, err = mybot.NewFileOAuthApp("")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { slackApp = tmpSlackApp }()

	testGet(t, s.URL, "Get /setup/")
}

func TestGetConfigFile(t *testing.T) {
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = &TestAuthenticator{}

	s := httptest.NewServer(http.HandlerFunc(configFileHandler))
	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	assertHTTPResponse(t, res, "GET /config/file/")
	strings.Contains(res.Header.Get("Content-Type"), "application/force-download")
	strings.Contains(res.Header.Get("Content-Disposition"), ".json")
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := mybot.NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(bs, cfg)
	if err != nil {
		t.Fatal(err)
	}
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

func testPostConfig(t *testing.T, f func(*testing.T, string, *agouti.Page, *sync.WaitGroup, mybot.Config)) {
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = &TestAuthenticator{}

	wg := new(sync.WaitGroup)
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			serverTestUserSpecificData.config = mybot.NewTestFileConfig("", t)
			postConfig(w, r, serverTestUserSpecificData.config, serverTestTwitterUser)
			wg.Done()
		} else if r.Method == http.MethodGet {
			getConfig(w, r, serverTestUserSpecificData.config, serverTestUserSpecificData.slackAPI, serverTestTwitterUser)
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

	f(t, s.URL, page, wg, serverTestUserSpecificData.config)
}

func TestPostConfigWithoutModification(t *testing.T) {
	testPostConfig(t, testPostConfigWithoutModification)
}

func testPostConfigWithoutModification(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c mybot.Config,
) {
	if err := page.Navigate(url); err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	if err := page.FindByID("overwrite").Submit(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()

	if !reflect.DeepEqual(c.GetProperties(), serverTestUserSpecificData.config.GetProperties()) {
		t.Fatalf("%v expected but %v found", c, serverTestUserSpecificData.config)
	}
}

func TestPostConfigDelete(t *testing.T) {
	testPostConfig(t, testPostConfigDelete)
}

func testPostConfigDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c mybot.Config,
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

	if len(serverTestUserSpecificData.config.GetTwitterTimelines()) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(serverTestUserSpecificData.config.GetTwitterTimelines()))
	}
	if len(serverTestUserSpecificData.config.GetTwitterFavorites()) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(serverTestUserSpecificData.config.GetTwitterFavorites()))
	}
	if len(serverTestUserSpecificData.config.GetTwitterSearches()) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(serverTestUserSpecificData.config.GetTwitterSearches()))
	}
	if len(serverTestUserSpecificData.config.GetSlackMessages()) != 0 {
		t.Fatalf("%d expected but %d found", 0, len(serverTestUserSpecificData.config.GetSlackMessages()))
	}

	serverTestUserSpecificData.config = c
}

func TestPostConfigSingleDelete(t *testing.T) {
	testPostConfig(t, testPostConfigSingleDelete)
}

func testPostConfigSingleDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c mybot.Config,
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

	if len(serverTestUserSpecificData.config.GetTwitterTimelines()) != len(c.GetTwitterTimelines())-1 {
		t.Fatalf("%s's length is not %d", serverTestUserSpecificData.config.GetTwitterTimelines(), len(c.GetTwitterTimelines())-1)
	}

	serverTestUserSpecificData.config = c
}

func TestPostConfigDoubleDelete(t *testing.T) {
	testPostConfig(t, testPostConfigDoubleDelete)
}

func testPostConfigDoubleDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c mybot.Config,
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

	if !reflect.DeepEqual(c.GetProperties(), serverTestUserSpecificData.config.GetProperties()) {
		t.Fatalf("%v expected but %v found", c, serverTestUserSpecificData.config)
	}
}

func TestPostConfigError(t *testing.T) {
	testPostConfig(t, testPostConfigError)
}

func testPostConfigError(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c mybot.Config,
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

	if !reflect.DeepEqual(c.GetProperties(), serverTestUserSpecificData.config.GetProperties()) {
		t.Fatalf("%v expected but %v found", c, serverTestUserSpecificData.config)
	}
}

func TestPostConfigTagsInput(t *testing.T) {
	testPostConfig(t, testPostConfigTagsInput)
}

func testPostConfigTagsInput(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
	c mybot.Config,
) {
	_, err := net.DialTimeout("tcp", "cdnjs.cloudflare.com", 30*time.Second)
	if err != nil {
		t.Skip("Skip because network is unavailable: ", err)
	}

	if err := page.Navigate(url); err != nil {
		t.Fatal(err)
	}

	name := "twitter.timelines.screen_names"
	keys := "foo,bar"
	if err := page.AllByName(name).SendKeys(keys); err == nil {
		t.Fatal("Tagsinput data-role elements must be uneditable currently")
	}
}

func TestPostConfigTimelineAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(serverTestUserSpecificData.config.GetTwitterTimelines()) },
		func() { addTimelineConfig(serverTestUserSpecificData.config) },
		"message",
	)
}

func TestPostConfigFavoriteAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(serverTestUserSpecificData.config.GetTwitterFavorites()) },
		func() { addFavoriteConfig(serverTestUserSpecificData.config) },
		"message",
	)
}

func TestPostConfigSearchAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(serverTestUserSpecificData.config.GetTwitterSearches()) },
		func() { addSearchConfig(serverTestUserSpecificData.config) },
		"message",
	)
}

func TestPostConfigMessageAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(serverTestUserSpecificData.config.GetSlackMessages()) },
		func() { addMessageConfig(serverTestUserSpecificData.config) },
		"message",
	)
}

func testPostConfigAdd(
	t *testing.T,
	length func() int,
	add func(),
	name string,
) {
	prev := length()
	add()
	cur := length()
	if cur != (prev + 1) {
		t.Fatalf("Failed to add %s", name)
	}

	_, err := configPage("", "", "", "", serverTestUserSpecificData.config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetIndex(t *testing.T) {
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = &TestAuthenticator{}

	ctrl := gomock.NewController(t)
	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	user := anaconda.User{}
	user.Name = "foo"
	twitterAPIMock.EXPECT().GetSelf(gomock.Any()).Return(user, nil)
	listResult := anaconda.CollectionListResult{}
	twitterAPIMock.EXPECT().GetCollectionListByUserId(gomock.Any(), gomock.Any()).Return(listResult, nil)
	serverTestUserSpecificData.twitterAPI = &mybot.TwitterAPI{API: twitterAPIMock, Cache: nil, Config: nil}

	tmpCache := serverTestUserSpecificData.cache
	defer func() { serverTestUserSpecificData.cache = tmpCache }()
	serverTestUserSpecificData.cache = mybot.NewTestFileCache("", t)
	img := mybot.ImageCacheData{}
	serverTestUserSpecificData.cache.SetImage(img)

	s := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Status code: %s (%d)", res.Status, res.StatusCode)
	}
}
