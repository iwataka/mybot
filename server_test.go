package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/deep"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/oauth"
	"github.com/iwataka/mybot/worker"
	"github.com/markbates/goth"
	"github.com/sclevine/agouti"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	serverTestTwitterUserID = "123456"
	screenshotsDir          = "screenshots"
)

var (
	serverTestUserSpecificData *userSpecificData
	serverTestTwitterUser      = goth.User{Name: "foo", NickName: "bar", UserID: serverTestTwitterUserID}
	driver                     *agouti.WebDriver
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	if driver != nil {
		driver.Stop()
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
	serverTestUserSpecificData.config, err = mybot.NewFileConfig("lib/testdata/config.template.toml")
	if err != nil {
		panic(err)
	}
	serverTestUserSpecificData.statuses = map[int]bool{}
	serverTestUserSpecificData.statuses = initialStatuses()
	serverTestUserSpecificData.workerChans = map[int]chan *worker.WorkerSignal{}
	serverTestUserSpecificData.slackAPI = mybot.NewSlackAPIWithAuth("", serverTestUserSpecificData.config, nil)
	userSpecificDataMap[twitterUserIDPrefix+serverTestTwitterUserID] = serverTestUserSpecificData

	if _, err := os.Stat(screenshotsDir); err != nil {
		err := os.Mkdir(screenshotsDir, os.FileMode(0755))
		if err != nil {
			fmt.Printf("Failed to make `%s` directory\n", screenshotsDir)
			os.Exit(1)
		}
	}

	deep.IgnoreDifferenceBetweenEmptyMapAndNil = true
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = true
}

func TestTwitterColsPage(t *testing.T) {
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

func TestGetTwitterColsIfAssetsNotExist(t *testing.T) {
	testIfAssetsNotExist(t, TestGetTwitterCols)
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
	serverTestUserSpecificData.twitterAPI = mybot.NewTwitterAPI(twitterAPIMock, nil, nil)

	s := httptest.NewServer(http.HandlerFunc(twitterColsHandler))
	defer s.Close()

	err := f(s.URL)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
}

func TestGetConfigIfAssetsNotExist(t *testing.T) {
	testIfAssetsNotExist(t, TestGetConfig)
}

func TestGetSetupTwitter(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(getSetup))
	defer s.Close()

	tmpTwitterApp := twitterApp
	var err error
	twitterApp, err = oauth.NewFileTwitterOAuthApp("")
	assert.NoError(t, err)
	defer func() { twitterApp = tmpTwitterApp }()

	tmpSlackApp := slackApp
	slackApp, err = oauth.NewFileOAuthApp("")
	assert.NoError(t, err)
	defer func() { slackApp = tmpSlackApp }()

	err = testGet(s.URL)
	assert.NoError(t, err)
}

func TestGetSetupTwitterIfAssetsNotExist(t *testing.T) {
	testIfAssetsNotExist(t, TestGetSetupTwitter)
}

func TestGetConfigJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	authMock := mocks.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CompleteUserAuth(gomock.Any(), gomock.Any(), gomock.Any()).Return(serverTestTwitterUser, nil)
	tmpAuth := authenticator
	defer func() { authenticator = tmpAuth }()
	authenticator = authMock

	s := httptest.NewServer(http.HandlerFunc(configJSONHandler))
	defer s.Close()

	res, err := http.Get(s.URL)
	assert.NoError(t, err)

	err = checkHTTPResponse(res)
	assert.NoError(t, err)

	bs, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	cfg, err := mybot.NewFileConfig("")
	assert.NoError(t, err)

	err = cfg.Unmarshal(".json", bs)
	assert.NoError(t, err)

	cfgProps := cfg.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	assert.Nil(t, deep.Equal(cfgProps, configProps))
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
	assert.NoError(t, err)

	err = checkHTTPResponse(res)
	assert.NoError(t, err)

	hasForceDownload := strings.Contains(res.Header.Get("Content-Type"), "application/force-download")
	assert.True(t, hasForceDownload)

	hasContentDisposition := strings.Contains(res.Header.Get("Content-Disposition"), ".json")
	assert.True(t, hasContentDisposition)

	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	cfg, err := mybot.NewFileConfig("")
	assert.NoError(t, err)

	err = json.Unmarshal(bs, cfg)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	checkHTTPResponse(res)
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

	page, err := getDriver().NewPage()
	assert.NoError(t, err)

	curUserData := serverTestUserSpecificData.config
	defer func() { serverTestUserSpecificData.config = curUserData }()
	f(t, s.URL, page, wg)
}

func TestPostConfigWithoutModification(t *testing.T) {
	testPostConfig(t, testPostConfigWithoutModification)
}

func testPostConfigWithoutModification(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	c := serverTestUserSpecificData.config

	assert.NoError(t, page.Navigate(url))

	page.Screenshot(filepath.Join(screenshotsDir, "config_before_post_without_modification.png"))
	wg.Add(1)
	assert.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	page.Screenshot(filepath.Join(screenshotsDir, "config_after_post_without_modification.png"))

	cProps := c.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	deep.IgnoreDifferenceBetweenEmptyMapAndNil = true
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = true
	assert.Nil(t, deep.Equal(cProps, configProps))
}

func TestPostConfigDelete(t *testing.T) {
	testPostConfig(t, testPostConfigDelete)
}

func testPostConfigDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	assert.NoError(t, page.Navigate(url))

	page.Screenshot(filepath.Join(screenshotsDir, "delete_config_before_post.png"))
	assert.NoError(t, page.AllByClass("config-row-delete").Click())
	page.Screenshot(filepath.Join(screenshotsDir, "delete_config_after_click_delete_buttons.png"))
	wg.Add(1)
	assert.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	page.Screenshot(filepath.Join(screenshotsDir, "delete_config_after_post.png"))

	assert.Empty(t, serverTestUserSpecificData.config.GetTwitterTimelines())
	assert.Empty(t, serverTestUserSpecificData.config.GetTwitterFavorites())
	assert.Empty(t, serverTestUserSpecificData.config.GetTwitterSearches())
	assert.Empty(t, serverTestUserSpecificData.config.GetSlackMessages())
}

func TestPostConfigSingleDelete(t *testing.T) {
	testPostConfig(t, testPostConfigSingleDelete)
}

func testPostConfigSingleDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	c := serverTestUserSpecificData.config

	assert.NoError(t, page.Navigate(url))

	page.Screenshot(filepath.Join(screenshotsDir, "single_delete_config_before_post.png"))
	assert.NoError(t, page.AllByClass("config-row-delete").At(0).Click())
	wg.Add(1)
	assert.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	page.Screenshot(filepath.Join(screenshotsDir, "single_delete_config_after_post.png"))

	assert.Equal(t, len(serverTestUserSpecificData.config.GetTwitterTimelines()), len(c.GetTwitterTimelines())-1)
}

func TestPostConfigDoubleDelete(t *testing.T) {
	testPostConfig(t, testPostConfigDoubleDelete)
}

func testPostConfigDoubleDelete(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	c := serverTestUserSpecificData.config

	assert.NoError(t, page.Navigate(url))

	page.Screenshot(filepath.Join(screenshotsDir, "double_delete_config_before_post.png"))
	assert.NoError(t, page.AllByClass("config-row-delete").DoubleClick())
	wg.Add(1)
	assert.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	page.Screenshot(filepath.Join(screenshotsDir, "double_delete_config_after_post.png"))

	cProps := c.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	assert.Nil(t, deep.Equal(cProps.Slack, configProps.Slack))
	assert.Nil(t, deep.Equal(cProps.Twitter, configProps.Twitter))
}

func TestPostConfigNameError(t *testing.T) {
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
	for i, _ := range timelines {
		timelines[i].Name = ""
	}
	favorites := c.GetTwitterFavorites()
	for i, _ := range favorites {
		favorites[i].Name = ""
	}
	searches := c.GetTwitterSearches()
	for i, _ := range searches {
		searches[i].Name = ""
	}
	msgs := c.GetSlackMessages()
	for i, _ := range msgs {
		msgs[i].Name = ""
	}

	assert.NoError(t, page.Navigate(url))

	page.Screenshot(filepath.Join(screenshotsDir, "name_error_config_before_post.png"))
	wg.Add(1)
	assert.NoError(t, page.FindByID("overwrite").Submit())
	wg.Wait()
	page.Screenshot(filepath.Join(screenshotsDir, "name_error_config_after_post.png"))

	msg, err := page.FindByID("error-message").Text()
	assert.NoError(t, err)
	assert.True(t, strings.Contains(msg, "No name"))

	cProps := c.GetProperties()
	configProps := serverTestUserSpecificData.config.GetProperties()
	assert.Nil(t, deep.Equal(cProps, configProps))
}

func TestPostConfigTagsInput(t *testing.T) {
	testPostConfig(t, testPostConfigTagsInput)
}

func testPostConfigTagsInput(
	t *testing.T,
	url string,
	page *agouti.Page,
	wg *sync.WaitGroup,
) {
	// _, err := net.DialTimeout("tcp", "cdnjs.cloudflare.com:https", 30*time.Second)
	// if err != nil {
	// 	t.Skip("Skip because network is unavailable: ", err)
	// }
	t.Skip("Skip because phantom.js doesn't support tagsinput currently.")

	assert.NoError(t, page.Navigate(url))

	page.Screenshot(filepath.Join(screenshotsDir, "tags_input_config_before_post.png"))
	name := "twitter.timelines.screen_names"
	keys := "foo,bar"
	assert.NoError(t, page.AllByName(name).SendKeys(keys))
	page.Screenshot(filepath.Join(screenshotsDir, "tags_input_config_after_post.png"))
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
	assert.True(t, err == nil || strings.HasSuffix(err.Error(), expectedErrMsg))
	cur := length()
	assert.Equal(t, prev+1, cur)
	testResponseIsRedirect(t, res, "/config")
}

func TestIndexPage(t *testing.T) {
	twitterAPIMock := generateTwitterAPIMock(t, anaconda.User{ScreenName: "foo"}, nil)
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = mybot.NewTwitterAPI(twitterAPIMock, nil, nil)

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
	twitterAPIMock := generateTwitterAPIMock(t, anaconda.User{}, fmt.Errorf("Your Twitter account is not authenticated."))
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = mybot.NewTwitterAPI(twitterAPIMock, nil, nil)

	testIndex(t, testGet)
}

func TestGetIndexWithTwitterAuthenticated(t *testing.T) {
	twitterAPIMock := generateTwitterAPIMock(t, anaconda.User{ScreenName: "foo"}, nil)
	tmpTwitterAPI := serverTestUserSpecificData.twitterAPI
	defer func() { serverTestUserSpecificData.twitterAPI = tmpTwitterAPI }()
	serverTestUserSpecificData.twitterAPI = mybot.NewTwitterAPI(twitterAPIMock, nil, nil)

	testIndex(t, testGet)
}

func generateTwitterAPIMock(t *testing.T, user anaconda.User, userErr error) *mocks.MockTwitterAPI {
	ctrl := gomock.NewController(t)
	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().GetSelf(gomock.Any()).Return(user, userErr)
	return twitterAPIMock
}

func TestGetIndexIfAssetsNotExist(t *testing.T) {
	testIfAssetsNotExist(t, TestGetIndexWithTwitterAuthenticated)
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
	serverTestUserSpecificData.cache = data.NewTestFileCache("", t)
	img := models.ImageCacheData{}
	serverTestUserSpecificData.cache.SetImage(img)

	s := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer s.Close()

	err := f(s.URL)
	assert.NoError(t, err)
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
	serverTestUserSpecificData.twitterAPI = mybot.NewTwitterAPI(twitterAPIMock, nil, nil)

	s := httptest.NewServer(http.HandlerFunc(twitterUserSearchHandler))
	defer s.Close()

	res, err := http.Get(s.URL)
	assert.NoError(t, err)
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	us := []anaconda.User{}
	err = json.Unmarshal(bs, &us)
	assert.NoError(t, err)

	assert.Nil(t, deep.Equal(users, us))
}

func testResponseIsRedirect(t *testing.T, res *http.Response, locPrefix string) {
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)
	loc := res.Header.Get("Location")
	assert.True(t, strings.HasPrefix(loc, locPrefix))
}

func testIfAssetsNotExist(t *testing.T, f func(t *testing.T)) {
	tmpdir := "tmp"
	require.NoError(t, os.Mkdir(tmpdir, os.FileMode(0777)))
	defer os.Remove(tmpdir)

	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmpdir))
	defer os.Chdir(wd)

	f(t)
}
