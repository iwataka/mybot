package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/deep"
	"github.com/iwataka/mybot/core"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/oauth"
	"github.com/iwataka/mybot/worker"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/sclevine/agouti"
	"github.com/stretchr/testify/require"
)

const (
	serverTestTwitterScreenName = "foo"
	serverTestUserProvider      = "provider"
	serverTestUserID            = "123456"
	screenshotsDir              = "screenshots"
)

var (
	serverTestUserSpecificData *userSpecificData
	serverTestTwitterUser      = goth.User{
		UserID:   serverTestUserID,
		Provider: serverTestUserProvider,
	}
	driver *agouti.WebDriver
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	if driver != nil {
		err := driver.Stop()
		if err != nil {
			panic(err)
		}
	}
	os.Exit(exitCode)
}

func getDriver() *agouti.WebDriver {
	if driver == nil {
		driver = agouti.PhantomJS()
		err := driver.Start()
		if err != nil {
			panic(err)
		}
	}
	return driver
}

func init() {
	var err error
	serverTestUserSpecificData = &userSpecificData{}
	// TODO: Mock config and other fields
	serverTestUserSpecificData.config, err = core.NewFileConfig("core/testdata/config.yaml")
	if err != nil {
		panic(err)
	}
	serverTestUserSpecificData.workerMgrs = map[int]*worker.WorkerManager{}
	serverTestUserSpecificData.slackAPI = core.NewSlackAPIWithAuth("", serverTestUserSpecificData.config, nil)
	userSpecificDataMap[fmt.Sprintf(appUserIDFormat, serverTestUserProvider, serverTestUserID)] = serverTestUserSpecificData

	if _, err := os.Stat(screenshotsDir); err != nil {
		err := os.Mkdir(screenshotsDir, os.FileMode(0755))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to make `%s` directory\n", screenshotsDir)
			os.Exit(1)
		}
	}

	deep.IgnoreDifferenceBetweenEmptyMapAndNil = true
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = true
}

func TestAuthenticator_SetProvider(t *testing.T) {
	expectedProvider := "twitter"
	auth := Authenticator{}
	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)

	auth.SetProvider(expectedProvider, req)
	provider, err := gothic.GetProviderName(req)
	require.NoError(t, err)
	require.Equal(t, expectedProvider, provider)

	anotherProvider := "slack"
	auth.SetProvider(anotherProvider, req)
	unchangedProvider, err := gothic.GetProviderName(req)
	require.NoError(t, err)
	require.Equal(t, expectedProvider, unchangedProvider)
}

func TestTwitterColsPage(t *testing.T) {
	skipIfDialTimeout(t, "twitter.com", "https", 30*time.Second)
	testTwitterCols(t, testTwitterColsPage)
}

func testTwitterColsPage(url string) error {
	page, err := getDriver().NewPage()
	if err != nil {
		return err
	}

	err = page.Navigate(url)
	if err != nil {
		return err
	}

	err = page.Screenshot(filepath.Join(screenshotsDir, "twitter-collections.png"))
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
	authMock.EXPECT().GetLoginUser(gomock.Any()).Times(2).Return(serverTestTwitterUser, nil)
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
	twitterAPIMock.EXPECT().GetSelf(gomock.Any()).AnyTimes().Return(anaconda.User{ScreenName: serverTestTwitterScreenName, IdStr: serverTestUserID}, nil)
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = core.NewTwitterAPI(twitterAPIMock, nil, nil)

	s := httptest.NewServer(http.HandlerFunc(twitterColsHandler))
	defer s.Close()

	err := f(s.URL)
	require.NoError(t, err)
}

func TestGetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().GetLoginUser(gomock.Any()).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	s := httptest.NewServer(http.HandlerFunc(configHandler))
	defer s.Close()

	err := testGet(s.URL)
	require.NoError(t, err)
}

func TestGetSetupTwitter(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(getSetup))
	defer s.Close()

	tmpTwitterApp := twitterApp
	var err error
	twitterApp, err = oauth.NewFileTwitterOAuthApp("")
	require.NoError(t, err)
	defer func() { twitterApp = tmpTwitterApp }()

	tmpSlackApp := slackApp
	slackApp, err = oauth.NewFileOAuthApp("")
	require.NoError(t, err)
	defer func() { slackApp = tmpSlackApp }()

	err = testGet(s.URL)
	require.NoError(t, err)
}

func TestGetConfigFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().GetLoginUser(gomock.Any()).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	s := httptest.NewServer(http.HandlerFunc(configFileHandler))
	defer s.Close()

	res, err := http.Get(s.URL)
	require.NoError(t, err)

	err = checkHTTPResponse(res)
	require.NoError(t, err)

	hasForceDownload := strings.Contains(res.Header.Get("Content-Type"), "application/force-download")
	require.True(t, hasForceDownload)

	hasContentDisposition := strings.Contains(res.Header.Get("Content-Disposition"), defaultConfigFormat)
	require.True(t, hasContentDisposition)

	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	cfg, err := core.NewFileConfig("")
	require.NoError(t, err)

	err = json.Unmarshal(bs, cfg)
	require.NoError(t, err)
}

func testGet(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	return checkHTTPResponse(res)
}

func checkHTTPResponse(res *http.Response) error {
	if res.StatusCode != http.StatusOK {
		bs, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("%s %d", string(bs), res.StatusCode)
	}
	return nil
}

func testPostConfig(t *testing.T, f func(*testing.T, string, *agouti.Page, *sync.WaitGroup)) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().GetLoginUser(gomock.Any()).Return(serverTestTwitterUser, nil)
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
				serverTestUserSpecificData.config = core.NewTestFileConfig("", t)
				postConfig(w, r, serverTestUserSpecificData.config, serverTestTwitterUser)
				wg.Done()
			} else if r.Method == http.MethodGet {
				getConfig(w, r, serverTestUserSpecificData.config, serverTestUserSpecificData.slackAPI, serverTestUserSpecificData.twitterAPI)
			}
		}
	}

	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	page, err := getDriver().NewPage()
	require.NoError(t, err)

	curUserData := serverTestUserSpecificData.config
	defer func() { serverTestUserSpecificData.config = curUserData }()
	f(t, s.URL, page, wg)
}

func TestPostConfigWithoutModification(t *testing.T) {
	skipIfDialTimeout(t, "twitter.com", "https", 30*time.Second)
	testPostConfig(t, testPostConfigWithoutModification)
}

func testPostConfigWithoutModification(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	c := serverTestUserSpecificData.config

	require.NoError(t, page.Navigate(url))

	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "config_before_post_without_modification.png")))
	wg.Add(1)
	require.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "config_after_post_without_modification.png")))

	cProps := c.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	deep.IgnoreDifferenceBetweenEmptyMapAndNil = true
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = true
	require.Nil(t, deep.Equal(&cProps, &configProps))
}

func TestPostConfigDelete(t *testing.T) {
	skipIfDialTimeout(t, "twitter.com", "https", 30*time.Second)
	testPostConfig(t, testPostConfigDelete)
}

func testPostConfigDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	require.NoError(t, page.Navigate(url))

	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "delete_config_before_post.png")))
	require.NoError(t, page.AllByClass("config-row-delete").Click())
	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "delete_config_after_click_delete_buttons.png")))
	wg.Add(1)
	require.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "delete_config_after_post.png")))

	require.Empty(t, serverTestUserSpecificData.config.GetTwitterTimelines())
	require.Empty(t, serverTestUserSpecificData.config.GetTwitterFavorites())
	require.Empty(t, serverTestUserSpecificData.config.GetTwitterSearches())
	require.Empty(t, serverTestUserSpecificData.config.GetSlackMessages())
}

func TestPostConfigSingleDelete(t *testing.T) {
	skipIfDialTimeout(t, "twitter.com", "https", 30*time.Second)
	testPostConfig(t, testPostConfigSingleDelete)
}

func testPostConfigSingleDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	c := serverTestUserSpecificData.config

	require.NoError(t, page.Navigate(url))

	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "single_delete_config_before_post.png")))
	require.NoError(t, page.AllByClass("config-row-delete").At(0).Click())
	wg.Add(1)
	require.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "single_delete_config_after_post.png")))

	require.Equal(t, len(serverTestUserSpecificData.config.GetTwitterTimelines()), len(c.GetTwitterTimelines())-1)
}

func TestPostConfigDoubleDelete(t *testing.T) {
	skipIfDialTimeout(t, "twitter.com", "https", 30*time.Second)
	testPostConfig(t, testPostConfigDoubleDelete)
}

func testPostConfigDoubleDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	c := serverTestUserSpecificData.config

	require.NoError(t, page.Navigate(url))

	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "double_delete_config_before_post.png")))
	require.NoError(t, page.AllByClass("config-row-delete").DoubleClick())
	wg.Add(1)
	require.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "double_delete_config_after_post.png")))

	cProps := c.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	require.Nil(t, deep.Equal(&cProps, &configProps))
}

func TestPostConfigNameError(t *testing.T) {
	skipIfDialTimeout(t, "twitter.com", "https", 30*time.Second)
	testPostConfig(t, testPostConfigNameError)
}

func testPostConfigNameError(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	c := serverTestUserSpecificData.config

	timelines := c.GetTwitterTimelines()
	for i := range timelines {
		timelines[i].Name = ""
	}
	favorites := c.GetTwitterFavorites()
	for i := range favorites {
		favorites[i].Name = ""
	}
	searches := c.GetTwitterSearches()
	for i := range searches {
		searches[i].Name = ""
	}
	msgs := c.GetSlackMessages()
	for i := range msgs {
		msgs[i].Name = ""
	}

	require.NoError(t, page.Navigate(url))

	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "name_error_config_before_post.png")))
	wg.Add(1)
	require.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "name_error_config_after_post.png")))

	msg, err := page.FindByID("error-message").Text()
	require.NoError(t, err)
	require.True(t, strings.Contains(msg, "No name"))

	cProps := c.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	require.Nil(t, deep.Equal(&cProps, &configProps))
}

func TestPostConfigTagsInput(t *testing.T) {
	skipIfDialTimeout(t, "twitter.com", "https", 30*time.Second)
	testPostConfig(t, testPostConfigTagsInput)
}

func testPostConfigTagsInput(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	t.Skip("Skip because phantom.js doesn't support tagsinput currently.")

	require.NoError(t, page.Navigate(url))

	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "tags_input_config_before_post.png")))
	name := "twitter.timelines.screen_names"
	keys := "foo,bar"
	require.NoError(t, page.AllByName(name).SendKeys(keys))
	require.NoError(t, page.Screenshot(filepath.Join(screenshotsDir, "tags_input_config_after_post.png")))
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
	authMock.EXPECT().GetLoginUser(gomock.Any()).Times(2).Return(serverTestTwitterUser, nil)
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
	require.True(t, err == nil || strings.HasSuffix(err.Error(), expectedErrMsg))
	cur := length()
	require.Equal(t, prev+1, cur)
	testResponseIsRedirect(t, res, "/config")
}

func TestIndexPage(t *testing.T) {
	skipIfDialTimeout(t, "twitter.com", "https", 30*time.Second)

	twitterAPIMock := generateTwitterAPIMock(t, anaconda.User{ScreenName: "foo"}, nil)
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = core.NewTwitterAPI(twitterAPIMock, nil, nil)

	testIndex(t, testIndexPage)
}

func testIndexPage(url string) error {
	page, err := getDriver().NewPage()
	if err != nil {
		return err
	}

	err = page.Navigate(url)
	if err != nil {
		return err
	}

	err = page.Screenshot(filepath.Join(screenshotsDir, "index.png"))
	if err != nil {
		return err
	}
	return nil
}

func TestGetIndexWithoutTwitterAuthenticated(t *testing.T) {
	twitterAPIMock := generateTwitterAPIMock(t, anaconda.User{}, fmt.Errorf("your Twitter account is not authenticated"))
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = core.NewTwitterAPI(twitterAPIMock, nil, nil)

	testIndex(t, testGet)
}

func TestGetIndexWithTwitterAuthenticated(t *testing.T) {
	twitterAPIMock := generateTwitterAPIMock(t, anaconda.User{ScreenName: "foo"}, nil)
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = core.NewTwitterAPI(twitterAPIMock, nil, nil)

	testIndex(t, testGet)
}

func generateTwitterAPIMock(t *testing.T, user anaconda.User, userErr error) *mocks.MockTwitterAPI {
	ctrl := gomock.NewController(t)
	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().GetSelf(gomock.Any()).Return(user, userErr)
	return twitterAPIMock
}

func testIndex(t *testing.T, f func(url string) error) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().GetLoginUser(gomock.Any()).Times(2).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	tmpCache := serverTestUserSpecificData.cache
	defer func() { serverTestUserSpecificData.cache = tmpCache }()
	serverTestUserSpecificData.cache = data.NewTestFileCache("", t)
	img := models.ImageCacheData{}
	serverTestUserSpecificData.cache.SetImage(img)

	s := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer s.Close()

	err := f(s.URL)
	require.NoError(t, err)
}

func TestGetTwitterUserSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().GetLoginUser(gomock.Any()).Return(serverTestTwitterUser, nil)
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
	serverTestUserSpecificData.twitterAPI = core.NewTwitterAPI(twitterAPIMock, nil, nil)

	s := httptest.NewServer(http.HandlerFunc(twitterUserSearchHandler))
	defer s.Close()

	res, err := http.Get(s.URL)
	require.NoError(t, err)
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	us := []anaconda.User{}
	err = json.Unmarshal(bs, &us)
	require.NoError(t, err)

	require.Nil(t, deep.Equal(users, us))
}

func testResponseIsRedirect(t *testing.T, res *http.Response, locPrefix string) {
	require.Equal(t, http.StatusSeeOther, res.StatusCode)
	loc := res.Header.Get("Location")
	require.True(t, strings.HasPrefix(loc, locPrefix))
}

// TODO: Show skip warning even when executing `go test` without `-v` argument
func skipIfDialTimeout(t *testing.T, domain, protocol string, timeout time.Duration) {
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", domain, protocol), timeout)
	if err != nil {
		t.Skip("Skip because network is unavailable: ", err)
	}
}
