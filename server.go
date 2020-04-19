package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/iwataka/mybot/data"
	mybot "github.com/iwataka/mybot/lib"
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
	twitterUserIDPrefix = "twitter-"
	sessNameForProvider = "mybot-%s-session"
	trueValue           = "true"
)

var (
	htmlTemplateDir = filepath.Join(assetsDir, "tmpl")
	htmlTemplate    *template.Template
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
func (a *Authenticator) InitProvider(host, name, callback string) {
	if callback == "" {
		callback = fmt.Sprintf("http://%s/auth/%s/callback", host, name)
	}
	var p goth.Provider
	switch name {
	case "twitter":
		ck, cs := twitterApp.GetCreds()
		p = twitter.New(
			ck,
			cs,
			callback,
		)
	case "slack":
		ck, cs := slackApp.GetCreds()
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
func (a *Authenticator) CompleteUserAuth(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error) {
	sessKey := "mybot-user"
	sess, err := serverSession.Get(r, fmt.Sprintf(sessNameForProvider, provider))
	if err != nil {
		return goth.User{}, utils.WithStack(err)
	}
	val, exists := sess.Values[sessKey]
	if exists {
		if user, ok := val.(goth.User); ok {
			return user, nil
		}
	}

	a.SetProvider(provider, r)
	q := r.URL.Query()
	q.Add("state", "state")
	r.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(w, r)
	if err == nil {
		user.RawData = nil // RawData cannot be converted into session data currently
		sess.Values[sessKey] = user
		err = sess.Save(r, w)
		if err != nil {
			return goth.User{}, utils.WithStack(err)
		}
	}
	return user, utils.WithStack(err)
}

// Logout executes logout operation of the current login-user.
func (a *Authenticator) Logout(provider string, w http.ResponseWriter, r *http.Request) error {
	sess, err := serverSession.Get(r, fmt.Sprintf(sessNameForProvider, provider))
	if err != nil {
		return utils.WithStack(err)
	}
	sess.Options.MaxAge = -1
	err = sess.Save(r, w)
	if err != nil {
		return utils.WithStack(err)
	}

	a.SetProvider(provider, r)
	return gothic.Logout(w, r)
}

func wrapHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ck, cs := twitterApp.GetCreds(); ck == "" || cs == "" {
			http.Redirect(w, r, "/setup/", http.StatusSeeOther)
			return
		}

		if ck, cs := slackApp.GetCreds(); ck == "" || cs == "" {
			http.Redirect(w, r, "/setup/", http.StatusSeeOther)
			return
		}

		if _, err := authenticator.CompleteUserAuth("twitter", w, r); err != nil {
			http.Redirect(w, r, "/auth/twitter/", http.StatusSeeOther)
			return
		}

		f(w, r)
	}
}

func startServer(host, port, cert, key string) error {
	var err error
	gothic.Store = serverSession

	// View endpoints
	http.HandleFunc("/", wrapHandler(indexHandler))
	http.HandleFunc("/twitter-collections/", wrapHandler(twitterColsHandler))
	http.HandleFunc("/config/", wrapHandler(configHandler))
	http.HandleFunc("/config/file/", wrapHandler(configFileHandler))
	http.HandleFunc("/config/timelines/add", wrapHandler(configTimelineAddHandler))
	http.HandleFunc("/config/favorites/add", wrapHandler(configFavoriteAddHandler))
	http.HandleFunc("/config/searches/add", wrapHandler(configSearchAddHandler))
	http.HandleFunc("/config/messages/add", wrapHandler(configMessageAddHandler))
	http.HandleFunc("/assets/css/", getAssetsCSS)
	http.HandleFunc("/assets/js/", getAssetsJS)
	http.HandleFunc("/auth/twitter/", getAuthTwitter)
	http.HandleFunc("/auth/slack", getAuthSlack)
	http.HandleFunc("/auth/twitter/callback", getAuthTwitterCallback)
	http.HandleFunc("/auth/slack/callback", getAuthSlackCallback)
	http.HandleFunc("/setup/", setupHandler)
	http.HandleFunc("/logout/twitter/", twitterLogoutHandler)
	// For Twitter user auto-completion usage
	http.HandleFunc("/twitter/users/search/", wrapHandler(twitterUserSearchHandler))

	addr := fmt.Sprintf("%s:%s", host, port)
	_, certErr := os.Stat(cert)
	_, keyErr := os.Stat(key)
	if certErr == nil && keyErr == nil {
		fmt.Printf("Listen on %s://%s\n", "https", addr)
		err = http.ListenAndServeTLS(addr, cert, key, nil)
	} else {
		fmt.Printf("Listen on %s://%s\n", "http", addr)
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		return utils.WithStack(err)
	}
	return nil
}

func generateHTMLTemplate() (*template.Template, error) {
	if htmlTemplate != nil {
		return htmlTemplate, nil
	}
	tmpl := template.New("mybot_template_root").Funcs(templateFuncMap).Delims("{{{", "}}}")
	var err error
	htmlTemplate, err = generateHTMLTemplateFromFiles(tmpl)
	return htmlTemplate, err
}

func generateHTMLTemplateFromFiles(tmpl *template.Template) (*template.Template, error) {
	tmplFiles, err := zglob.Glob(filepath.Join(htmlTemplateDir, "**", "*.tmpl"))
	if err != nil {
		return nil, utils.WithStack(err)
	}
	template, err := tmpl.ParseFiles(tmplFiles...)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	return template, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	switch r.URL.Path {
	case "/":
		getIndex(w, r, data.cache, data.twitterAPI, data.slackAPI, twitterUser, data.statuses)
	default:
		http.NotFound(w, r)
	}
}

func getIndex(
	w http.ResponseWriter,
	r *http.Request,
	cache data.Cache,
	twitterAPI *mybot.TwitterAPI,
	slackAPI *mybot.SlackAPI,
	twitterUser goth.User,
	statuses map[int]bool,
) {
	slackTeam, slackURL := getSlackInfo(slackAPI)
	imgURL, imgSrc, imgAnalysisResult, imgAnalysisDate := imageAnalysis(cache)

	data := &struct {
		NavbarName               string
		TwitterName              string
		SlackTeam                string
		SlackURL                 string
		GoogleEnabled            bool
		ImageURL                 string
		ImageSource              string
		ImageAnalysisResult      string
		ImageAnalysisDate        string
		TwitterListenDMStatus    bool
		TwitterListenUsersStatus bool
		TwitterPeriodicStatus    bool
		SlackListenerStatus      bool
	}{
		"",
		twitterUser.NickName,
		slackTeam,
		slackURL,
		googleEnabled(),
		imgURL,
		imgSrc,
		imgAnalysisResult,
		imgAnalysisDate,
		statuses[twitterDMRoutineKey],
		statuses[twitterUserRoutineKey],
		statuses[twitterPeriodicRoutineKey],
		statuses[slackRoutineKey],
	}

	buf := new(bytes.Buffer)
	tmpl, err := generateHTMLTemplate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(buf, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func twitterColsHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	switch r.Method {
	case http.MethodGet:
		getTwitterCols(w, r, data.slackAPI, data.twitterAPI, twitterUser)
	default:
		http.NotFound(w, r)
	}
}

func getTwitterCols(w http.ResponseWriter, r *http.Request, slackAPI *mybot.SlackAPI, twitterAPI *mybot.TwitterAPI, twitterUser goth.User) {
	colMap := make(map[string]string)
	activeCol := ""
	id, err := strconv.ParseInt(twitterUser.UserID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	slackTeam, slackURL := getSlackInfo(slackAPI)
	data := &struct {
		NavbarName       string
		TwitterName      string
		SlackTeam        string
		SlackURL         string
		GoogleEnabled    bool
		CollectionMap    map[string]string
		ActiveCollection string
	}{
		"TwitterCols",
		twitterUser.NickName,
		slackTeam,
		slackURL,
		googleEnabled(),
		colMap,
		activeCol,
	}

	buf := new(bytes.Buffer)
	tmpl, err := generateHTMLTemplate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(buf, "twitterCols", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func configHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	switch r.Method {
	case http.MethodPost:
		postConfig(w, r, data.config, twitterUser)
	case http.MethodGet:
		getConfig(w, r, data.config, data.slackAPI, twitterUser)
	default:
		http.NotFound(w, r)
	}
}

func postConfig(w http.ResponseWriter, r *http.Request, config mybot.Config, twitterUser goth.User) {
	var err error

	defer func() {
		if err == nil {
			err = config.Save()
			if err != nil {
				go reloadWorkers(twitterUserIDPrefix + twitterUser.UserID)
			}
		} else {
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

func setTimelinesToConfig(config mybot.Config, r *http.Request) error {
	val := r.MultipartForm.Value

	prefix := "twitter.timelines"
	deletedFlags := val[prefix+".deleted"]
	timelines := []mybot.TimelineConfig{}
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
		timeline := mybot.NewTimelineConfig()
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

func setFavoritesToConfig(config mybot.Config, r *http.Request) error {
	val := r.MultipartForm.Value

	prefix := "twitter.favorites"
	deletedFlags := val[prefix+".deleted"]
	favorites := []mybot.FavoriteConfig{}
	actions, err := postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return err
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == trueValue {
			continue
		}
		favorite := mybot.NewFavoriteConfig()
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

func setSearchesToConfig(config mybot.Config, r *http.Request) error {
	val := r.MultipartForm.Value

	prefix := "twitter.searches"
	deletedFlags := val[prefix+".deleted"]
	searches := []mybot.SearchConfig{}
	actions, err := postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return err
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == trueValue {
			continue
		}
		search := mybot.NewSearchConfig()
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

func setMessagesToConfig(config mybot.Config, r *http.Request) error {
	val := r.MultipartForm.Value

	prefix := "slack.messages"
	deletedFlags := val[prefix+".deleted"]
	msgs := []mybot.MessageConfig{}
	actions, err := postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return err
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == trueValue {
			continue
		}
		msg := mybot.NewMessageConfig()
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

func postConfigForFilter(r *http.Request, i int, prefix string) (mybot.Filter, error) {
	val := r.MultipartForm.Value

	prefix = prefix + ".filter."
	filter := mybot.NewFilter()
	filter.Patterns = tmpl.GetListTextboxValue(val, i, prefix+"patterns")
	filter.URLPatterns = tmpl.GetListTextboxValue(val, i, prefix+"url_patterns")
	filter.HasMedia = tmpl.GetBoolSelectboxValue(val, i, prefix+"has_media")
	fThreshold, err := tmpl.GetIntPtr(val, i, prefix+"favorite_threshold")
	if err != nil {
		return mybot.NewFilter(), utils.WithStack(err)
	}
	filter.FavoriteThreshold = fThreshold
	rThreshold, err := tmpl.GetIntPtr(val, i, prefix+"retweeted_threshold")
	if err != nil {
		return mybot.NewFilter(), utils.WithStack(err)
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
		return mybot.NewFilter(), utils.WithStack(err)
	}
	filter.Language.MinSentiment = minSentiment
	maxSentiment, err := tmpl.GetFloat64Ptr(val, i, prefix+"language.max_sentiment")
	if err != nil {
		return mybot.NewFilter(), utils.WithStack(err)
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

func getConfig(w http.ResponseWriter, r *http.Request, config mybot.Config, slackAPI *mybot.SlackAPI, twitterUser goth.User) {
	msg := ""
	msgCookie, err := r.Cookie("mybot.config.message")
	if err == nil {
		msg = msgCookie.Value
	}

	if msgCookie != nil {
		msgCookie.Value = ""
		msgCookie.Path = "/config/"
		http.SetCookie(w, msgCookie)
	}

	slackTeam, slackURL := getSlackInfo(slackAPI)
	bs, err := configPage(twitterUser.NickName, slackTeam, slackURL, msg, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func configPage(twitterName, slackTeam, slackURL, msg string, config mybot.Config) ([]byte, error) {
	data := &struct {
		NavbarName    string
		TwitterName   string
		SlackTeam     string
		SlackURL      string
		GoogleEnabled bool
		Message       string
		Config        mybot.ConfigProperties
	}{
		"Config",
		twitterName,
		slackTeam,
		slackURL,
		googleEnabled(),
		msg,
		*config.GetProperties(),
	}

	buf := new(bytes.Buffer)
	tmpl, err := generateHTMLTemplate()
	if err != nil {
		return nil, utils.WithStack(err)
	}
	err = tmpl.ExecuteTemplate(buf, "config", data)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	return buf.Bytes(), nil
}

func configTimelineAddHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodPost {
		postConfigTimelineAdd(w, r, data.config)
	}
}

func postConfigTimelineAdd(w http.ResponseWriter, r *http.Request, config mybot.Config) {
	addTimelineConfig(config)
	http.Redirect(w, r, "/config/", http.StatusSeeOther)
}

func addTimelineConfig(config mybot.Config) {
	config.AddTwitterTimeline(mybot.NewTimelineConfig())
}

func configFavoriteAddHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodPost {
		postConfigFavoriteAdd(w, r, data.config)
	}
}

func postConfigFavoriteAdd(w http.ResponseWriter, r *http.Request, config mybot.Config) {
	addFavoriteConfig(config)
	http.Redirect(w, r, "/config/", http.StatusSeeOther)
}

func addFavoriteConfig(config mybot.Config) {
	config.AddTwitterFavorite(mybot.NewFavoriteConfig())
}

func configSearchAddHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodPost {
		postConfigSearchAdd(w, r, data.config)
	} else {
		http.NotFound(w, r)
	}
}

func postConfigSearchAdd(w http.ResponseWriter, r *http.Request, config mybot.Config) {
	addSearchConfig(config)
	http.Redirect(w, r, "/config/", http.StatusSeeOther)
}

func addSearchConfig(config mybot.Config) {
	config.AddTwitterSearch(mybot.NewSearchConfig())
}

func configMessageAddHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodPost {
		postConfigMessageAdd(w, r, data.config)
	}
}

func postConfigMessageAdd(w http.ResponseWriter, r *http.Request, c mybot.Config) {
	addMessageConfig(c)
	http.Redirect(w, r, "/config/", http.StatusSeeOther)
}

func addMessageConfig(config mybot.Config) {
	config.AddSlackMessage(mybot.NewMessageConfig())
}

func configFileHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	switch r.Method {
	case http.MethodPost:
		postConfigFile(w, r, data.config)
	case http.MethodGet:
		getConfigFile(w, r, data.config)
	default:
		http.NotFound(w, r)
	}
}

func postConfigFile(w http.ResponseWriter, r *http.Request, config mybot.Config) {
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
	err = config.Unmarshal(".json", bytes)
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

func getConfigFile(w http.ResponseWriter, r *http.Request, config mybot.Config) {
	ext := ".json"
	w.Header().Add("Content-Type", "application/force-download; charset=utf-8")
	w.Header().Add("Content-Disposition", `attachment; filename="config`+ext+`"`)
	bytes, err := config.Marshal(strings.Repeat(" ", 4), ext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	len, err := w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len), 16))
}

func getAssetsJS(w http.ResponseWriter, r *http.Request) {
	getAssets(w, r, "application/javascript")
}

func getAssetsCSS(w http.ResponseWriter, r *http.Request) {
	getAssets(w, r, "text/css")
}

func getAssets(w http.ResponseWriter, r *http.Request, contentType string) {
	path := r.URL.Path[len("/"):]
	data, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", contentType)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func setupHandler(w http.ResponseWriter, r *http.Request) {
	twitterCk, twitterCs := twitterApp.GetCreds()
	slackCk, slackCs := slackApp.GetCreds()
	if twitterCk != "" && twitterCs != "" && slackCk != "" && slackCs != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodPost:
		postSetup(w, r)
	case http.MethodGet:
		getSetup(w, r)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		msg = "Both of Slack Client ID and Client Secret can't be empty"
	}
}

func getSetup(w http.ResponseWriter, r *http.Request) {
	msg := ""
	msgCookie, err := r.Cookie("mybot.setup.message")
	if msgCookie != nil {
		msg = msgCookie.Value
	}

	twitterCk, twitterCs := twitterApp.GetCreds()
	slackCk, slackCs := slackApp.GetCreds()
	data := &struct {
		NavbarName            string
		TwitterName           string
		SlackTeam             string
		SlackURL              string
		GoogleEnabled         bool
		TwitterConsumerKey    string
		TwitterConsumerSecret string
		SlackClientID         string
		SlackClientSecret     string
		Message               string
	}{
		"",
		"",
		"",
		"",
		googleEnabled(),
		twitterCk,
		twitterCs,
		slackCk,
		slackCs,
		msg,
	}

	if msgCookie != nil {
		msgCookie.Value = ""
		msgCookie.Path = "/setup/"
		http.SetCookie(w, msgCookie)
	}

	buf := new(bytes.Buffer)
	tmpl, err := generateHTMLTemplate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(buf, "setup", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func twitterUserSearchHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	switch r.Method {
	case http.MethodGet:
		getTwitterUserSearch(w, r, data.twitterAPI)
	default:
		http.NotFound(w, r)
	}
}

func getTwitterUserSearch(w http.ResponseWriter, r *http.Request, twitterAPI *mybot.TwitterAPI) {
	vals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	searchTerm := vals.Get("q")
	vals.Del("q")
	res, err := twitterAPI.GetUserSearch(searchTerm, vals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bs, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAuthTwitterCallback(w http.ResponseWriter, r *http.Request) {
	user, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := twitterUserIDPrefix + user.UserID
	data, exists := userSpecificDataMap[id]
	if !exists {
		data, err = newUserSpecificData(cliContext, dbSession, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userSpecificDataMap[id] = data
	}

	data.twitterAuth.SetCreds(user.AccessToken, user.AccessTokenSecret)
	err = data.twitterAuth.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	*data.twitterAPI = *mybot.NewTwitterAPIWithAuth(data.twitterAuth, data.config, data.cache)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getAuthSlackCallback(w http.ResponseWriter, r *http.Request) {
	user, err := authenticator.CompleteUserAuth("slack", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	data.slackAuth.SetCreds(user.AccessToken, user.AccessTokenSecret)
	err = data.slackAuth.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	*data.slackAPI = *mybot.NewSlackAPIWithAuth(user.AccessToken, data.config, data.cache)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getAuthTwitter(w http.ResponseWriter, r *http.Request) {
	getAuth("twitter", w, r)
}

func getAuthSlack(w http.ResponseWriter, r *http.Request) {
	getAuth("slack", w, r)
}

func getAuth(provider string, w http.ResponseWriter, r *http.Request) {
	callback := r.URL.Query().Get("callback")
	authenticator.SetProvider(provider, r)
	authenticator.InitProvider(r.Host, provider, callback)
	gothic.BeginAuthHandler(w, r)
}

func twitterLogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTwitterLogout(w, r)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func getTwitterLogout(w http.ResponseWriter, r *http.Request) {
	err := authenticator.Logout("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getSlackInfo(slackAPI *mybot.SlackAPI) (string, string) {
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
