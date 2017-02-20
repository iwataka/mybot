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

	"github.com/gin-gonic/contrib/sessions"
	"github.com/iwataka/mybot/lib"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
)

//go:generate go-bindata assets/...

const (
	// go1.5 or lower doesn't support http.MethodPost and else.
	methodPost      = "POST"
	methodGet       = "GET"
	htmlTemplateDir = "assets/tmpl"
)

var (
	htmlTemplate *template.Template
)

func init() {
	gothic.Store = sessions.NewCookieStore([]byte("mybot_session_key"))

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
		panic(err)
	}
}

func wrapHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status.UpdateTwitterAuth(twitterAPI)

		if !status.PassTwitterApp {
			w.Header().Add("Location", "/setup/twitter/")
			w.WriteHeader(http.StatusSeeOther)
			return
		}

		if !status.PassTwitterAuth {
			w.Header().Add("Location", "/auth/twitter/")
			w.WriteHeader(http.StatusSeeOther)
			return
		}

		f(w, r)
	}
}

func startServer(host, port, cert, key string) error {
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
		wrapHandler(getLog),
	)
	http.HandleFunc(
		"/status/",
		wrapHandler(getStatus),
	)
	http.HandleFunc(
		"/auth/twitter/",
		getAuthTwitter,
	)
	http.HandleFunc(
		"/auth/twitter/callback",
		getAuthTwitterCallback,
	)
	http.HandleFunc(
		"/setup/twitter/",
		setupTwitterHandler,
	)
	http.HandleFunc(
		"/hooks/",
		hooksHandler,
	)

	if len(host) == 0 {
		host = "localhost"
	}

	if len(port) == 0 {
		port = "3256"
	}

	var err error
	addr := fmt.Sprintf("%s:%s", host, port)
	_, certErr := os.Stat(cert)
	_, keyErr := os.Stat(key)
	if certErr == nil && keyErr == nil {
		fmt.Printf("Open %s://%s for more details\n", "https", addr)
		err = http.ListenAndServeTLS(addr, cert, key, nil)
	} else {
		fmt.Printf("Open %s://%s for more details\n", "http", addr)
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		return err
	}
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		getIndex(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	var botName string
	self, err := twitterAPI.GetSelf()
	if err == nil {
		botName = self.ScreenName
	} else {
		botName = ""
	}

	imageSource := ""
	imageURL := ""
	imageAnalysisResult := ""
	imageAnalysisDate := ""
	images, err := cache.GetLatestImages(1)
	if err == nil && len(images) != 0 {
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
	colList, err := twitterAPI.GetCollectionListByUserId(self.Id, nil)
	if err == nil {
		for _, c := range colList.Objects.Timelines {
			name := strings.Replace(c.Name, " ", "-", -1)
			colMap[name] = c.CollectionUrl
		}
	}

	data := &struct {
		NavbarName          string
		Log                 string
		BotName             string
		ImageURL            string
		ImageSource         string
		ImageAnalysisResult string
		ImageAnalysisDate   string
		CollectionMap       map[string]string
	}{
		"",
		"Currently you cannot see the log here",
		botName,
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
	if r.Method == methodPost {
		postConfig(w, r)
	} else if r.Method == methodGet {
		getConfig(w, r)
	}
}

func postConfig(w http.ResponseWriter, r *http.Request) {
	var err error
	valid := false

	defer func() {
		if valid {
			err = config.Save()
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
		w.Header().Add("Location", "/config/")
		w.WriteHeader(http.StatusSeeOther)
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
	config.Twitter.Timelines = timelines

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
	config.Twitter.Favorites = favorites

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
	config.Twitter.Searches = searches

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
	config.Slack.Messages = msgs

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
	config.IncomingWebhooks = incomings

	prefix = "twitter.notification"
	config.Twitter.Notification.Place.AllowSelf = len(val[prefix+".place.allow_self"]) > 1
	config.Twitter.Notification.Place.Users = mybot.GetListTextboxValue(val, 0, prefix+".place.users")

	prefix = "twitter.interaction"
	config.Twitter.Interaction.AllowSelf = len(val[prefix+".allow_self"]) > 1
	config.Twitter.Interaction.Users = mybot.GetListTextboxValue(val, 0, prefix+".users")

	config.Twitter.Duration = val["twitter.duration"][0]

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
		if deletedFlags[i] == "true" {
			tweetCounter.returnValue(i, val, false)
			retweetCounter.returnValue(i, val, false)
			favoriteCounter.returnValue(i, val, false)
			pinCounter.returnValue(i, val, false)
			starCounter.returnValue(i, val, false)
			continue
		}
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

func getConfig(w http.ResponseWriter, r *http.Request) {
	msg := ""
	msgCookie, err := r.Cookie("mybot.config.message")
	if err == nil {
		msg = msgCookie.Value
	}

	data := &struct {
		NavbarName string
		Message    string
		Config     mybot.FileConfig
	}{
		"Config",
		msg,
		*config,
	}

	if msgCookie != nil {
		msgCookie.Value = ""
		msgCookie.Path = "/config/"
		http.SetCookie(w, msgCookie)
	}

	buf := new(bytes.Buffer)
	err = htmlTemplate.ExecuteTemplate(buf, "config", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func configTimelineAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		postConfnigTimelineAdd(w, r)
	}
}

func postConfnigTimelineAdd(w http.ResponseWriter, r *http.Request) {
	timelines := config.Twitter.Timelines
	timelines = append(timelines, *mybot.NewTimelineConfig())
	config.Twitter.Timelines = timelines
	w.Header().Add("Location", "/config/")
	w.WriteHeader(http.StatusSeeOther)
}

func configFavoriteAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
	}
}

func postConfigFavoriteAdd(w http.ResponseWriter, r *http.Request) {
	favorites := config.Twitter.Favorites
	favorites = append(favorites, *mybot.NewFavoriteConfig())
	config.Twitter.Favorites = favorites
	w.Header().Add("Location", "/config/")
	w.WriteHeader(http.StatusSeeOther)
}

func configSearchAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		postConfigSearchAdd(w, r)
	}
}

func postConfigSearchAdd(w http.ResponseWriter, r *http.Request) {
	searches := config.Twitter.Searches
	searches = append(searches, *mybot.NewSearchConfig())
	config.Twitter.Searches = searches
	w.Header().Add("Location", "/config/")
	w.WriteHeader(http.StatusSeeOther)
}

func configMessageAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		postConfigMessageAdd(w, r)
	}
}

func postConfigMessageAdd(w http.ResponseWriter, r *http.Request) {
	msgs := config.Slack.Messages
	msgs = append(msgs, *mybot.NewMessageConfig())
	config.Slack.Messages = msgs
	w.Header().Add("Location", "/config/")
	w.WriteHeader(http.StatusSeeOther)
}

func configIncomingAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		postConfigIncomingAdd(w, r)
	}
}

func postConfigIncomingAdd(w http.ResponseWriter, r *http.Request) {
	hooks := config.IncomingWebhooks
	hooks = append(hooks, *mybot.NewIncomingWebhook())
	config.IncomingWebhooks = hooks
	w.Header().Add("Location", "/config/")
	w.WriteHeader(http.StatusSeeOther)
}

func configFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		postConfigFile(w, r)
	} else if r.Method == methodGet {
		getConfigFile(w, r)
	}
}

func postConfigFile(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Add("Location", "/config/")
		w.WriteHeader(http.StatusSeeOther)
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
	err = config.FromText(bytes)
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

func getConfigFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/force-download; charset=utf-8")
	w.Header().Add("Content-Disposition", `attachment; filename="config.toml"`)
	bytes, err := ioutil.ReadFile(config.File)
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

func getLog(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		NavbarName string
		Log        string
	}{
		"Log",
		"Currently you cannot see the log here",
	}

	buf := new(bytes.Buffer)
	if err := htmlTemplate.ExecuteTemplate(buf, "log", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		NavbarName               string
		Status                   mybot.Status
		TwitterListenDMStatus    bool
		TwitterListenUsersStatus bool
		SlackListenerStatus      bool
	}{
		"Status",
		*status,
		status.CheckTwitterListenDMStatus(),
		status.CheckTwitterListenUsersStatus(),
		status.CheckSlackListen(),
	}

	buf := new(bytes.Buffer)
	err := htmlTemplate.ExecuteTemplate(buf, "status", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func setupTwitterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		postSetupTwitter(w, r)
	} else if r.Method == methodGet {
		getSetupTwitter(w, r)
	}
}

func postSetupTwitter(w http.ResponseWriter, r *http.Request) {
	msg := ""
	defer func() {
		if len(msg) != 0 {
			msgCookie := &http.Cookie{
				Name:  "mybot.setup.twitter.message",
				Value: msg,
				Path:  "/setup/twitter/",
			}
			http.SetCookie(w, msgCookie)
		}
		w.Header().Add("Location", "/")
		w.WriteHeader(http.StatusSeeOther)
	}()

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		msg = err.Error()
		return
	}
	val := r.MultipartForm.Value

	ck := val["twitter_setup.consumer_key"][0]
	cs := val["twitter_setup.consumer_secret"][0]

	if ck != "" && cs != "" {
		twitterApp.ConsumerKey = ck
		twitterApp.ConsumerSecret = cs
		twitterApp.Encode()
	} else {
		msg = "Both of Consumer Key and Consumer Secret can't be empty"
	}
}

func getSetupTwitter(w http.ResponseWriter, r *http.Request) {
	msg := ""
	msgCookie, err := r.Cookie("mybot.setup.twitter.message")
	if msgCookie != nil {
		msg = msgCookie.Value
	}

	data := &struct {
		NavbarName     string
		Message        string
		ConsumerKey    string
		ConsumerSecret string
	}{
		"",
		msg,
		twitterApp.ConsumerKey,
		twitterApp.ConsumerSecret,
	}

	if msgCookie != nil {
		msgCookie.Value = ""
		msgCookie.Path = "/setup/twitter/"
		http.SetCookie(w, msgCookie)
	}

	buf := new(bytes.Buffer)
	err = htmlTemplate.ExecuteTemplate(buf, "twitter_setup", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func getAuthTwitterCallback(w http.ResponseWriter, r *http.Request) {
	setProvider(r, "twitter")

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	twitterAuth.AccessToken = user.AccessToken
	twitterAuth.AccessTokenSecret = user.AccessTokenSecret
	err = twitterAuth.Encode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	*twitterAPI = *mybot.NewTwitterAPI(twitterAuth, cache, config)

	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func getAuthTwitter(w http.ResponseWriter, r *http.Request) {
	setProvider(r, "twitter")
	initProvider(r.Host, "twitter")

	gothic.BeginAuthHandler(w, r)
}

func setProvider(req *http.Request, name string) {
	q := req.URL.Query()
	q.Add("provider", name)
	req.URL.RawQuery = q.Encode()
}

func initProvider(host, name string) {
	callback := fmt.Sprintf("http://%s/auth/%s/callback", host, name)
	var p goth.Provider
	switch name {
	case "twitter":
		p = twitter.New(
			twitterApp.ConsumerKey,
			twitterApp.ConsumerSecret,
			callback,
		)
	}
	if p != nil {
		goth.UseProviders(p)
	}
}

func hooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		postHooks(w, r)
	}
}

func postHooks(w http.ResponseWriter, r *http.Request) {
	cs := []mybot.IncomingWebhook{}
	for _, c := range config.IncomingWebhooks {
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
				if err := slackAPI.PostMesage(ch, msg, nil); err != nil {
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
