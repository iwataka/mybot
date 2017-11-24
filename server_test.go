package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/deep"
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

	if _, err := os.Stat("screenshots"); err != nil {
		err := os.Mkdir("screenshots", os.FileMode(0755))
		if err != nil {
			fmt.Println("Failed to make `screenshots` directory")
			os.Exit(1)
		}
	}
}

func TestTwitterColsPage(t *testing.T) {
	testTwitterCols(t, testTwitterColsPage)
}

func testTwitterColsPage(url string) error {
	driver := agouti.PhantomJS()
	if err := driver.Start(); err != nil {
		return err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		return err
	}

	err = page.Navigate(url)
	if err != nil {
		return err
	}

	err = page.Screenshot("screenshots/twitter-collections.png")
	if err != nil {
		return err
	}
	return nil
}

func TestGetTwitterCols(t *testing.T) {
	testTwitterCols(t, testGet)
}

func testTwitterCols(t *testing.T, f func(url string) error) {
	ctrl := gomock.NewController(t)

	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CompleteUserAuth(gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	fooCol := anaconda.Collection{
		Name:          "foo",
		CollectionUrl: "https://twitter.com/NYTNow/timelines/576828964162965504",
	}
	barCol := anaconda.Collection{
		Name:          "bar",
		CollectionUrl: "https://twitter.com/NYTNow/timelines/576828964162965504",
	}
	listResult := anaconda.CollectionListResult{}
	listResult.Objects.Timelines = map[string]anaconda.Collection{
		"fooID": fooCol,
		"barID": barCol,
	}
	twitterAPIMock.EXPECT().GetCollectionListByUserId(gomock.Any(), gomock.Any()).Times(2).Return(listResult, nil)
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = &mybot.TwitterAPI{API: twitterAPIMock, Cache: nil, Config: nil}

	s := httptest.NewServer(http.HandlerFunc(twitterColsHandler))
	defer s.Close()

	err := f(s.URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CompleteUserAuth(gomock.Any(), gomock.Any(), gomock.Any()).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	s := httptest.NewServer(http.HandlerFunc(configHandler))
	defer s.Close()

	err := testGet(s.URL)
	if err != nil {
		t.Fatal(err)
	}
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

	err = testGet(s.URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetConfigJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CompleteUserAuth(gomock.Any(), gomock.Any(), gomock.Any()).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	s := httptest.NewServer(http.HandlerFunc(configJsonHandler))
	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	err = checkHTTPResponse(res)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := mybot.NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	err = cfg.Unmarshal(bs)
	if err != nil {
		t.Fatal(err)
	}

	cfgProps := cfg.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	deep.IgnoreDifferenceBetweenEmptyMapAndNil = true
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = true
	if diff := deep.Equal(cfgProps, configProps); diff != nil {
		t.Fatal(diff)
	}
}

func TestGetConfigFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CompleteUserAuth(gomock.Any(), gomock.Any(), gomock.Any()).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	s := httptest.NewServer(http.HandlerFunc(configFileHandler))
	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	err = checkHTTPResponse(res)
	if err != nil {
		t.Fatal(err)
	}
	hasForceDownload := strings.Contains(res.Header.Get("Content-Type"), "application/force-download")
	if !hasForceDownload {
		t.Fatalf("It must have force-download but not")
	}
	hasContentDisposition := strings.Contains(res.Header.Get("Content-Disposition"), ".json")
	if !hasContentDisposition {
		t.Fatalf("It must have Content-Disposition but not")
	}
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

func testGet(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	return checkHTTPResponse(res)
}

func testPost(t *testing.T, url string, bodyType string, body io.Reader, msg string) {
	res, err := http.Post(url, bodyType, body)
	if err != nil {
		t.Fatal(err)
	}
	checkHTTPResponse(res)
}

func checkHTTPResponse(res *http.Response) error {
	if res.StatusCode != http.StatusOK {
		bs, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("%s %d", string(bs), res.StatusCode)
	}
	return nil
}

func testPostConfig(t *testing.T, f func(*testing.T, string, *agouti.Page, *sync.WaitGroup, mybot.Config)) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CompleteUserAuth(gomock.Any(), gomock.Any(), gomock.Any()).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	wg := new(sync.WaitGroup)
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/assets/js") {
			getAssetsJS(w, r)
		} else {
			if r.Method == http.MethodPost {
				serverTestUserSpecificData.config = mybot.NewTestFileConfig("", t)
				postConfig(w, r, serverTestUserSpecificData.config, serverTestTwitterUser)
				wg.Done()
			} else if r.Method == http.MethodGet {
				getConfig(w, r, serverTestUserSpecificData.config, serverTestUserSpecificData.slackAPI, serverTestTwitterUser)
			}
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

	cProps := c.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	deep.IgnoreDifferenceBetweenEmptyMapAndNil = true
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = true
	if diff := deep.Equal(cProps, configProps); diff != nil {
		t.Fatal(diff)
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

	cProps := c.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	if diff := deep.Equal(cProps.Slack, configProps.Slack); diff != nil {
		t.Fatal(diff)
	}
	if diff := deep.Equal(cProps.Twitter, configProps.Twitter); diff != nil {
		t.Fatal(diff)
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

	cProps := c.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	if diff := deep.Equal(cProps, configProps); diff != nil {
		t.Fatal(diff)
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
		configTimelineAddHandler,
		"message",
	)
}

func TestPostConfigFavoriteAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(serverTestUserSpecificData.config.GetTwitterFavorites()) },
		configFavoriteAddHandler,
		"message",
	)
}

func TestPostConfigSearchAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(serverTestUserSpecificData.config.GetTwitterSearches()) },
		configSearchAddHandler,
		"message",
	)
}

func TestPostConfigMessageAdd(t *testing.T) {
	testPostConfigAdd(
		t,
		func() int { return len(serverTestUserSpecificData.config.GetSlackMessages()) },
		configMessageAddHandler,
		"message",
	)
}

func testPostConfigAdd(
	t *testing.T,
	length func() int,
	handler func(w http.ResponseWriter, r *http.Request),
	name string,
) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CompleteUserAuth(gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	expectedErrMsg := "expected error"
	expectedErr := fmt.Errorf(expectedErrMsg)

	prev := length()
	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return expectedErr
		},
	}
	res, err := client.Post(s.URL, "", nil)
	if err != nil && !strings.HasSuffix(err.Error(), expectedErrMsg) {
		t.Fatal(err)
	}
	cur := length()
	if cur != (prev + 1) {
		t.Fatalf("Failed to add %s", name)
	}
	testResponseIsRedirect(t, res, "/config")
}

func TestIndexPage(t *testing.T) {
	testIndex(t, testIndexPage)
}

func testIndexPage(url string) error {
	driver := agouti.PhantomJS()
	if err := driver.Start(); err != nil {
		return err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		return err
	}

	err = page.Navigate(url)
	if err != nil {
		return err
	}

	err = page.Screenshot("screenshots/index.png")
	if err != nil {
		return err
	}
	return nil
}

func TestGetIndex(t *testing.T) {
	testIndex(t, testGet)
}

func testIndex(t *testing.T, f func(url string) error) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CompleteUserAuth(gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	tmpCache := serverTestUserSpecificData.cache
	defer func() { serverTestUserSpecificData.cache = tmpCache }()
	serverTestUserSpecificData.cache = mybot.NewTestFileCache("", t)
	img := mybot.ImageCacheData{}
	serverTestUserSpecificData.cache.SetImage(img)

	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	user := anaconda.User{ScreenName: "foo"}
	twitterAPIMock.EXPECT().GetSelf(gomock.Any()).Return(user, nil)
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = &mybot.TwitterAPI{API: twitterAPIMock, Cache: nil, Config: nil}

	s := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer s.Close()

	err := f(s.URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTwitterUserSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CompleteUserAuth(gomock.Any(), gomock.Any(), gomock.Any()).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	user1 := anaconda.User{Name: "foo"}
	user2 := anaconda.User{Name: "bar"}
	users := []anaconda.User{user1, user2}
	twitterAPIMock.EXPECT().GetUserSearch(gomock.Any(), gomock.Any()).Return(users, nil)
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = &mybot.TwitterAPI{API: twitterAPIMock, Cache: nil, Config: nil}

	s := httptest.NewServer(http.HandlerFunc(twitterUserSearchHandler))
	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	us := []anaconda.User{}
	err = json.Unmarshal(bs, &us)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(users, us); diff != nil {
		t.Fatal(diff)
	}
}

func testResponseIsRedirect(t *testing.T, res *http.Response, locPrefix string) {
	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("Status code is expected to be %d but %d", http.StatusSeeOther, res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, locPrefix) {
		t.Fatalf("Location header value should start with %s but %s", locPrefix, loc)
	}
}
