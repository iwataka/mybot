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
	if len(cache.GetLatestImages(1)) != 0 {
		imgCache := cache.GetLatestImages(1)[0]
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
	actionRetweetCounter := checkboxCounter{"twitter.timelines.action.retweet", 0}
	actionFavoriteCounter := checkboxCounter{"twitter.timelines.action.favorite", 0}
	actionFollowCounter := checkboxCounter{"twitter.timelines.action.follow", 0}
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
		timeline.Filter.Patterns = mybot.GetListTextboxValue(val, i, "twitter.timelines.filter.patterns")
		timeline.Filter.URLPatterns = mybot.GetListTextboxValue(val, i, "twitter.timelines.filter.url_patterns")
		timeline.Filter.HasMedia = mybot.GetBoolSelectboxValue(val, i, "twitter.timelines.filter.has_media")
		timeline.Filter.HasURL = mybot.GetBoolSelectboxValue(val, i, "twitter.timelines.filter.has_url")
		timeline.Filter.Retweeted = mybot.GetBoolSelectboxValue(val, i, "twitter.timelines.filter.retweeted")
		fThreshold, err := mybot.GetIntPtr(val, i, "twitter.timelines.filter.favorite_threshold")
		if err != nil {
			msg = err.Error()
			return
		}
		timeline.Filter.FavoriteThreshold = fThreshold
		rThreshold, err := mybot.GetIntPtr(val, i, "twitter.timelines.filter.retweeted_threshold")
		if err != nil {
			msg = err.Error()
			return
		}
		timeline.Filter.RetweetedThreshold = rThreshold
		timeline.Filter.Lang = val["twitter.timelines.filter.lang"][i]
		timeline.Filter.Vision.Label = mybot.GetListTextboxValue(val, i, "twitter.timelines.filter.vision.label")
		timeline.Filter.Vision.Face.AngerLikelihood = val["twitter.timelines.filter.vision.face.anger_likelihood"][i]
		timeline.Filter.Vision.Face.BlurredLikelihood = val["twitter.timelines.filter.vision.face.blurred_likelihood"][i]
		timeline.Filter.Vision.Face.HeadwearLikelihood = val["twitter.timelines.filter.vision.face.headwear_likelihood"][i]
		timeline.Filter.Vision.Face.JoyLikelihood = val["twitter.timelines.filter.vision.face.joy_likelihood"][i]
		timeline.Filter.Vision.Text = mybot.GetListTextboxValue(val, i, "twitter.timelines.filter.vision.text")
		timeline.Filter.Vision.Landmark = mybot.GetListTextboxValue(val, i, "twitter.timelines.filter.vision.landmark")
		timeline.Filter.Vision.Logo = mybot.GetListTextboxValue(val, i, "twitter.timelines.filter.vision.logo")
		minSentiment, err := mybot.GetFloat64Ptr(val, i, "twitter.timelines.filter.language.min_sentiment")
		if err != nil {
			msg = err.Error()
			return
		}
		timeline.Filter.Language.MinSentiment = minSentiment
		maxSentiment, err := mybot.GetFloat64Ptr(val, i, "twitter.timelines.filter.language.max_sentiment")
		if err != nil {
			msg = err.Error()
			return
		}
		timeline.Filter.Language.MaxSentiment = maxSentiment
		timeline.Action.Retweet = actionRetweetCounter.returnValue(i, val)
		timeline.Action.Favorite = actionFavoriteCounter.returnValue(i, val)
		timeline.Action.Follow = actionFollowCounter.returnValue(i, val)
		timeline.Action.Collections = mybot.GetListTextboxValue(val, i, "twitter.timelines.action.collections")
		timelines = append(timelines, timeline)
	}
	config.Twitter.Timelines = timelines

	deletedFlags = val["twitter.favorites.deleted"]
	if len(deletedFlags) != len(config.Twitter.Favorites) {
		http.Error(w, "Collapsed request", http.StatusInternalServerError)
		return
	}
	actionRetweetCounter = checkboxCounter{"twitter.favorites.action.retweet", 0}
	actionFavoriteCounter = checkboxCounter{"twitter.favorites.action.favorite", 0}
	actionFollowCounter = checkboxCounter{"twitter.favorites.action.follow", 0}
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
		config.Twitter.Favorites[i] = favorite
		favorite.Filter.Patterns = mybot.GetListTextboxValue(val, i, "twitter.favorites.filter.patterns")
		favorite.Filter.URLPatterns = mybot.GetListTextboxValue(val, i, "twitter.favorites.filter.url_patterns")
		favorite.Filter.HasMedia = mybot.GetBoolSelectboxValue(val, i, "twitter.favorites.filter.has_media")
		favorite.Filter.HasURL = mybot.GetBoolSelectboxValue(val, i, "twitter.favorites.filter.has_url")
		favorite.Filter.Retweeted = mybot.GetBoolSelectboxValue(val, i, "twitter.favorites.filter.retweeted")
		fThreshold, err := mybot.GetIntPtr(val, i, "twitter.favorites.filter.favorite_threshold")
		if err != nil {
			msg = err.Error()
			return
		}
		favorite.Filter.FavoriteThreshold = fThreshold
		rThreshold, err := mybot.GetIntPtr(val, i, "twitter.favorites.filter.retweeted_threshold")
		if err != nil {
			msg = err.Error()
			return
		}
		favorite.Filter.RetweetedThreshold = rThreshold
		favorite.Filter.Lang = val["twitter.favorites.filter.lang"][i]
		favorite.Filter.Vision.Label = mybot.GetListTextboxValue(val, i, "twitter.favorites.filter.vision.label")
		favorite.Filter.Vision.Face.AngerLikelihood = val["twitter.favorites.filter.vision.face.anger_likelihood"][i]
		favorite.Filter.Vision.Face.BlurredLikelihood = val["twitter.favorites.filter.vision.face.blurred_likelihood"][i]
		favorite.Filter.Vision.Face.HeadwearLikelihood = val["twitter.favorites.filter.vision.face.headwear_likelihood"][i]
		favorite.Filter.Vision.Face.JoyLikelihood = val["twitter.favorites.filter.vision.face.joy_likelihood"][i]
		favorite.Filter.Vision.Text = mybot.GetListTextboxValue(val, i, "twitter.favorites.filter.vision.text")
		favorite.Filter.Vision.Landmark = mybot.GetListTextboxValue(val, i, "twitter.favorites.filter.vision.landmark")
		favorite.Filter.Vision.Logo = mybot.GetListTextboxValue(val, i, "twitter.favorites.filter.vision.logo")
		minSentiment, err := mybot.GetFloat64Ptr(val, i, "twitter.favorites.filter.language.min_sentiment")
		if err != nil {
			msg = err.Error()
			return
		}
		favorite.Filter.Language.MinSentiment = minSentiment
		maxSentiment, err := mybot.GetFloat64Ptr(val, i, "twitter.favorites.filter.language.max_sentiment")
		if err != nil {
			msg = err.Error()
			return
		}
		favorite.Filter.Language.MaxSentiment = maxSentiment
		favorite.Action.Retweet = actionRetweetCounter.returnValue(i, val)
		favorite.Action.Favorite = actionFavoriteCounter.returnValue(i, val)
		favorite.Action.Follow = actionFollowCounter.returnValue(i, val)
		favorite.Action.Collections = mybot.GetListTextboxValue(val, i, "twitter.favorites.action.collections")
		favorites = append(favorites, favorite)
	}
	config.Twitter.Favorites = favorites

	deletedFlags = val["twitter.searches.deleted"]
	if len(deletedFlags) != len(config.Twitter.Searches) {
		http.Error(w, "Collapsed request", http.StatusInternalServerError)
		return
	}
	actionRetweetCounter = checkboxCounter{"twitter.searches.action.retweet", 0}
	actionFavoriteCounter = checkboxCounter{"twitter.searches.action.favorite", 0}
	actionFollowCounter = checkboxCounter{"twitter.searches.action.follow", 0}
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
		config.Twitter.Searches[i] = search
		search.Filter.Patterns = mybot.GetListTextboxValue(val, i, "twitter.searches.filter.patterns")
		search.Filter.URLPatterns = mybot.GetListTextboxValue(val, i, "twitter.searches.filter.url_patterns")
		search.Filter.HasMedia = mybot.GetBoolSelectboxValue(val, i, "twitter.searches.filter.has_media")
		search.Filter.HasURL = mybot.GetBoolSelectboxValue(val, i, "twitter.searches.filter.has_url")
		search.Filter.Retweeted = mybot.GetBoolSelectboxValue(val, i, "twitter.searches.filter.retweeted")
		fThreshold, err := mybot.GetIntPtr(val, i, "twitter.searches.filter.favorite_threshold")
		if err != nil {
			msg = err.Error()
			return
		}
		search.Filter.FavoriteThreshold = fThreshold
		rThreshold, err := mybot.GetIntPtr(val, i, "twitter.searches.filter.retweeted_threshold")
		if err != nil {
			msg = err.Error()
			return
		}
		search.Filter.RetweetedThreshold = rThreshold
		search.Filter.Lang = val["twitter.searches.filter.lang"][i]
		search.Filter.Vision.Label = mybot.GetListTextboxValue(val, i, "twitter.searches.filter.vision.label")
		search.Filter.Vision.Face.AngerLikelihood = val["twitter.searches.filter.vision.face.anger_likelihood"][i]
		search.Filter.Vision.Face.BlurredLikelihood = val["twitter.searches.filter.vision.face.blurred_likelihood"][i]
		search.Filter.Vision.Face.HeadwearLikelihood = val["twitter.searches.filter.vision.face.headwear_likelihood"][i]
		search.Filter.Vision.Face.JoyLikelihood = val["twitter.searches.filter.vision.face.joy_likelihood"][i]
		search.Filter.Vision.Text = mybot.GetListTextboxValue(val, i, "twitter.searches.filter.vision.text")
		search.Filter.Vision.Landmark = mybot.GetListTextboxValue(val, i, "twitter.searches.filter.vision.landmark")
		search.Filter.Vision.Logo = mybot.GetListTextboxValue(val, i, "twitter.searches.filter.vision.logo")
		minSentiment, err := mybot.GetFloat64Ptr(val, i, "twitter.searches.filter.language.min_sentiment")
		if err != nil {
			msg = err.Error()
			return
		}
		search.Filter.Language.MinSentiment = minSentiment
		maxSentiment, err := mybot.GetFloat64Ptr(val, i, "twitter.searches.filter.language.max_sentiment")
		if err != nil {
			msg = err.Error()
			return
		}
		search.Filter.Language.MaxSentiment = maxSentiment
		search.Action.Retweet = actionRetweetCounter.returnValue(i, val)
		search.Action.Favorite = actionFavoriteCounter.returnValue(i, val)
		search.Action.Follow = actionFollowCounter.returnValue(i, val)
		search.Action.Collections = mybot.GetListTextboxValue(val, i, "twitter.searches.action.collections")
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
		c := make(chan bool, 2)
		defer close(c)
		status.AddMonitorTwitterCredChan(c)
		twitterAuth.Encode()
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
	status.AddMonitorTwitterCredChan(c)
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
