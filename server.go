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

	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/models"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/twitter"
)

const (
	htmlTemplateDir     = "assets/tmpl"
	twitterUserIDPrefix = "twitter-"
)

var (
	htmlTemplate  *template.Template
	authenticator models.Authenticator = &Authenticator{}
)

type Authenticator struct{}

func (a *Authenticator) SetProvider(req *http.Request, name string) {
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return name, nil
	}
}

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

func (a *Authenticator) CompleteUserAuth(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error) {
	sess, err := serverSession.Get(r, fmt.Sprintf("mybot-%s-session", provider))
	if err != nil {
		return goth.User{}, mybot.WithStack(err)
	}
	val, exists := sess.Values["mybot-user"]
	if exists {
		if user, ok := val.(goth.User); ok {
			return user, nil
		}
	}

	a.SetProvider(r, provider)
	q := r.URL.Query()
	q.Add("state", "state")
	r.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(w, r)
	if err == nil {
		user.RawData = nil // RawData cannot be converted into session data cerrently
		sess.Values["mybot-user"] = user
		err := sess.Save(r, w)
		if err != nil {
			return goth.User{}, mybot.WithStack(err)
		}
	}
	return user, mybot.WithStack(err)
}

func (a *Authenticator) Logout(provider string, w http.ResponseWriter, r *http.Request) error {
	sess, err := serverSession.Get(r, fmt.Sprintf("mybot-%s-session", provider))
	if err != nil {
		return mybot.WithStack(err)
	}
	sess.Options.MaxAge = -1
	err = sess.Save(r, w)
	if err != nil {
		return mybot.WithStack(err)
	}

	a.SetProvider(r, provider)
	return gothic.Logout(w, r)
}

func wrapHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ck, cs := twitterApp.GetCreds(); ck == "" || cs == "" {
			http.Redirect(w, r, "/setup/", http.StatusSeeOther)
		} else if ck, cs := slackApp.GetCreds(); ck == "" || cs == "" {
			http.Redirect(w, r, "/setup/", http.StatusSeeOther)
		} else if _, err := authenticator.CompleteUserAuth("twitter", w, r); err != nil {
			http.Redirect(w, r, "/auth/twitter/", http.StatusSeeOther)
		} else {
			f(w, r)
		}
	}
}

func startServer(host, port, cert, key string) error {
	err := initServer()
	if err != nil {
		return mybot.WithStack(err)
	}

	// View endpoints
	http.HandleFunc(
		"/",
		wrapHandler(indexHandler),
	)
	http.HandleFunc(
		"/twitter-collections/",
		wrapHandler(twitterColsHandler),
	)
	http.HandleFunc(
		"/config/",
		wrapHandler(configHandler),
	)
	http.HandleFunc(
		"/config/file/",
		wrapHandler(configFileHandler),
	)
	http.HandleFunc(
		"/config/timelines/add",
		wrapHandler(configTimelineAddHandler),
	)
	http.HandleFunc(
		"/config/favorites/add",
		wrapHandler(configFavoriteAddHandler),
	)
	http.HandleFunc(
		"/config/searches/add",
		wrapHandler(configSearchAddHandler),
	)
	http.HandleFunc(
		"/config/messages/add",
		wrapHandler(configMessageAddHandler),
	)
	http.HandleFunc(
		"/assets/css/",
		getAssetsCSS,
	)
	http.HandleFunc(
		"/assets/js/",
		getAssetsJS,
	)
	http.HandleFunc(
		"/auth/twitter/",
		getAuthTwitter,
	)
	http.HandleFunc(
		"/auth/slack",
		getAuthSlack,
	)
	http.HandleFunc(
		"/auth/twitter/callback",
		getAuthTwitterCallback,
	)
	http.HandleFunc(
		"/auth/slack/callback",
		getAuthSlackCallback,
	)
	http.HandleFunc(
		"/setup/",
		setupHandler,
	)
	http.HandleFunc(
		"/logout/twitter/",
		twitterLogoutHandler,
	)

	// API endpoints
	http.HandleFunc(
		"/config/json/",
		wrapHandler(configJsonHandler),
	)
	http.HandleFunc(
		"/setting/",
		wrapHandler(settingHandler),
	)
	http.HandleFunc(
		"/twitter/users/search/",
		twitterUserSearchHandler,
	)
	http.HandleFunc(
		"/twitter/collections/list/",
		twitterCollectionListByUserId,
	)

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
		return mybot.WithStack(err)
	}
	return nil
}

func initServer() error {
	gothic.Store = serverSession

	tmplTexts := []string{}
	for _, name := range AssetNames() {
		if filepath.Ext(name) == ".tmpl" {
			tmplBytes := MustAsset(name)
			tmplTexts = append(tmplTexts, string(tmplBytes))
		}
	}

	funcMap := template.FuncMap{
		"checkbox":            mybot.Checkbox,
		"boolSelectbox":       mybot.BoolSelectbox,
		"selectbox":           mybot.Selectbox,
		"listTextbox":         mybot.ListTextbox,
		"textboxOfFloat64Ptr": mybot.TextboxOfFloat64Ptr,
		"textboxOfIntPtr":     mybot.TextboxOfIntPtr,
		"newMap":              mybot.NewMap,
	}

	tmpl, err := template.
		New("mybot_template_root").
		Funcs(funcMap).
		Parse(strings.Join(tmplTexts, "\n"))
	htmlTemplate = tmpl

	if err != nil {
		return mybot.WithStack(err)
	}
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.URL.Path == "/" {
		getIndex(w, r, data.cache, data.twitterAPI, data.slackAPI, data.statuses)
	} else {
		http.NotFound(w, r)
	}
}

func getIndex(w http.ResponseWriter, r *http.Request, cache mybot.Cache, twitterAPI *mybot.TwitterAPI, slackAPI *mybot.SlackAPI, statuses map[int]*bool) {
	setting, err := generateSetting(twitterAPI, slackAPI, cache, statuses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
		setting.TwitterName,
		setting.SlackTeam,
		setting.SlackURL,
		setting.GoogleEnabled,
		setting.Image.URL,
		setting.Image.Src,
		setting.Image.AnalysisResult,
		setting.Image.AnalysisDate,
		setting.Status.TwitterDMListener,
		setting.Status.TwitterUserListener,
		setting.Status.TwitterPeriodicJob,
		setting.Status.SlackListener,
	}

	buf := new(bytes.Buffer)
	err = htmlTemplate.ExecuteTemplate(buf, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func twitterCollectionListByUserId(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodGet {
		getTwitterCollectionListByUserId(w, r, data.twitterAPI)
	} else {
		http.NotFound(w, r)
	}
}

func getTwitterCollectionListByUserId(w http.ResponseWriter, r *http.Request, twitterAPI *mybot.TwitterAPI) {
	user, err := twitterAPI.GetSelf()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := twitterAPI.GetCollectionListByUserId(user.Id, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bs, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bs)
}

func settingHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodGet {
		getSetting(w, r, data.twitterAPI, data.slackAPI, data.cache, data.statuses)
	} else {
		http.NotFound(w, r)
	}
}

func getSetting(w http.ResponseWriter, r *http.Request, twitterAPI *mybot.TwitterAPI, slackAPI *mybot.SlackAPI, cache mybot.Cache, statuses map[int]*bool) {
	setting, err := generateSetting(twitterAPI, slackAPI, cache, statuses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bs, err := json.Marshal(setting)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bs)
}

type SettingResponse struct {
	TwitterName   string         `json:"twitter_name" toml:"twitter_name" bson:"twitter_name"`
	SlackTeam     string         `json:"slack_team" toml:"slack_team" bson:"slack_team"`
	SlackURL      string         `json:"slack_url" toml:"slack_url" bson:"slack_url"`
	GoogleEnabled bool           `json:"google_enabled" toml:"google_enabled" bson:"google_enabled"`
	Status        StatusResponse `json:"status" toml:"status" bson:"status"`
	Image         ImageResponse  `json:"image" toml:"image" bson:"image"`
}

type StatusResponse struct {
	TwitterDMListener   bool `json:"twitter_dm_listener" toml:"twitter_dm_listener" bson:"twitter_dm_listener"`
	TwitterUserListener bool `json:"twitter_user_listener" toml:"twitter_user_listener" bson:"twitter_user_listener"`
	TwitterPeriodicJob  bool `json:"twitter_periodic_job" toml:"twitter_periodic_job" bson:"twitter_periodic_job"`
	SlackListener       bool `json:"slack_listener" toml:"slack_listener" bson:"slack_listener"`
}

type ImageResponse struct {
	URL            string `json:"url" toml:"url" bson:"url"`
	Src            string `json:"src" toml:"src" bson:"src"`
	AnalysisResult string `json:"analysis_result" toml:"analysis_result" bson:"analysis_result"`
	AnalysisDate   string `json:"analysis_date" toml:"analysis_date" bson:"analysis_date"`
}

func generateSetting(twitterAPI *mybot.TwitterAPI, slackAPI *mybot.SlackAPI, cache mybot.Cache, statuses map[int]*bool) (*SettingResponse, error) {
	twitterUser, err := twitterAPI.GetSelf()
	if err != nil {
		return nil, mybot.WithStack(err)
	}
	slackTeam, slackURL := getSlackInfo(slackAPI)

	status := StatusResponse{
		TwitterDMListener:   *statuses[twitterDMRoutineKey],
		TwitterUserListener: *statuses[twitterUserRoutineKey],
		TwitterPeriodicJob:  *statuses[twitterPeriodicRoutineKey],
		SlackListener:       *statuses[slackRoutineKey],
	}

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
	img := ImageResponse{
		URL:            imageURL,
		Src:            imageSource,
		AnalysisResult: imageAnalysisResult,
		AnalysisDate:   imageAnalysisDate,
	}

	return &SettingResponse{
		TwitterName:   twitterUser.ScreenName,
		SlackTeam:     slackTeam,
		SlackURL:      slackURL,
		GoogleEnabled: googleEnabled(),
		Status:        status,
		Image:         img,
	}, nil
}

func twitterColsHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodGet {
		getTwitterCols(w, r, data.slackAPI, data.twitterAPI, twitterUser)
	} else {
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
	colList, err := twitterAPI.GetCollectionListByUserId(id, nil)
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
	err = htmlTemplate.ExecuteTemplate(buf, "twitterCols", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

type checkboxCounter struct {
	name       string
	extraCount int
}

func (c *checkboxCounter) returnValue(index int, val map[string][]string, def bool) bool {
	vs := val[c.name]
	if len(vs) <= index {
		return def
	}
	if val[c.name][index+c.extraCount] == "true" {
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

	if r.Method == http.MethodPost {
		postConfig(w, r, data.config, twitterUser)
	} else if r.Method == http.MethodGet {
		getConfig(w, r, data.config, data.slackAPI, twitterUser)
	}
}

func postConfig(w http.ResponseWriter, r *http.Request, config mybot.Config, twitterUser goth.User) {
	var err error
	valid := false

	defer func() {
		if valid {
			err = config.Save()
			go reloadWorkers(twitterUserIDPrefix + twitterUser.UserID)
		} else {
			err = config.Load()
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
	val := r.MultipartForm.Value

	prefix := "twitter.timelines"
	deletedFlags := val[prefix+".deleted"]
	timelines := []mybot.TimelineConfig{}
	actions, err := postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == "true" {
			continue
		}
		timeline := mybot.NewTimelineConfig()
		timeline.ScreenNames = mybot.GetListTextboxValue(val, i, prefix+".screen_names")
		timeline.ExcludeReplies = mybot.GetBoolSelectboxValue(val, i, prefix+".exclude_replies")
		timeline.IncludeRts = mybot.GetBoolSelectboxValue(val, i, prefix+".include_rts")
		if timeline.Count, err = mybot.GetIntPtr(val, i, prefix+".count"); err != nil {
			return
		}
		if timeline.Filter, err = postConfigForFilter(val, i, prefix); err != nil {
			return
		}
		timeline.Action = actions[i]
		timelines = append(timelines, timeline)
	}
	config.SetTwitterTimelines(timelines)

	prefix = "twitter.favorites"
	deletedFlags = val[prefix+".deleted"]
	favorites := []mybot.FavoriteConfig{}
	actions, err = postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == "true" {
			continue
		}
		favorite := mybot.NewFavoriteConfig()
		favorite.ScreenNames = mybot.GetListTextboxValue(val, i, prefix+".screen_names")
		if favorite.Count, err = mybot.GetIntPtr(val, i, prefix+".count"); err != nil {
			return
		}
		if favorite.Filter, err = postConfigForFilter(val, i, prefix); err != nil {
			return
		}
		favorite.Action = actions[i]
		favorites = append(favorites, favorite)
	}
	config.SetTwitterFavorites(favorites)

	prefix = "twitter.searches"
	deletedFlags = val[prefix+".deleted"]
	searches := []mybot.SearchConfig{}
	actions, err = postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == "true" {
			continue
		}
		search := mybot.NewSearchConfig()
		search.Queries = mybot.GetListTextboxValue(val, i, prefix+".queries")
		search.ResultType = val[prefix+".result_type"][i]
		if search.Count, err = mybot.GetIntPtr(val, i, prefix+".count"); err != nil {
			return
		}
		if search.Filter, err = postConfigForFilter(val, i, prefix); err != nil {
			return
		}
		search.Action = actions[i]
		searches = append(searches, search)
	}
	config.SetTwitterSearches(searches)

	prefix = "slack.messages"
	deletedFlags = val[prefix+".deleted"]
	msgs := []mybot.MessageConfig{}
	actions, err = postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == "true" {
			continue
		}
		msg := mybot.NewMessageConfig()
		msg.Channels = mybot.GetListTextboxValue(val, i, prefix+".channels")
		if msg.Filter, err = postConfigForFilter(val, i, prefix); err != nil {
			return
		}
		msg.Action = actions[i]
		msgs = append(msgs, msg)
	}
	config.SetSlackMessages(msgs)

	prefix = "twitter.notification"
	notif := config.GetTwitterNotification()
	notif.Place.AllowSelf = len(val[prefix+".place.allow_self"]) > 1
	notif.Place.Users = mybot.GetListTextboxValue(val, 0, prefix+".place.users")
	config.SetTwitterNotification(notif)

	prefix = "twitter.interaction"
	intr := config.GetTwitterInteraction()
	intr.AllowSelf = len(val[prefix+".allow_self"]) > 1
	intr.Users = mybot.GetListTextboxValue(val, 0, prefix+".users")
	config.SetTwitterInteraction(intr)

	config.SetTwitterDuration(val["twitter.duration"][0])

	err = config.Validate()
	if err == nil {
		valid = true
	} else {
		return
	}
}

func postConfigForFilter(val map[string][]string, i int, prefix string) (mybot.Filter, error) {
	prefix = prefix + ".filter."
	filter := mybot.NewFilter()
	filter.Patterns = mybot.GetListTextboxValue(val, i, prefix+"patterns")
	filter.URLPatterns = mybot.GetListTextboxValue(val, i, prefix+"url_patterns")
	filter.HasMedia = mybot.GetBoolSelectboxValue(val, i, prefix+"has_media")
	filter.Retweeted = mybot.GetBoolSelectboxValue(val, i, prefix+"retweeted")
	fThreshold, err := mybot.GetIntPtr(val, i, prefix+"favorite_threshold")
	if err != nil {
		return mybot.NewFilter(), mybot.WithStack(err)
	}
	filter.FavoriteThreshold = fThreshold
	rThreshold, err := mybot.GetIntPtr(val, i, prefix+"retweeted_threshold")
	if err != nil {
		return mybot.NewFilter(), mybot.WithStack(err)
	}
	filter.RetweetedThreshold = rThreshold
	filter.Lang = mybot.GetString(val, prefix+"lang", i, "")
	filter.Vision.Label = mybot.GetListTextboxValue(val, i, prefix+"vision.label")
	filter.Vision.Face.AngerLikelihood = mybot.GetString(val, prefix+"vision.face.anger_likelihood", i, "")
	filter.Vision.Face.BlurredLikelihood = mybot.GetString(val, prefix+"vision.face.blurred_likelihood", i, "")
	filter.Vision.Face.HeadwearLikelihood = mybot.GetString(val, prefix+"vision.face.headwear_likelihood", i, "")
	filter.Vision.Face.JoyLikelihood = mybot.GetString(val, prefix+"vision.face.joy_likelihood", i, "")
	filter.Vision.Text = mybot.GetListTextboxValue(val, i, prefix+"vision.text")
	filter.Vision.Landmark = mybot.GetListTextboxValue(val, i, prefix+"vision.landmark")
	filter.Vision.Logo = mybot.GetListTextboxValue(val, i, prefix+"vision.logo")
	minSentiment, err := mybot.GetFloat64Ptr(val, i, prefix+"language.min_sentiment")
	if err != nil {
		return mybot.NewFilter(), mybot.WithStack(err)
	}
	filter.Language.MinSentiment = minSentiment
	maxSentiment, err := mybot.GetFloat64Ptr(val, i, prefix+"language.max_sentiment")
	if err != nil {
		return mybot.NewFilter(), mybot.WithStack(err)
	}
	filter.Language.MaxSentiment = maxSentiment
	return filter, nil
}

func postConfigForActions(
	val map[string][]string,
	prefix string,
	deletedFlags []string,
) ([]mybot.Action, error) {
	prefix = prefix + ".action."
	tweetCounter := checkboxCounter{prefix + "twitter.tweet", 0}
	retweetCounter := checkboxCounter{prefix + "twitter.retweet", 0}
	favoriteCounter := checkboxCounter{prefix + "twitter.favorite", 0}
	pinCounter := checkboxCounter{prefix + "slack.pin", 0}
	starCounter := checkboxCounter{prefix + "slack.star", 0}
	results := []mybot.Action{}
	for i := 0; i < len(deletedFlags); i++ {
		a, err := postConfigForAction(val, i, prefix)
		if err != nil {
			return nil, mybot.WithStack(err)
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

func postConfigForAction(val map[string][]string, i int, prefix string) (mybot.Action, error) {
	action := mybot.NewAction()
	action.Twitter.Collections = mybot.GetListTextboxValue(val, i, prefix+"twitter.collections")
	action.Slack.Channels = mybot.GetListTextboxValue(val, i, prefix+"slack.channels")
	action.Slack.Reactions = mybot.GetListTextboxValue(val, i, prefix+"slack.reactions")
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
	w.Write(bs)
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
	err := htmlTemplate.ExecuteTemplate(buf, "config", data)
	if err != nil {
		return nil, mybot.WithStack(err)
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

func configJsonHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodPost {
		postConfigJson(w, r, data.config)
	} else if r.Method == http.MethodGet {
		getConfigJson(w, r, data.config)
	}
}

func postConfigJson(w http.ResponseWriter, r *http.Request, config mybot.Config) {
	var err error
	defer func() {
		r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}()

	var bs []byte
	bs, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = config.Unmarshal(bs)
	if err != nil {
		return
	}
	err = config.Validate()
	if err != nil {
		config.Load()
		return
	}
	err = config.Save()
	if err != nil {
		return
	}
}

func getConfigJson(w http.ResponseWriter, r *http.Request, config mybot.Config) {
	var err error
	defer func() {
		r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}()

	var bs []byte
	bs, err = config.Marshal("  ", ".json")
	if err != nil {
		return
	}
	_, err = w.Write(bs)
	if err != nil {
		return
	}
}

func configFileHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodPost {
		postConfigFile(w, r, data.config)
	} else if r.Method == http.MethodGet {
		getConfigFile(w, r, data.config)
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
	err = config.Unmarshal(bytes)
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
	data, err := readFile(path)
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

	if r.Method == http.MethodPost {
		postSetup(w, r)
	} else if r.Method == http.MethodGet {
		getSetup(w, r)
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
	err = htmlTemplate.ExecuteTemplate(buf, "setup", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func twitterUserSearchHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodGet {
		getTwitterUserSearch(w, r, data.twitterAPI)
	} else {
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
	vals.Encode()
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
	w.Write(bs)
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
	*data.twitterAPI = *mybot.NewTwitterAPI(data.twitterAuth, data.cache, data.config)

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
	*data.slackAPI = *mybot.NewSlackAPI(user.AccessToken, data.config, data.cache)

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
	authenticator.SetProvider(r, provider)
	authenticator.InitProvider(r.Host, provider, callback)
	gothic.BeginAuthHandler(w, r)
}

func twitterLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getTwitterLogout(w, r)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func getTwitterLogout(w http.ResponseWriter, r *http.Request) {
	authenticator.Logout("twitter", w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func readFile(path string) ([]byte, error) {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, mybot.WithStack(err)
		}
		return data, nil
	}
	data, err := Asset(path)
	if err != nil {
		return nil, mybot.WithStack(err)
	}
	return data, nil
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
	} else {
		return visionAPI.Enabled()
	}
}
