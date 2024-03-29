package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iwataka/mybot/core"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/tmpl"
	"github.com/iwataka/mybot/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/twitter"
	"github.com/mattn/go-zglob"
)

const (
	assetsDir           = "assets"
	appUserIDFormat     = "%s-%s"
	sessionKeyForUser   = "mybot_user"
	sessionName         = "_mybot_sess"
	trueValue           = "true"
	defaultConfigFormat = ".json"
)

var (
	htmlTemplateDir                      = filepath.Join(assetsDir, "tmpl")
	authenticator   models.Authenticator = &Authenticator{}
)

var templateFuncMap = template.FuncMap{
	"checkbox":              tmpl.Checkbox,
	"boolSelectbox":         tmpl.BoolSelectbox,
	"selectbox":             tmpl.Selectbox,
	"listTextbox":           tmpl.ListTextbox,
	"textboxOfFloat64Ptr":   tmpl.TextboxOfFloat64Ptr,
	"textboxOfIntPtr":       tmpl.TextboxOfIntPtr,
	"likelihoodMultiSelect": tmpl.LikelihoodMultiSelect,
	"newMap":                tmpl.NewMap,
	"add":                   func(i1, i2 int) int { return i1 + i2 },
	"replace":               func(s, old, new string) string { return strings.Replace(s, old, new, -1) },
	"title":                 func(s string) string { return strings.Title(s) },
}

// Authenticator is an implementation of models.Authenticator and provides some
// common functions for authenticating users.
type Authenticator struct{}

// SetProvider sets a specified provider name to the gothic module.
func (a *Authenticator) SetProvider(name string, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", name)
	r.URL.RawQuery = q.Encode()
}

// InitProvider initializes a provider and makes it to be used.
func (a *Authenticator) InitProvider(provider, callback, ck, cs string) {
	var p goth.Provider
	switch provider {
	case "twitter":
		p = twitter.New(
			ck,
			cs,
			callback,
		)
	case "slack":
		p = slack.New(
			ck,
			cs,
			callback,
			"client",
		)
	}
	if p != nil {
		goth.UseProviders(p)
	}
}

// CompleteUserAuth executes user authentication and returns the user
// information.
func (a *Authenticator) CompleteUserAuth(provider string, w http.ResponseWriter, r *http.Request) (user goth.User, err error) {
	a.SetProvider(provider, r)
	q := r.URL.Query()
	q.Add("state", "state")
	r.URL.RawQuery = q.Encode()
	user, err = gothic.CompleteUserAuth(w, r)
	return
}

func (a *Authenticator) Login(user goth.User, w http.ResponseWriter, r *http.Request) error {
	sess, err := serverSession.Get(r, sessionName)
	if err != nil {
		return utils.WithStack(err)
	}
	user.RawData = nil // RawData cannot be converted into session data currently
	sess.Values[sessionKeyForUser] = user
	err = sess.Save(r, w)
	if err != nil {
		return utils.WithStack(err)
	}
	return nil
}

func (a *Authenticator) GetLoginUser(r *http.Request) (goth.User, error) {
	sess, err := serverSession.Get(r, sessionName)
	if err != nil {
		return goth.User{}, utils.WithStack(err)
	}
	val, exists := sess.Values[sessionKeyForUser]
	if exists {
		if user, ok := val.(goth.User); ok {
			return user, nil
		}
	}
	return goth.User{}, fmt.Errorf("No login user")
}

// Logout executes logout operation of the current login-user.
func (a *Authenticator) Logout(w http.ResponseWriter, r *http.Request) error {
	sess, err := serverSession.Get(r, sessionName)
	if err != nil {
		return utils.WithStack(err)
	}
	val, exists := sess.Values[sessionKeyForUser]
	if exists {
		if user, ok := val.(goth.User); ok {
			sess.Options.MaxAge = -1
			err = sess.Save(r, w)
			if err != nil {
				return utils.WithStack(err)
			}
			a.SetProvider(user.Provider, r)
			return nil
		}
	}
	return fmt.Errorf("No login user")
}

func wrapHandler(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ck, cs := twitterApp.GetCreds(); ck == "" || cs == "" {
			c.Redirect(http.StatusSeeOther, "/setup/")
			return
		}

		if ck, cs := slackApp.GetCreds(); ck == "" || cs == "" {
			c.Redirect(http.StatusSeeOther, "/setup/")
			return
		}

		if _, err := authenticator.GetLoginUser(c.Request); err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		f(c)
	}
}

func setupRouter() *gin.Engine {
	return setupRouterWithWrapper(wrapHandler)
}

func setupRouterWithWrapper(wrapper func(gin.HandlerFunc) gin.HandlerFunc) *gin.Engine {
	r := setupBaseRouter()
	r.GET("/", wrapper(getIndexHandler))
	r.GET("/account/delete", wrapper(accountDeleteHandler)) // currently hidden endpoint
	r.GET("/twitter-collections/", wrapper(twitterColsHandler))
	r.Any("/config/", wrapper(configHandler))
	r.Any("/config/file/", wrapper(configFileHandler))
	r.POST("/config/timelines/add", wrapper(configTimelineAddHandler))
	r.POST("/config/favorites/add", wrapper(configFavoriteAddHandler))
	r.POST("/config/searches/add", wrapper(configSearchAddHandler))
	r.POST("/config/messages/add", wrapper(configMessageAddHandler))
	r.GET("/auth/:provider", authHandler)
	r.GET("/auth/callback/:provider", authCallbackHandler)
	r.GET("/login/", loginHandler)
	r.Any("/setup/", setupHandler)
	r.GET("/logout/", logoutHandler)
	r.GET("/twitter/users/search/", wrapper(twitterUserSearchHandler)) // For Twitter user auto-completion usage
	return r
}

func setupBaseRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Static("/assets", "./assets")
	r.Delims("{{{", "}}}")
	r.SetFuncMap(templateFuncMap)
	tmplFiles, err := zglob.Glob(filepath.Join(htmlTemplateDir, "**", "*.tmpl"))
	if err != nil {
		panic(err)
	}
	r.LoadHTMLFiles(tmplFiles...)
	return r
}

func accountDeleteHandler(c *gin.Context) {
	w, r := c.Writer, c.Request
	getAccountDelete(w, r)
}

func getAccountDelete(w http.ResponseWriter, r *http.Request) {
	user, err := authenticator.GetLoginUser(r)
	if err != nil {
		panic(err)
	}
	userID := fmt.Sprintf(appUserIDFormat, user.Provider, user.UserID)
	data := userSpecificDataMap[userID]
	err = data.delete()
	if err != nil {
		panic(err)
	}
	delete(userSpecificDataMap, userID)
	err = authenticator.Logout(w, r)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/login/", http.StatusSeeOther)
}

func getIndexHandler(c *gin.Context) {
	user, err := authenticator.GetLoginUser(c.Request)
	if err != nil {
		panic(err)
	}
	data := userSpecificDataMap[fmt.Sprintf(appUserIDFormat, user.Provider, user.UserID)]
	getIndex(c, data.cache, data.slackAPI, data.twitterAPI, data.statuses())
}

func getIndex(
	c *gin.Context,
	cache data.Cache,
	slackAPI *core.SlackAPI,
	twitterAPI *core.TwitterAPI,
	statuses map[int]bool,
) {
	_, twitterScreenName := getTwitterInfo(twitterAPI)
	slackTeam, slackURL := getSlackInfo(slackAPI)
	imgURL, imgSrc, imgAnalysisResult, imgAnalysisDate := imageAnalysis(cache)

	data := gin.H{
		"NavbarName":               "",
		"TwitterName":              twitterScreenName,
		"SlackTeam":                slackTeam,
		"SlackURL":                 slackURL,
		"GoogleEnabled":            googleEnabled(),
		"ImageURL":                 imgURL,
		"ImageSource":              imgSrc,
		"ImageAnalysisResult":      imgAnalysisResult,
		"ImageAnalysisDate":        imgAnalysisDate,
		"TwitterListenDMStatus":    statuses[twitterDMRoutineKey],
		"TwitterListenUsersStatus": statuses[twitterUserRoutineKey],
		"TwitterPeriodicStatus":    statuses[twitterPeriodicRoutineKey],
		"SlackListenerStatus":      statuses[slackRoutineKey],
	}

	c.HTML(http.StatusOK, "index", data)
}

func imageAnalysis(cache data.Cache) (string, string, string, string) {
	imageSource := ""
	imageURL := ""
	imageAnalysisResult := ""
	imageAnalysisDate := ""
	images := cache.GetLatestImages(1)
	if len(images) != 0 {
		imgCache := images[0]
		imageSource = imgCache.Src
		imageURL = imgCache.URL
		buf := new(bytes.Buffer)
		result := imgCache.AnalysisResult
		err := json.Indent(buf, []byte(result), "", "  ")
		if err == nil {
			imageAnalysisResult = buf.String()
		}
		imageAnalysisDate = imgCache.AnalysisDate
	}
	return imageURL, imageSource, imageAnalysisResult, imageAnalysisDate
}

func twitterColsHandler(c *gin.Context) {
	r := c.Request
	twitterUser, err := authenticator.GetLoginUser(r)
	if err != nil {
		panic(err)
	}
	data := userSpecificDataMap[fmt.Sprintf(appUserIDFormat, twitterUser.Provider, twitterUser.UserID)]
	getTwitterCols(c, data.slackAPI, data.twitterAPI)
}

func getTwitterCols(c *gin.Context, slackAPI *core.SlackAPI, twitterAPI *core.TwitterAPI) {
	twitterID, twitterScreenName := getTwitterInfo(twitterAPI)
	slackTeam, slackURL := getSlackInfo(slackAPI)

	colMap := make(map[string]string)
	activeCol := ""
	id, err := strconv.ParseInt(twitterID, 10, 64)
	if err != nil {
		panic(err)
	}
	colList, err := twitterAPI.BaseAPI().GetCollectionListByUserId(id, nil)
	if err == nil && len(colList.Objects.Timelines) != 0 {
		names := []string{}
		for _, c := range colList.Objects.Timelines {
			name := strings.Replace(c.Name, " ", "-", -1)
			colMap[name] = c.CollectionUrl
			names = append(names, name)
		}
		sort.Strings(names)
		activeCol = names[0]
	}

	data := gin.H{
		"NavbarName":       "TwitterCols",
		"TwitterName":      twitterScreenName,
		"SlackTeam":        slackTeam,
		"SlackURL":         slackURL,
		"GoogleEnabled":    googleEnabled(),
		"CollectionMap":    colMap,
		"ActiveCollection": activeCol,
	}
	c.HTML(http.StatusOK, "twitterCols", data)
}

type checkboxCounter struct {
	name       string
	extraCount int
}

func newCheckboxCounter(name string) checkboxCounter {
	return checkboxCounter{name, 0}
}

func (c *checkboxCounter) returnValue(index int, val map[string][]string, def bool) bool {
	vs := val[c.name]
	if len(vs) <= index {
		return def
	}
	if val[c.name][index+c.extraCount] == trueValue {
		c.extraCount++
		return true
	}
	return false
}

func configHandler(c *gin.Context) {
	twitterUser, err := authenticator.GetLoginUser(c.Request)
	if err != nil {
		panic(err)
	}
	data := userSpecificDataMap[fmt.Sprintf(appUserIDFormat, twitterUser.Provider, twitterUser.UserID)]

	switch c.Request.Method {
	case http.MethodPost:
		postConfig(c.Writer, c.Request, data.config, twitterUser)
	case http.MethodGet:
		getConfig(c, data.config, data.slackAPI, data.twitterAPI)
	default:
		http.NotFound(c.Writer, c.Request)
	}
}

func postConfig(w http.ResponseWriter, r *http.Request, config core.Config, twitterUser goth.User) {
	var err error

	defer func() {
		if err == nil {
			err = config.Save()
			if err == nil {
				userID := fmt.Sprintf(appUserIDFormat, twitterUser.Provider, twitterUser.UserID)
				go userSpecificDataMap[userID].restart()
			}
		} else {
			// TODO: add load error to cookie if any
			_ = config.Load()
		}

		if err != nil {
			msgCookie := &http.Cookie{
				Name:  "mybot.config.message",
				Value: err.Error(),
				Path:  "/config/",
			}
			http.SetCookie(w, msgCookie)
		}
		http.Redirect(w, r, "/config/", http.StatusSeeOther)
	}()

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		return
	}

	err = setTimelinesToConfig(config, r)
	if err != nil {
		return
	}

	err = setFavoritesToConfig(config, r)
	if err != nil {
		return
	}

	err = setSearchesToConfig(config, r)
	if err != nil {
		return
	}

	err = setMessagesToConfig(config, r)
	if err != nil {
		return
	}

	config.SetPollingDuration(r.Form["duration"][0])

	err = config.Validate()
}

func setTimelinesToConfig(config core.Config, r *http.Request) error {
	val := r.MultipartForm.Value

	prefix := "twitter.timelines"
	deletedFlags := val[prefix+".deleted"]
	timelines := []core.TimelineConfig{}
	actions, err := postConfigForActions(val, prefix, deletedFlags)
	excludeRepliesCounter := newCheckboxCounter(prefix + ".exclude_replies")
	includeRtsCounter := newCheckboxCounter(prefix + ".include_rts")
	if err != nil {
		return err
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == trueValue {
			continue
		}
		timeline := core.NewTimelineConfig()
		timeline.Name = val[prefix+".name"][i]
		timeline.ScreenNames = tmpl.GetListTextboxValue(val, i, prefix+".screen_names")
		timeline.ExcludeReplies = excludeRepliesCounter.returnValue(i, val, false)
		timeline.IncludeRts = includeRtsCounter.returnValue(i, val, false)
		if timeline.Count, err = tmpl.GetIntPtr(val, i, prefix+".count"); err != nil {
			return err
		}
		if timeline.Filter, err = postConfigForFilter(r, i, prefix); err != nil {
			return err
		}
		timeline.Action = actions[i]
		timelines = append(timelines, timeline)
	}
	config.SetTwitterTimelines(timelines)
	return nil
}

func setFavoritesToConfig(config core.Config, r *http.Request) error {
	val := r.MultipartForm.Value

	prefix := "twitter.favorites"
	deletedFlags := val[prefix+".deleted"]
	favorites := []core.FavoriteConfig{}
	actions, err := postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return err
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == trueValue {
			continue
		}
		favorite := core.NewFavoriteConfig()
		favorite.Name = val[prefix+".name"][i]
		favorite.ScreenNames = tmpl.GetListTextboxValue(val, i, prefix+".screen_names")
		if favorite.Count, err = tmpl.GetIntPtr(val, i, prefix+".count"); err != nil {
			return err
		}
		if favorite.Filter, err = postConfigForFilter(r, i, prefix); err != nil {
			return err
		}
		favorite.Action = actions[i]
		favorites = append(favorites, favorite)
	}
	config.SetTwitterFavorites(favorites)
	return nil
}

func setSearchesToConfig(config core.Config, r *http.Request) error {
	val := r.MultipartForm.Value

	prefix := "twitter.searches"
	deletedFlags := val[prefix+".deleted"]
	searches := []core.SearchConfig{}
	actions, err := postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return err
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == trueValue {
			continue
		}
		search := core.NewSearchConfig()
		search.Name = val[prefix+".name"][i]
		search.Queries = tmpl.GetListTextboxValue(val, i, prefix+".queries")
		search.ResultType = val[prefix+".result_type"][i]
		if search.Count, err = tmpl.GetIntPtr(val, i, prefix+".count"); err != nil {
			return err
		}
		if search.Filter, err = postConfigForFilter(r, i, prefix); err != nil {
			return err
		}
		search.Action = actions[i]
		searches = append(searches, search)
	}
	config.SetTwitterSearches(searches)
	return nil
}

func setMessagesToConfig(config core.Config, r *http.Request) error {
	val := r.MultipartForm.Value

	prefix := "slack.messages"
	deletedFlags := val[prefix+".deleted"]
	msgs := []core.MessageConfig{}
	actions, err := postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return err
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == trueValue {
			continue
		}
		msg := core.NewMessageConfig()
		msg.Name = val[prefix+".name"][i]
		msg.Channels = tmpl.GetListTextboxValue(val, i, prefix+".channels")
		if msg.Filter, err = postConfigForFilter(r, i, prefix); err != nil {
			return err
		}
		msg.Action = actions[i]
		msgs = append(msgs, msg)
	}
	config.SetSlackMessages(msgs)
	return nil
}

func postConfigForFilter(r *http.Request, i int, prefix string) (core.Filter, error) {
	val := r.MultipartForm.Value

	prefix = prefix + ".filter."
	filter := core.NewFilter()
	filter.Patterns = tmpl.GetListTextboxValue(val, i, prefix+"patterns")
	filter.URLPatterns = tmpl.GetListTextboxValue(val, i, prefix+"url_patterns")
	filter.HasMedia = tmpl.GetBoolSelectboxValue(val, i, prefix+"has_media")
	fThreshold, err := tmpl.GetIntPtr(val, i, prefix+"favorite_threshold")
	if err != nil {
		return core.NewFilter(), utils.WithStack(err)
	}
	filter.FavoriteThreshold = fThreshold
	rThreshold, err := tmpl.GetIntPtr(val, i, prefix+"retweeted_threshold")
	if err != nil {
		return core.NewFilter(), utils.WithStack(err)
	}
	filter.RetweetedThreshold = rThreshold
	filter.Lang = tmpl.GetString(val, prefix+"lang", i, "")
	filter.Vision.Label = tmpl.GetListTextboxValue(val, i, prefix+"vision.label")
	filter.Vision.Face.AngerLikelihood = tmpl.GetLikelihood(r, prefix+"vision.face.anger_likelihood", i, "")
	filter.Vision.Face.BlurredLikelihood = tmpl.GetLikelihood(r, prefix+"vision.face.blurred_likelihood", i, "")
	filter.Vision.Face.HeadwearLikelihood = tmpl.GetLikelihood(r, prefix+"vision.face.headwear_likelihood", i, "")
	filter.Vision.Face.JoyLikelihood = tmpl.GetLikelihood(r, prefix+"vision.face.joy_likelihood", i, "")
	filter.Vision.Text = tmpl.GetListTextboxValue(val, i, prefix+"vision.text")
	filter.Vision.Landmark = tmpl.GetListTextboxValue(val, i, prefix+"vision.landmark")
	filter.Vision.Logo = tmpl.GetListTextboxValue(val, i, prefix+"vision.logo")
	minSentiment, err := tmpl.GetFloat64Ptr(val, i, prefix+"language.min_sentiment")
	if err != nil {
		return core.NewFilter(), utils.WithStack(err)
	}
	filter.Language.MinSentiment = minSentiment
	maxSentiment, err := tmpl.GetFloat64Ptr(val, i, prefix+"language.max_sentiment")
	if err != nil {
		return core.NewFilter(), utils.WithStack(err)
	}
	filter.Language.MaxSentiment = maxSentiment
	return filter, nil
}

func postConfigForActions(
	val map[string][]string,
	prefix string,
	deletedFlags []string,
) ([]data.Action, error) {
	prefix = prefix + ".action."
	tweetCounter := newCheckboxCounter(prefix + "twitter.tweet")
	retweetCounter := newCheckboxCounter(prefix + "twitter.retweet")
	favoriteCounter := newCheckboxCounter(prefix + "twitter.favorite")
	pinCounter := newCheckboxCounter(prefix + "slack.pin")
	starCounter := newCheckboxCounter(prefix + "slack.star")
	results := []data.Action{}
	for i := 0; i < len(deletedFlags); i++ {
		a, err := postConfigForAction(val, i, prefix)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		a.Twitter.Tweet = tweetCounter.returnValue(i, val, false)
		a.Twitter.Retweet = retweetCounter.returnValue(i, val, false)
		a.Twitter.Favorite = favoriteCounter.returnValue(i, val, false)
		a.Slack.Pin = pinCounter.returnValue(i, val, false)
		a.Slack.Star = starCounter.returnValue(i, val, false)
		results = append(results, a)
	}
	return results, nil
}

func postConfigForAction(val map[string][]string, i int, prefix string) (data.Action, error) {
	action := data.NewAction()
	action.Twitter.Collections = tmpl.GetListTextboxValue(val, i, prefix+"twitter.collections")
	action.Slack.Channels = tmpl.GetListTextboxValue(val, i, prefix+"slack.channels")
	action.Slack.Reactions = tmpl.GetListTextboxValue(val, i, prefix+"slack.reactions")
	return action, nil
}

func getConfig(c *gin.Context, config core.Config, slackAPI *core.SlackAPI, twitterAPI *core.TwitterAPI) {
	msg := ""
	msgCookie, err := c.Request.Cookie("mybot.config.message")
	if err == nil {
		msg = msgCookie.Value
	}

	if msgCookie != nil {
		msgCookie.Value = ""
		msgCookie.Path = "/config/"
		http.SetCookie(c.Writer, msgCookie)
	}

	_, twitterName := getTwitterInfo(twitterAPI)
	slackTeam, slackURL := getSlackInfo(slackAPI)
	data := gin.H{
		"NavbarName":    "Config",
		"TwitterName":   twitterName,
		"SlackTeam":     slackTeam,
		"SlackURL":      slackURL,
		"GoogleEnabled": googleEnabled(),
		"Message":       msg,
		"Config":        config.GetProperties(),
	}
	c.HTML(http.StatusOK, "config", data)
}

func configTimelineAddHandler(c *gin.Context) {
	w, r := c.Writer, c.Request
	twitterUser, err := authenticator.GetLoginUser(r)
	if err != nil {
		panic(err)
	}
	data := userSpecificDataMap[fmt.Sprintf(appUserIDFormat, twitterUser.Provider, twitterUser.UserID)]
	postConfigTimelineAdd(w, r, data.config)
}

func postConfigTimelineAdd(w http.ResponseWriter, r *http.Request, config core.Config) {
	addTimelineConfig(config)
	http.Redirect(w, r, "/config/", http.StatusSeeOther)
}

func addTimelineConfig(config core.Config) {
	config.AddTwitterTimeline(core.NewTimelineConfig())
}

func configFavoriteAddHandler(c *gin.Context) {
	w, r := c.Writer, c.Request
	twitterUser, err := authenticator.GetLoginUser(r)
	if err != nil {
		panic(err)
	}
	data := userSpecificDataMap[fmt.Sprintf(appUserIDFormat, twitterUser.Provider, twitterUser.UserID)]
	postConfigFavoriteAdd(w, r, data.config)
}

func postConfigFavoriteAdd(w http.ResponseWriter, r *http.Request, config core.Config) {
	addFavoriteConfig(config)
	http.Redirect(w, r, "/config/", http.StatusSeeOther)
}

func addFavoriteConfig(config core.Config) {
	config.AddTwitterFavorite(core.NewFavoriteConfig())
}

func configSearchAddHandler(c *gin.Context) {
	w, r := c.Writer, c.Request
	twitterUser, err := authenticator.GetLoginUser(r)
	if err != nil {
		panic(err)
	}
	data := userSpecificDataMap[fmt.Sprintf(appUserIDFormat, twitterUser.Provider, twitterUser.UserID)]
	postConfigSearchAdd(w, r, data.config)
}

func postConfigSearchAdd(w http.ResponseWriter, r *http.Request, config core.Config) {
	addSearchConfig(config)
	http.Redirect(w, r, "/config/", http.StatusSeeOther)
}

func addSearchConfig(config core.Config) {
	config.AddTwitterSearch(core.NewSearchConfig())
}

func configMessageAddHandler(c *gin.Context) {
	w, r := c.Writer, c.Request
	twitterUser, err := authenticator.GetLoginUser(r)
	if err != nil {
		panic(err)
	}
	data := userSpecificDataMap[fmt.Sprintf(appUserIDFormat, twitterUser.Provider, twitterUser.UserID)]
	postConfigMessageAdd(w, r, data.config)
}

func postConfigMessageAdd(w http.ResponseWriter, r *http.Request, c core.Config) {
	addMessageConfig(c)
	http.Redirect(w, r, "/config/", http.StatusSeeOther)
}

func addMessageConfig(config core.Config) {
	config.AddSlackMessage(core.NewMessageConfig())
}

func configFileHandler(c *gin.Context) {
	w, r := c.Writer, c.Request
	twitterUser, err := authenticator.GetLoginUser(r)
	if err != nil {
		panic(err)
	}
	data := userSpecificDataMap[fmt.Sprintf(appUserIDFormat, twitterUser.Provider, twitterUser.UserID)]

	switch r.Method {
	case http.MethodPost:
		postConfigFile(w, r, data.config)
	case http.MethodGet:
		getConfigFile(w, data.config)
	default:
		http.NotFound(w, r)
	}
}

func postConfigFile(w http.ResponseWriter, r *http.Request, config core.Config) {
	msg := ""
	defer func() {
		if len(msg) != 0 {
			msgCookie := &http.Cookie{
				Name:  "mybot.config.message",
				Value: msg,
				Path:  "/config/",
			}
			http.SetCookie(w, msgCookie)
		}
		http.Redirect(w, r, "/config/", http.StatusSeeOther)
	}()

	file, _, err := r.FormFile("mybot.config")
	if err != nil {
		msg = err.Error()
		return
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		msg = err.Error()
		return
	}
	err = config.Unmarshal(defaultConfigFormat, bytes)
	if err != nil {
		msg = err.Error()
		return
	}
	err = config.Validate()
	if err != nil {
		msg = err.Error()
		err = config.Load()
		if err != nil {
			msg = err.Error()
		}
		return
	}
	err = config.Save()
	if err != nil {
		msg = err.Error()
		return
	}
}

func getConfigFile(w http.ResponseWriter, config core.Config) {
	ext := defaultConfigFormat
	w.Header().Add("Content-Type", "application/force-download; charset=utf-8")
	w.Header().Add("Content-Disposition", `attachment; filename="config`+ext+`"`)
	bytes, err := config.Marshal(ext)
	if err != nil {
		panic(err)
	}
	len, err := w.Write(bytes)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len), 16))
}

func setupHandler(c *gin.Context) {
	twitterCk, twitterCs := twitterApp.GetCreds()
	slackCk, slackCs := slackApp.GetCreds()
	if twitterCk != "" && twitterCs != "" && slackCk != "" && slackCs != "" {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	w, r := c.Writer, c.Request
	switch c.Request.Method {
	case http.MethodPost:
		postSetup(w, r)
	case http.MethodGet:
		getSetup(c)
	default:
		http.NotFound(w, r)
	}
}

func postSetup(w http.ResponseWriter, r *http.Request) {
	msg := ""
	defer func() {
		if len(msg) != 0 {
			msgCookie := &http.Cookie{
				Name:  "mybot.setup.message",
				Value: msg,
				Path:  "/setup/",
			}
			http.SetCookie(w, msgCookie)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}()

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		msg = err.Error()
		return
	}
	val := r.MultipartForm.Value

	twitterCk := val["setup.twitter_consumer_key"][0]
	twitterCs := val["setup.twitter_consumer_secret"][0]
	if twitterCk != "" && twitterCs != "" {
		twitterApp.SetCreds(twitterCk, twitterCs)
		err = twitterApp.Save()
		if err != nil {
			panic(err)
		}
	} else {
		msg = "Both of Twitter Consumer Key and Consumer Secret can't be empty"
	}

	slackCk := val["setup.slack_client_id"][0]
	slackCs := val["setup.slack_client_secret"][0]
	if slackCk != "" && slackCs != "" {
		slackApp.SetCreds(slackCk, slackCs)
		err = slackApp.Save()
		if err != nil {
			panic(err)
		}
	} else {
		msg = "Both of Slack Client ID and Client Secret can't be empty"
	}
}

func getSetup(c *gin.Context) {
	msg := ""
	msgCookie, err := c.Request.Cookie("mybot.setup.message")
	if err == nil && msgCookie != nil {
		msg = msgCookie.Value
	}

	twitterCk, twitterCs := twitterApp.GetCreds()
	slackCk, slackCs := slackApp.GetCreds()
	data := gin.H{
		"NavbarName":            "",
		"TwitterName":           "",
		"SlackTeam":             "",
		"SlackURL":              "",
		"GoogleEnabled":         googleEnabled(),
		"TwitterConsumerKey":    twitterCk,
		"TwitterConsumerSecret": twitterCs,
		"SlackClientID":         slackCk,
		"SlackClientSecret":     slackCs,
		"Message":               msg,
	}

	if msgCookie != nil {
		msgCookie.Value = ""
		msgCookie.Path = "/setup/"
		http.SetCookie(c.Writer, msgCookie)
	}

	c.HTML(http.StatusOK, "setup", data)
}

func loginHandler(c *gin.Context) {
	data := gin.H{
		"NavbarName":    "",
		"TwitterName":   "",
		"SlackTeam":     "",
		"SlackURL":      "",
		"GoogleEnabled": googleEnabled(),
	}
	c.HTML(http.StatusOK, "login", data)
}

func twitterUserSearchHandler(c *gin.Context) {
	r := c.Request
	twitterUser, err := authenticator.GetLoginUser(r)
	if err != nil {
		panic(err)
	}
	data := userSpecificDataMap[fmt.Sprintf(appUserIDFormat, twitterUser.Provider, twitterUser.UserID)]
	getTwitterUserSearch(c, data.twitterAPI)
}

func getTwitterUserSearch(c *gin.Context, twitterAPI *core.TwitterAPI) {
	r := c.Request
	vals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	searchTerm := vals.Get("q")
	vals.Del("q")
	res, err := twitterAPI.GetUserSearch(searchTerm, vals)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, res)
}

func authCallbackHandler(c *gin.Context) {
	w, r := c.Writer, c.Request
	provider := c.Param("provider")
	switch provider {
	case "twitter":
		getAuthTwitterCallback(w, r)
	case "slack":
		getAuthSlackCallback(w, r)
	default:
		http.NotFound(w, r)
	}
}

func getAuthCallback(w http.ResponseWriter, r *http.Request) (goth.User, *userSpecificData, error) {
	user, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		return goth.User{}, nil, err
	}

	err = authenticator.Login(user, w, r)
	if err != nil {
		return goth.User{}, nil, err
	}

	id := fmt.Sprintf(appUserIDFormat, user.Provider, user.UserID)
	data, exists := userSpecificDataMap[id]
	if exists {
		return user, data, nil
	}

	data, err = newUserSpecificData(cliContext, database, id)
	if err != nil {
		return goth.User{}, nil, err
	}
	err = startUserSpecificData(cliContext, data, id)
	if err != nil {
		return goth.User{}, nil, err
	}
	userSpecificDataMap[id] = data
	return user, data, nil
}

func getAuthTwitterCallback(w http.ResponseWriter, r *http.Request) {
	user, data, err := getAuthCallback(w, r)
	if err != nil {
		panic(err)
	}
	data.twitterAuth.SetCreds(user.AccessToken, user.AccessTokenSecret)
	err = data.twitterAuth.Save()
	if err != nil {
		panic(err)
	}
	*data.twitterAPI = *core.NewTwitterAPIWithAuth(data.twitterAuth, data.config, data.cache)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getAuthSlackCallback(w http.ResponseWriter, r *http.Request) {
	user, data, err := getAuthCallback(w, r)
	if err != nil {
		panic(err)
	}
	data.slackAuth.SetCreds(user.AccessToken, user.AccessTokenSecret)
	err = data.slackAuth.Save()
	if err != nil {
		panic(err)
	}
	*data.slackAPI = *core.NewSlackAPIWithAuth(user.AccessToken, data.config, data.cache)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func authHandler(c *gin.Context) {
	w, r := c.Writer, c.Request
	provider := c.Param("provider")
	getAuth(w, r, provider)
}

func getAuth(w http.ResponseWriter, r *http.Request, provider string) {
	callback := r.URL.Query().Get("callback")
	authenticator.SetProvider(provider, r)
	if len(callback) == 0 {
		callback = fmt.Sprintf("http://%s/auth/callback/%s", r.Host, provider)
	}
	ck, cs := "", ""
	switch provider {
	case "twitter":
		ck, cs = twitterApp.GetCreds()
	case "slack":
		ck, cs = slackApp.GetCreds()
	}
	authenticator.InitProvider(provider, callback, ck, cs)
	gothic.BeginAuthHandler(w, r)
}

func logoutHandler(c *gin.Context) {
	w, r := c.Writer, c.Request
	err := authenticator.Logout(w, r)
	if err != nil {
		panic(err)
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func getTwitterInfo(twitterAPI *core.TwitterAPI) (id, screenName string) {
	if twitterAPI != nil {
		user, err := twitterAPI.GetSelf()
		if err == nil {
			return user.IdStr, user.ScreenName
		}
	}
	return "", ""
}

func getSlackInfo(slackAPI *core.SlackAPI) (string, string) {
	if slackAPI != nil {
		user, err := slackAPI.AuthTest()
		if err == nil {
			return user.Team, user.URL
		}
	}
	return "", ""
}

func googleEnabled() bool {
	if visionAPI == nil {
		return false
	}

	return visionAPI.Enabled()
}
