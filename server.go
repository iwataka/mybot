package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
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
		return goth.User{}, err
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
			return goth.User{}, err
		}
	}
	return user, err
}

func (a *Authenticator) Logout(provider string, w http.ResponseWriter, r *http.Request) error {
	sess, err := serverSession.Get(r, fmt.Sprintf("mybot-%s-session", provider))
	if err != nil {
		return err
	}
	sess.Options.MaxAge = -1
	err = sess.Save(r, w)
	if err != nil {
		return err
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
		return err
	}

	http.HandleFunc(
		"/",
		wrapHandler(indexHandler),
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
		"/config/incomings/add",
		wrapHandler(configIncomingAddHandler),
	)
	http.HandleFunc(
		"/assets/css/",
		getAssetsCSS,
	)
	http.HandleFunc(
		"/log/",
		wrapHandler(logHandler),
	)
	http.HandleFunc(
		"/status/",
		wrapHandler(statusHandler),
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
		"/hooks/",
		hooksHandler,
	)
	http.HandleFunc(
		"/logout/twitter/",
		twitterLogoutHandler,
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
		return err
	}
	return nil
}

func initServer() error {
	gothic.Store = serverSession

	tmplTexts := []string{}
	for _, name := range AssetNames() {
		tmplBytes := MustAsset(name)
		tmplTexts = append(tmplTexts, string(tmplBytes))
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
		return err
	}
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.URL.Path == "/" {
		getIndex(w, r, data.cache, data.twitterAPI, data.slackAPI, twitterUser)
	} else {
		http.NotFound(w, r)
	}
}

func getIndex(w http.ResponseWriter, r *http.Request, cache mybot.Cache, twitterAPI *mybot.TwitterAPI, slackAPI *mybot.SlackAPI, twitterUser goth.User) {
	imageSource := ""
	imageURL := ""
	imageAnalysisResult := ""
	imageAnalysisDate := ""
	images := cache.GetLatestImages(1)
	if len(images) != 0 {
		imgCache := images[0]
		imageSource = imgCache.Src
		imageURL = imgCache.URL
		if cache != nil {
			buf := new(bytes.Buffer)
			result := imgCache.AnalysisResult
			err := json.Indent(buf, []byte(result), "", "  ")
			if err != nil {
				imageAnalysisResult = "Error while formatting the result"
			} else {
				imageAnalysisResult = buf.String()
			}
		}
		imageAnalysisDate = imgCache.AnalysisDate
	}

	colMap := make(map[string]string)
	id, err := strconv.ParseInt(twitterUser.UserID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	colList, err := twitterAPI.GetCollectionListByUserId(id, nil)
	if err == nil {
		for _, c := range colList.Objects.Timelines {
			name := strings.Replace(c.Name, " ", "-", -1)
			colMap[name] = c.CollectionUrl
		}
	}

	slackTeam, slackURL := getSlackInfo(w, r, slackAPI)
	data := &struct {
		NavbarName          string
		Log                 string
		TwitterName         string
		SlackTeam           string
		SlackURL            string
		ImageURL            string
		ImageSource         string
		ImageAnalysisResult string
		ImageAnalysisDate   string
		CollectionMap       map[string]string
	}{
		"",
		"Currently you cannot see the log here",
		twitterUser.NickName,
		slackTeam,
		slackURL,
		imageURL,
		imageSource,
		imageAnalysisResult,
		imageAnalysisDate,
		colMap,
	}

	err = htmlTemplate.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
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
		timeline := *mybot.NewTimelineConfig()
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
		favorite := *mybot.NewFavoriteConfig()
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
		search := *mybot.NewSearchConfig()
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
		msg := *mybot.NewMessageConfig()
		msg.Channels = mybot.GetListTextboxValue(val, i, prefix+".channels")
		if msg.Filter, err = postConfigForFilter(val, i, prefix); err != nil {
			return
		}
		msg.Action = actions[i]
		msgs = append(msgs, msg)
	}
	config.SetSlackMessages(msgs)

	prefix = "incoming_webhooks"
	deletedFlags = val[prefix+".deleted"]
	incomings := []mybot.IncomingWebhook{}
	actions, err = postConfigForActions(val, prefix, deletedFlags)
	if err != nil {
		return
	}
	for i := 0; i < len(deletedFlags); i++ {
		if deletedFlags[i] == "true" {
			continue
		}
		in := *mybot.NewIncomingWebhook()
		in.Endpoint = mybot.GetString(val, prefix+".endpoint", i, "")
		in.Template = mybot.GetString(val, prefix+".template", i, "")
		in.Action = actions[i]
		incomings = append(incomings, in)
	}
	config.SetIncomingWebhooks(incomings)

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

func postConfigForFilter(val map[string][]string, i int, prefix string) (*mybot.Filter, error) {
	prefix = prefix + ".filter."
	filter := mybot.NewFilter()
	filter.Patterns = mybot.GetListTextboxValue(val, i, prefix+"patterns")
	filter.URLPatterns = mybot.GetListTextboxValue(val, i, prefix+"url_patterns")
	filter.HasMedia = mybot.GetBoolSelectboxValue(val, i, prefix+"has_media")
	filter.Retweeted = mybot.GetBoolSelectboxValue(val, i, prefix+"retweeted")
	fThreshold, err := mybot.GetIntPtr(val, i, prefix+"favorite_threshold")
	if err != nil {
		return nil, err
	}
	filter.FavoriteThreshold = fThreshold
	rThreshold, err := mybot.GetIntPtr(val, i, prefix+"retweeted_threshold")
	if err != nil {
		return nil, err
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
		return nil, err
	}
	filter.Language.MinSentiment = minSentiment
	maxSentiment, err := mybot.GetFloat64Ptr(val, i, prefix+"language.max_sentiment")
	if err != nil {
		return nil, err
	}
	filter.Language.MaxSentiment = maxSentiment
	return filter, nil
}

func postConfigForActions(
	val map[string][]string,
	prefix string,
	deletedFlags []string,
) ([]*mybot.Action, error) {
	prefix = prefix + ".action."
	tweetCounter := checkboxCounter{prefix + "twitter.tweet", 0}
	retweetCounter := checkboxCounter{prefix + "twitter.retweet", 0}
	favoriteCounter := checkboxCounter{prefix + "twitter.favorite", 0}
	pinCounter := checkboxCounter{prefix + "slack.pin", 0}
	starCounter := checkboxCounter{prefix + "slack.star", 0}
	results := []*mybot.Action{}
	for i := 0; i < len(deletedFlags); i++ {
		a, err := postConfigForAction(val, i, prefix)
		if err != nil {
			return nil, err
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

func postConfigForAction(val map[string][]string, i int, prefix string) (*mybot.Action, error) {
	action := mybot.NewAction()
	action.Twitter.Collections = mybot.GetListTextboxValue(val, i, prefix+"twitter.collections")
	action.Slack.Channels = mybot.GetListTextboxValue(val, i, prefix+"slack.channels")
	action.Slack.Reactions = mybot.GetListTextboxValue(val, i, prefix+"slack.reactions")
	action.OutgoingWebhook.Endpoint = mybot.GetString(val, prefix+"outgoing_webhook.endpoint", i, "")
	action.OutgoingWebhook.Method = mybot.GetString(val, prefix+"outgoing_webhook.method", i, "")
	action.OutgoingWebhook.Body = mybot.GetString(val, prefix+"outgoing_webhook.body", i, "")
	action.OutgoingWebhook.Template = mybot.GetString(val, prefix+"outgoing_webhook.template", i, "")
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

	slackTeam, slackURL := getSlackInfo(w, r, slackAPI)
	bs, err := configPage(twitterUser.NickName, slackTeam, slackURL, msg, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bs)
}

func configPage(twitterName, slackTeam, slackURL, msg string, config mybot.Config) ([]byte, error) {
	data := &struct {
		NavbarName  string
		TwitterName string
		SlackTeam   string
		SlackURL    string
		Message     string
		Config      mybot.ConfigProperties
	}{
		"Config",
		twitterName,
		slackTeam,
		slackURL,
		msg,
		*config.GetProperties(),
	}

	buf := new(bytes.Buffer)
	err := htmlTemplate.ExecuteTemplate(buf, "config", data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func configTimelineAddHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
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
	if r.Method == http.MethodPost {
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
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
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

func configIncomingAddHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodPost {
		postConfigIncomingAdd(w, r, data.config)
	}
}

func postConfigIncomingAdd(w http.ResponseWriter, r *http.Request, c mybot.Config) {
	addIncomingConfig(c)
	http.Redirect(w, r, "/config/", http.StatusSeeOther)
}

func addIncomingConfig(config mybot.Config) {
	hooks := config.GetIncomingWebhooks()
	hooks = append(hooks, *mybot.NewIncomingWebhook())
	config.SetIncomingWebhooks(hooks)
}

func configFileHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
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

func getAssetsCSS(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/"):]
	data, err := readFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodGet {
		getLog(w, r, data.slackAPI)
	} else {
		http.NotFound(w, r)
	}
}

func getLog(w http.ResponseWriter, r *http.Request, slackAPI *mybot.SlackAPI) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
		return
	}

	slackTeam, slackURL := getSlackInfo(w, r, slackAPI)
	data := &struct {
		NavbarName  string
		TwitterName string
		SlackTeam   string
		SlackURL    string
		Log         string
	}{
		"Log",
		twitterUser.NickName,
		slackTeam,
		slackURL,
		"Currently you cannot see the log here",
	}

	buf := new(bytes.Buffer)
	if err := htmlTemplate.ExecuteTemplate(buf, "log", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodGet {
		getStatus(w, r, data.slackAPI, data.statuses)
	} else {
		http.NotFound(w, r)
	}
}

func getStatus(w http.ResponseWriter, r *http.Request, slackAPI *mybot.SlackAPI, statuses map[int]*bool) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Redirect(w, r, "/setup/", http.StatusSeeOther)
		return
	}

	slackTeam, slackURL := getSlackInfo(w, r, slackAPI)
	data := &struct {
		NavbarName               string
		TwitterName              string
		SlackTeam                string
		SlackURL                 string
		TwitterListenDMStatus    bool
		TwitterListenUsersStatus bool
		TwitterPeriodicStatus    bool
		SlackListenerStatus      bool
	}{
		"Status",
		twitterUser.NickName,
		slackTeam,
		slackURL,
		*statuses[twitterDMRoutineKey],
		*statuses[twitterUserRoutineKey],
		*statuses[twitterPeriodicRoutineKey],
		*statuses[slackRoutineKey],
	}

	buf := new(bytes.Buffer)
	err = htmlTemplate.ExecuteTemplate(buf, "status", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func setupHandler(w http.ResponseWriter, r *http.Request) {
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

func hooksHandler(w http.ResponseWriter, r *http.Request) {
	twitterUser, err := authenticator.CompleteUserAuth("twitter", w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := userSpecificDataMap[twitterUserIDPrefix+twitterUser.UserID]

	if r.Method == http.MethodPost {
		postHooks(w, r, data.config, data.twitterAPI, data.slackAPI)
	}
}

func postHooks(w http.ResponseWriter, r *http.Request, config mybot.Config, twitterAPI *mybot.TwitterAPI, slackAPI *mybot.SlackAPI) {
	cs := []mybot.IncomingWebhook{}
	for _, c := range config.GetIncomingWebhooks() {
		if r.URL.Path == c.Endpoint {
			cs = append(cs, c)
		}
	}
	if len(cs) == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := new(interface{})
	if len(bs) != 0 {
		err = json.Unmarshal(bs, data)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, c := range cs {
		buf := new(bytes.Buffer)
		tmpl, err := template.New("template").Parse(c.Template)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = tmpl.Execute(buf, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		msg := buf.String()
		if slackAPI.Enabled() {
			for _, ch := range c.Action.Slack.Channels {
				if err := slackAPI.PostMessage(ch, msg, nil, true); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
		if c.Action.Twitter.Tweet {
			if _, err := twitterAPI.PostTweet(msg, nil); mybot.CheckTwitterError(err) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func readFile(path string) ([]byte, error) {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	data, err := Asset(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getSlackInfo(w http.ResponseWriter, r *http.Request, slackAPI *mybot.SlackAPI) (string, string) {
	if slackAPI != nil {
		user, err := slackAPI.AuthTest()
		if err == nil {
			return user.Team, user.URL
		}
	}
	return "", ""
}
