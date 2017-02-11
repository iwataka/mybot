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
		"withPrefix":          mybot.NewWithPrefix,
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
		"/config/apis/add",
		wrapHandler(configAPIAddHandler),
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
	log := logger.ReadString()
	lines := strings.Split(log, "\n")
	linenum := config.Log.Linenum
	head := len(lines) - linenum
	if head < 0 {
		head = 0
	}
	log = strings.Join(lines[head:len(lines)], "\n")
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
		log,
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

func (c *checkboxCounter) returnValue(index int, val map[string][]string) bool {
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

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		msg = err.Error()
		return
	}
	val := r.MultipartForm.Value

	deletedFlags := val["twitter.timelines.deleted"]
	if len(deletedFlags) != len(config.Twitter.Timelines) {
		http.Error(w, "Collapsed request", http.StatusInternalServerError)
		return
	}
	actionRetweetCounter := checkboxCounter{"twitter.timelines.action.twitter.retweet", 0}
	actionFavoriteCounter := checkboxCounter{"twitter.timelines.action.twitter.favorite", 0}
	actionFollowCounter := checkboxCounter{"twitter.timelines.action.twitter.follow", 0}
	length := len(val["twitter.timelines.count"])
	timelines := []mybot.TimelineConfig{}
	for i := 0; i < length; i++ {
		if deletedFlags[i] == "true" {
			actionRetweetCounter.returnValue(i, val)
			actionFavoriteCounter.returnValue(i, val)
			actionFollowCounter.returnValue(i, val)
			continue
		}
		timeline := *mybot.NewTimelineConfig()
		timeline.ScreenNames = mybot.GetListTextboxValue(val, i, "twitter.timelines.screen_names")
		timeline.ExcludeReplies = mybot.GetBoolSelectboxValue(val, i, "twitter.timelines.exclude_replies")
		timeline.IncludeRts = mybot.GetBoolSelectboxValue(val, i, "twitter.timelines.include_rts")
		count, err := mybot.GetIntPtr(val, i, "twitter.timelines.count")
		if err != nil {
			msg = err.Error()
			return
		}
		timeline.Count = count
		filter, err := postConfigForFilter(val, i, "twitter.timelines")
		if err != nil {
			msg = err.Error()
			return
		}
		timeline.Filter = filter
		action, err := postConfigForAction(val, i, "twitter.timelines")
		if err != nil {
			msg = err.Error()
			return
		}
		action.Twitter.Retweet = actionRetweetCounter.returnValue(i, val)
		action.Twitter.Favorite = actionFavoriteCounter.returnValue(i, val)
		action.Twitter.Follow = actionFollowCounter.returnValue(i, val)
		timeline.Action = action
		timelines = append(timelines, timeline)
	}
	config.Twitter.Timelines = timelines

	deletedFlags = val["twitter.favorites.deleted"]
	if len(deletedFlags) != len(config.Twitter.Favorites) {
		http.Error(w, "Collapsed request", http.StatusInternalServerError)
		return
	}
	actionRetweetCounter = checkboxCounter{"twitter.favorites.action.twitter.retweet", 0}
	actionFavoriteCounter = checkboxCounter{"twitter.favorites.action.twitter.favorite", 0}
	actionFollowCounter = checkboxCounter{"twitter.favorites.action.twitter.follow", 0}
	length = len(val["twitter.favorites.count"])
	favorites := []mybot.FavoriteConfig{}
	for i := 0; i < length; i++ {
		if deletedFlags[i] == "true" {
			actionRetweetCounter.returnValue(i, val)
			actionFavoriteCounter.returnValue(i, val)
			actionFollowCounter.returnValue(i, val)
			continue
		}
		favorite := *mybot.NewFavoriteConfig()
		favorite.ScreenNames = mybot.GetListTextboxValue(val, i, "twitter.favorites.screen_names")
		count, err := mybot.GetIntPtr(val, i, "twitter.favorites.count")
		if err != nil {
			msg = err.Error()
			return
		}
		favorite.Count = count
		filter, err := postConfigForFilter(val, i, "twitter.favorites")
		if err != nil {
			msg = err.Error()
			return
		}
		favorite.Filter = filter
		action, err := postConfigForAction(val, i, "twitter.favorites")
		if err != nil {
			msg = err.Error()
			return
		}
		action.Twitter.Retweet = actionRetweetCounter.returnValue(i, val)
		action.Twitter.Favorite = actionFavoriteCounter.returnValue(i, val)
		action.Twitter.Follow = actionFollowCounter.returnValue(i, val)
		favorite.Action = action
		favorites = append(favorites, favorite)
	}
	config.Twitter.Favorites = favorites

	deletedFlags = val["twitter.searches.deleted"]
	if len(deletedFlags) != len(config.Twitter.Searches) {
		http.Error(w, "Collapsed request", http.StatusInternalServerError)
		return
	}
	actionRetweetCounter = checkboxCounter{"twitter.searches.action.twitter.retweet", 0}
	actionFavoriteCounter = checkboxCounter{"twitter.searches.action.twitter.favorite", 0}
	actionFollowCounter = checkboxCounter{"twitter.searches.action.twitter.follow", 0}
	length = len(val["twitter.searches.count"])
	searches := []mybot.SearchConfig{}
	for i := 0; i < length; i++ {
		if deletedFlags[i] == "true" {
			actionRetweetCounter.returnValue(i, val)
			actionFavoriteCounter.returnValue(i, val)
			actionFollowCounter.returnValue(i, val)
			continue
		}
		search := *mybot.NewSearchConfig()
		search.Queries = mybot.GetListTextboxValue(val, i, "twitter.searches.queries")
		search.ResultType = val["twitter.searches.result_type"][i]
		count, err := mybot.GetIntPtr(val, i, "twitter.searches.count")
		if err != nil {
			msg = err.Error()
			return
		}
		search.Count = count
		filter, err := postConfigForFilter(val, i, "twitter.searches")
		if err != nil {
			msg = err.Error()
			return
		}
		search.Filter = filter
		action, err := postConfigForAction(val, i, "twitter.searches")
		if err != nil {
			msg = err.Error()
			return
		}
		action.Twitter.Retweet = actionRetweetCounter.returnValue(i, val)
		action.Twitter.Favorite = actionFavoriteCounter.returnValue(i, val)
		action.Twitter.Follow = actionFollowCounter.returnValue(i, val)
		search.Action = action
		searches = append(searches, search)
	}
	config.Twitter.Searches = searches

	deletedFlags = val["twitter.apis.deleted"]
	if len(deletedFlags) != len(config.Twitter.APIs) {
		http.Error(w, "Collapsed request", http.StatusInternalServerError)
		return
	}
	length = len(val["twitter.apis.source_url"])
	apis := []mybot.APIConfig{}
	for i := 0; i < length; i++ {
		if deletedFlags[i] == "true" {
			continue
		}
		api := *mybot.NewAPIConfig()
		api.SourceURL = val["twitter.apis.source_url"][i]
		api.MessageTemplate = val["twitter.apis.message_template"][i]
		apis = append(apis, api)
	}
	config.Twitter.APIs = apis

	config.Twitter.Notification.Place.AllowSelf = len(val["twitter.notification.place.allow_self"]) > 1
	config.Twitter.Notification.Place.Users = mybot.GetListTextboxValue(val, 0, "twitter.notification.place.users")

	config.Interaction.AllowSelf = len(val["interaction.allow_self"]) > 1
	config.Interaction.Users = mybot.GetListTextboxValue(val, 0, "interaction.users")

	config.Log.AllowSelf = len(val["log.allow_self"]) > 1
	config.Log.Users = mybot.GetListTextboxValue(val, 0, "log.users")
	linenum, err := strconv.Atoi(val["log.linenum"][0])
	if err != nil {
		msg = err.Error()
		return
	}
	config.Log.Linenum = linenum

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

func postConfigForFilter(val map[string][]string, i int, prefix string) (*mybot.TweetFilter, error) {
	filter := mybot.NewTweetFilter()
	filter.Patterns = mybot.GetListTextboxValue(val, i, prefix+".filter.patterns")
	filter.URLPatterns = mybot.GetListTextboxValue(val, i, prefix+".filter.url_patterns")
	filter.HasMedia = mybot.GetBoolSelectboxValue(val, i, prefix+".filter.has_media")
	filter.HasURL = mybot.GetBoolSelectboxValue(val, i, prefix+".filter.has_url")
	filter.Retweeted = mybot.GetBoolSelectboxValue(val, i, prefix+".filter.retweeted")
	fThreshold, err := mybot.GetIntPtr(val, i, prefix+".filter.favorite_threshold")
	if err != nil {
		return nil, err
	}
	filter.FavoriteThreshold = fThreshold
	rThreshold, err := mybot.GetIntPtr(val, i, prefix+".filter.retweeted_threshold")
	if err != nil {
		return nil, err
	}
	filter.RetweetedThreshold = rThreshold
	filter.Lang = val[prefix+".filter.lang"][i]
	filter.Vision.Label = mybot.GetListTextboxValue(val, i, prefix+".filter.vision.label")
	filter.Vision.Face.AngerLikelihood = val[prefix+".filter.vision.face.anger_likelihood"][i]
	filter.Vision.Face.BlurredLikelihood = val[prefix+".filter.vision.face.blurred_likelihood"][i]
	filter.Vision.Face.HeadwearLikelihood = val[prefix+".filter.vision.face.headwear_likelihood"][i]
	filter.Vision.Face.JoyLikelihood = val[prefix+".filter.vision.face.joy_likelihood"][i]
	filter.Vision.Text = mybot.GetListTextboxValue(val, i, prefix+".filter.vision.text")
	filter.Vision.Landmark = mybot.GetListTextboxValue(val, i, prefix+".filter.vision.landmark")
	filter.Vision.Logo = mybot.GetListTextboxValue(val, i, prefix+".filter.vision.logo")
	minSentiment, err := mybot.GetFloat64Ptr(val, i, prefix+".filter.language.min_sentiment")
	if err != nil {
		return nil, err
	}
	filter.Language.MinSentiment = minSentiment
	maxSentiment, err := mybot.GetFloat64Ptr(val, i, prefix+".filter.language.max_sentiment")
	if err != nil {
		return nil, err
	}
	filter.Language.MaxSentiment = maxSentiment
	return filter, nil
}

func postConfigForAction(val map[string][]string, i int, prefix string) (*mybot.TweetAction, error) {
	action := mybot.NewTweetAction()
	action.Twitter.Collections = mybot.GetListTextboxValue(val, i, prefix+".action.twitter.collections")
	action.Slack.Channels = mybot.GetListTextboxValue(val, i, prefix+".action.slack.channels")
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
		Config     mybot.Config
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

	err = htmlTemplate.ExecuteTemplate(w, "config", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func configAPIAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == methodPost {
		postConfigAPIAdd(w, r)
	}
}

func postConfigAPIAdd(w http.ResponseWriter, r *http.Request) {
	apis := config.Twitter.APIs
	apis = append(apis, *mybot.NewAPIConfig())
	config.Twitter.APIs = apis
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
		logger.ReadString(),
	}
	err := htmlTemplate.ExecuteTemplate(w, "log", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		NavbarName               string
		Status                   mybot.Status
		TwitterListenDMStatus    bool
		TwitterListenUsersStatus bool
	}{
		"Status",
		*status,
		status.CheckTwitterListenDMStatus(),
		status.CheckTwitterListenUsersStatus(),
	}
	err := htmlTemplate.ExecuteTemplate(w, "status", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		c := make(chan bool)
		defer close(c)
		status.AddMonitorChan(ctxt.String("twitter-app"), c)
		twitterApp.Encode()
		<-c
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

	err = htmlTemplate.ExecuteTemplate(w, "twitter_setup", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	c := make(chan bool)
	defer close(c)
	status.AddMonitorChan(ctxt.String("twitter"), c)
	err = twitterAuth.Encode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	<-c

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
