package mybot

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
)

//go:generate go-bindata -pkg mybot assets/... pages/...

// MybotServer shows various pieces of information to users, such as an error
// log, Google Vision API result and Twitter collections.
type MybotServer struct {
	// Logger is a logging utility instance of this application. This
	// returns a log file's content if users request.
	Logger *Logger
	// TwitterAPI is a client for Twitter API. This server requires some
	// pieces of information related to TWitter, so this is here.
	TwitterAPI *TwitterAPI
	// VisionAPI is a client for Google Vision API.
	//
	// TODO: This field may not be required (at this time only
	// VisionAPI.File is required).
	VisionAPI *VisionAPI
	// Cache is a cache of this application and contains some Vision API
	// analysis result. This server need to show them.
	//
	// TODO: In the future, this server will fetch Vision API results from
	// DB and thus this field will be removed.
	Cache *MybotCache
	// Config is a configuration of this application and this server use
	// this as the others do.
	Config *MybotConfig
	// Status is a status of all processes in this application. This
	// enables users monitor their status via browser.
	Status *MybotStatus
	// pass is a flag which represents whether Twitter API is authenticated
	// or not. When Twitter API is authenticated, then users can pass a
	// setup page and go to other pages, thus this is called 'pass'.
	pass bool
}

func (s *MybotServer) wrapHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.pass {
			f(w, r)
			return
		}

		ok, err := s.TwitterAPI.VerifyCredentials()
		if ok && err == nil {
			s.pass = true
			f(w, r)
			return
		}

		msg := ""
		if err != nil {
			msg = err.Error()
		} else {
			msg = "You should specify the below information"
		}
		msgCookie := &http.Cookie{
			Name:  "mybot.setup.message",
			Value: msg,
			Path:  "/setup/",
		}
		http.SetCookie(w, msgCookie)
		w.Header().Add("Location", "/setup/")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *MybotServer) Init(host, port, cert, key string) error {
	http.HandleFunc(
		"/",
		s.wrapHandler(s.handler),
	)
	http.HandleFunc(
		"/config/",
		s.wrapHandler(s.configHandler),
	)
	http.HandleFunc(
		"/config/timelines/add",
		s.wrapHandler(s.configTimelineAddHandler),
	)
	http.HandleFunc(
		"/config/favorites/add",
		s.wrapHandler(s.configFavoriteAddHandler),
	)
	http.HandleFunc(
		"/config/searches/add",
		s.wrapHandler(s.configSearchAddHandler),
	)
	http.HandleFunc(
		"/assets/",
		s.assetHandler,
	)
	http.HandleFunc(
		"/log/",
		s.wrapHandler(s.logHandler),
	)
	http.HandleFunc(
		"/status/",
		s.wrapHandler(s.statusHandler),
	)
	http.HandleFunc(
		"/setup/",
		s.setupHandler,
	)

	h := s.Config.Server.Host
	if len(host) != 0 {
		h = host
	}

	p := s.Config.Server.Port
	if len(port) != 0 {
		p = port
	}

	var err error
	addr := h + ":" + p
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

func (s *MybotServer) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		tmpl, err := generateTemplate("index", "src/pages/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log := s.Logger.ReadString()
		lines := strings.Split(log, "\n")
		lineNum := s.Config.Server.LogLines
		head := len(lines) - lineNum
		if head < 0 {
			head = 0
		}
		log = strings.Join(lines[head:len(lines)], "\n")
		var botName string
		self, err := s.TwitterAPI.GetSelf()
		if err == nil {
			botName = self.ScreenName
		} else {
			botName = ""
		}

		imageAnalysisResult := ""
		if s.Cache != nil {
			buf := new(bytes.Buffer)
			err := json.Indent(buf, []byte(s.Cache.ImageAnalysisResult), "", "  ")
			if err != nil {
				imageAnalysisResult = "Error while formatting the result"
			} else {
				imageAnalysisResult = buf.String()
			}
		}

		colMap := make(map[string]string)
		colList, err := s.TwitterAPI.api.GetCollectionListByUserId(self.Id, nil)
		if err == nil {
			for _, c := range colList.Objects.Timelines {
				name := strings.Replace(c.Name, " ", "-", -1)
				colMap[name] = c.CollectionUrl
			}
		}

		data := &struct {
			UserName            string
			Log                 string
			BotName             string
			ImageURL            string
			ImageSource         string
			ImageAnalysisResult string
			ImageAnalysisDate   string
			CollectionMap       map[string]string
		}{
			s.Config.Server.Name,
			log,
			botName,
			s.Cache.ImageURL,
			s.Cache.ImageSource,
			imageAnalysisResult,
			s.Cache.ImageAnalysisDate,
			colMap,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.NotFound(w, r)
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
	} else {
		return false
	}
}

func atoiOrDefault(str string, def int) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return def
	} else {
		return i
	}
}

func (s *MybotServer) configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		val := r.MultipartForm.Value

		deletedFlags := val["twitter.timelines.deleted"]
		if len(deletedFlags) != len(s.Config.Twitter.Timelines) {
			http.Error(w, "Collapsed request", http.StatusInternalServerError)
			return
		}
		actionRetweetCounter := checkboxCounter{"twitter.timelines.action.retweet", 0}
		actionFavoriteCounter := checkboxCounter{"twitter.timelines.action.favorite", 0}
		actionFollowCounter := checkboxCounter{"twitter.timelines.action.follow", 0}
		length := len(val["twitter.timelines.count"])
		timelines := []TimelineConfig{}
		for i := 0; i < length; i++ {
			if deletedFlags[i] == "true" {
				actionRetweetCounter.returnValue(i, val)
				actionFavoriteCounter.returnValue(i, val)
				actionFollowCounter.returnValue(i, val)
				continue
			}
			timeline := *NewTimelineConfig()
			timeline.ScreenNames = getListTextboxValue(val, i, "twitter.timelines.screen_names")
			timeline.ExcludeReplies = getBoolSelectboxValue(val, i, "twitter.timelines.exclude_replies")
			timeline.IncludeRts = getBoolSelectboxValue(val, i, "twitter.timelines.include_rts")
			timeline.Count = atoiOrDefault(val["twitter.timelines.count"][i], timeline.Count)
			timeline.Filter.Patterns = getListTextboxValue(val, i, "twitter.timelines.filter.patterns")
			timeline.Filter.URLPatterns = getListTextboxValue(val, i, "twitter.timelines.filter.url_patterns")
			timeline.Filter.HasMedia = getBoolSelectboxValue(val, i, "twitter.timelines.filter.has_media")
			timeline.Filter.HasURL = getBoolSelectboxValue(val, i, "twitter.timelines.filter.has_url")
			timeline.Filter.Retweeted = getBoolSelectboxValue(val, i, "twitter.timelines.filter.retweeted")
			timeline.Filter.FavoriteThreshold = atoiOrDefault(val["twitter.timelines.filter.favorite_threshold"][i], timeline.Filter.FavoriteThreshold)
			timeline.Filter.RetweetedThreshold = atoiOrDefault(val["twitter.timelines.filter.retweeted_threshold"][i], timeline.Filter.RetweetedThreshold)
			timeline.Filter.Lang = val["twitter.timelines.filter.lang"][i]
			timeline.Filter.Vision.Label = getListTextboxValue(val, i, "twitter.timelines.filter.vision.label")
			timeline.Filter.Vision.Face.AngerLikelihood = val["twitter.timelines.filter.vision.face.anger_likelihood"][i]
			timeline.Filter.Vision.Face.BlurredLikelihood = val["twitter.timelines.filter.vision.face.blurred_likelihood"][i]
			timeline.Filter.Vision.Face.HeadwearLikelihood = val["twitter.timelines.filter.vision.face.headwear_likelihood"][i]
			timeline.Filter.Vision.Face.JoyLikelihood = val["twitter.timelines.filter.vision.face.joy_likelihood"][i]
			timeline.Filter.Vision.Text = getListTextboxValue(val, i, "twitter.timelines.filter.vision.text")
			timeline.Filter.Vision.Landmark = getListTextboxValue(val, i, "twitter.timelines.filter.vision.landmark")
			timeline.Filter.Vision.Logo = getListTextboxValue(val, i, "twitter.timelines.filter.vision.logo")
			timeline.Action.Retweet = actionRetweetCounter.returnValue(i, val)
			timeline.Action.Favorite = actionFavoriteCounter.returnValue(i, val)
			timeline.Action.Follow = actionFollowCounter.returnValue(i, val)
			timeline.Action.Collections = getListTextboxValue(val, i, "twitter.timelines.action.collections")
			timelines = append(timelines, timeline)
		}
		s.Config.Twitter.Timelines = timelines

		deletedFlags = val["twitter.favorites.deleted"]
		if len(deletedFlags) != len(s.Config.Twitter.Favorites) {
			http.Error(w, "Collapsed request", http.StatusInternalServerError)
			return
		}
		actionRetweetCounter = checkboxCounter{"twitter.favorites.action.retweet", 0}
		actionFavoriteCounter = checkboxCounter{"twitter.favorites.action.favorite", 0}
		actionFollowCounter = checkboxCounter{"twitter.favorites.action.follow", 0}
		length = len(val["twitter.favorites.count"])
		favorites := []FavoriteConfig{}
		for i := 0; i < length; i++ {
			if deletedFlags[i] == "true" {
				actionRetweetCounter.returnValue(i, val)
				actionFavoriteCounter.returnValue(i, val)
				actionFollowCounter.returnValue(i, val)
				continue
			}
			favorite := *NewFavoriteConfig()
			favorite.ScreenNames = getListTextboxValue(val, i, "twitter.favorites.screen_names")
			favorite.Count = atoiOrDefault(val["twitter.favorites.count"][i], favorite.Count)
			s.Config.Twitter.Favorites[i] = favorite
			favorite.Filter.Patterns = getListTextboxValue(val, i, "twitter.favorites.filter.patterns")
			favorite.Filter.URLPatterns = getListTextboxValue(val, i, "twitter.favorites.filter.url_patterns")
			favorite.Filter.HasMedia = getBoolSelectboxValue(val, i, "twitter.favorites.filter.has_media")
			favorite.Filter.HasURL = getBoolSelectboxValue(val, i, "twitter.favorites.filter.has_url")
			favorite.Filter.Retweeted = getBoolSelectboxValue(val, i, "twitter.favorites.filter.retweeted")
			favorite.Filter.FavoriteThreshold = atoiOrDefault(val["twitter.favorites.filter.favorite_threshold"][i], favorite.Filter.FavoriteThreshold)
			favorite.Filter.RetweetedThreshold = atoiOrDefault(val["twitter.favorites.filter.retweeted_threshold"][i], favorite.Filter.RetweetedThreshold)
			favorite.Filter.Lang = val["twitter.favorites.filter.lang"][i]
			favorite.Filter.Vision.Label = getListTextboxValue(val, i, "twitter.favorites.filter.vision.label")
			favorite.Filter.Vision.Face.AngerLikelihood = val["twitter.favorites.filter.vision.face.anger_likelihood"][i]
			favorite.Filter.Vision.Face.BlurredLikelihood = val["twitter.favorites.filter.vision.face.blurred_likelihood"][i]
			favorite.Filter.Vision.Face.HeadwearLikelihood = val["twitter.favorites.filter.vision.face.headwear_likelihood"][i]
			favorite.Filter.Vision.Face.JoyLikelihood = val["twitter.favorites.filter.vision.face.joy_likelihood"][i]
			favorite.Filter.Vision.Text = getListTextboxValue(val, i, "twitter.favorites.filter.vision.text")
			favorite.Filter.Vision.Landmark = getListTextboxValue(val, i, "twitter.favorites.filter.vision.landmark")
			favorite.Filter.Vision.Logo = getListTextboxValue(val, i, "twitter.favorites.filter.vision.logo")
			favorite.Action.Retweet = actionRetweetCounter.returnValue(i, val)
			favorite.Action.Favorite = actionFavoriteCounter.returnValue(i, val)
			favorite.Action.Follow = actionFollowCounter.returnValue(i, val)
			favorite.Action.Collections = getListTextboxValue(val, i, "twitter.favorites.action.collections")
			favorites = append(favorites, favorite)
		}
		s.Config.Twitter.Favorites = favorites

		deletedFlags = val["twitter.searches.deleted"]
		if len(deletedFlags) != len(s.Config.Twitter.Searches) {
			http.Error(w, "Collapsed request", http.StatusInternalServerError)
			return
		}
		actionRetweetCounter = checkboxCounter{"twitter.searches.action.retweet", 0}
		actionFavoriteCounter = checkboxCounter{"twitter.searches.action.favorite", 0}
		actionFollowCounter = checkboxCounter{"twitter.searches.action.follow", 0}
		length = len(val["twitter.searches.count"])
		searches := []SearchConfig{}
		for i := 0; i < length; i++ {
			if deletedFlags[i] == "true" {
				actionRetweetCounter.returnValue(i, val)
				actionFavoriteCounter.returnValue(i, val)
				actionFollowCounter.returnValue(i, val)
				continue
			}
			search := *NewSearchConfig()
			search.Queries = getListTextboxValue(val, i, "twitter.searches.queries")
			search.ResultType = val["twitter.searches.result_type"][i]
			search.Count = atoiOrDefault(val["twitter.searches.count"][i], search.Count)
			s.Config.Twitter.Searches[i] = search
			search.Filter.Patterns = getListTextboxValue(val, i, "twitter.searches.filter.patterns")
			search.Filter.URLPatterns = getListTextboxValue(val, i, "twitter.searches.filter.url_patterns")
			search.Filter.HasMedia = getBoolSelectboxValue(val, i, "twitter.searches.filter.has_media")
			search.Filter.HasURL = getBoolSelectboxValue(val, i, "twitter.searches.filter.has_url")
			search.Filter.Retweeted = getBoolSelectboxValue(val, i, "twitter.searches.filter.retweeted")
			search.Filter.FavoriteThreshold = atoiOrDefault(val["twitter.searches.filter.favorite_threshold"][i], search.Filter.FavoriteThreshold)
			search.Filter.RetweetedThreshold = atoiOrDefault(val["twitter.searches.filter.retweeted_threshold"][i], search.Filter.RetweetedThreshold)
			search.Filter.Lang = val["twitter.searches.filter.lang"][i]
			search.Filter.Vision.Label = getListTextboxValue(val, i, "twitter.searches.filter.vision.label")
			search.Filter.Vision.Face.AngerLikelihood = val["twitter.searches.filter.vision.face.anger_likelihood"][i]
			search.Filter.Vision.Face.BlurredLikelihood = val["twitter.searches.filter.vision.face.blurred_likelihood"][i]
			search.Filter.Vision.Face.HeadwearLikelihood = val["twitter.searches.filter.vision.face.headwear_likelihood"][i]
			search.Filter.Vision.Face.JoyLikelihood = val["twitter.searches.filter.vision.face.joy_likelihood"][i]
			search.Filter.Vision.Text = getListTextboxValue(val, i, "twitter.searches.filter.vision.text")
			search.Filter.Vision.Landmark = getListTextboxValue(val, i, "twitter.searches.filter.vision.landmark")
			search.Filter.Vision.Logo = getListTextboxValue(val, i, "twitter.searches.filter.vision.logo")
			search.Action.Retweet = actionRetweetCounter.returnValue(i, val)
			search.Action.Favorite = actionFavoriteCounter.returnValue(i, val)
			search.Action.Follow = actionFollowCounter.returnValue(i, val)
			search.Action.Collections = getListTextboxValue(val, i, "twitter.searches.action.collections")
			searches = append(searches, search)
		}
		s.Config.Twitter.Searches = searches

		s.Config.Twitter.Notification.Place.AllowSelf = len(val["twitter.notification.place.allow_self"]) > 1
		s.Config.Twitter.Notification.Place.Users = getListTextboxValue(val, 0, "twitter.notification.place.users")

		s.Config.DB.Driver = val["db.driver"][0]
		s.Config.DB.DataSource = val["db.data_source"][0]
		s.Config.DB.VisionTable = val["db.vision_table"][0]

		s.Config.Interaction.Duration = val["interaction.duration"][0]
		s.Config.Interaction.AllowSelf = len(val["interaction.allow_self"]) > 1
		s.Config.Interaction.Users = getListTextboxValue(val, 0, "interaction.users")
		s.Config.Interaction.Count = atoiOrDefault(val["interaction.count"][0], s.Config.Interaction.Count)

		s.Config.Log.AllowSelf = len(val["log.allow_self"]) > 1
		s.Config.Log.Users = getListTextboxValue(val, 0, "log.users")

		s.Config.Server.Name = val["server.name"][0]
		s.Config.Server.Host = val["server.host"][0]
		s.Config.Server.Port = val["server.port"][0]
		s.Config.Server.LogLines = atoiOrDefault(val["server.log_lines"][0], s.Config.Server.LogLines)

		err = s.Config.Validate()
		if err != nil {
			s.Config.Load()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.Config.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// This two lines must be in this order, I don't know the reason.
		w.Header().Add("Location", "/config/")
		w.WriteHeader(http.StatusSeeOther)
	} else if r.Method == http.MethodGet {
		tmpl, err := generateTemplate("config", "src/pages/config.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := &struct {
			UserName string
			Config   MybotConfig
		}{
			s.Config.Server.Name,
			*s.Config,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *MybotServer) configTimelineAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		timelines := s.Config.Twitter.Timelines
		timelines = append(timelines, *NewTimelineConfig())
		s.Config.Twitter.Timelines = timelines
		w.Header().Add("Location", "/config/")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *MybotServer) configFavoriteAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		favorites := s.Config.Twitter.Favorites
		favorites = append(favorites, *NewFavoriteConfig())
		s.Config.Twitter.Favorites = favorites
		w.Header().Add("Location", "/config/")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *MybotServer) configSearchAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		searches := s.Config.Twitter.Searches
		searches = append(searches, *NewSearchConfig())
		s.Config.Twitter.Searches = searches
		w.Header().Add("Location", "/config/")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *MybotServer) assetHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/"):]
	data, err := readFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	w.Write(data)
}

func (s *MybotServer) logHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := generateTemplate("log", "src/pages/log.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &struct {
		UserName string
		Log      string
	}{
		s.Config.Server.Name,
		s.Logger.ReadString(),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *MybotServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := generateTemplate("status", "src/pages/status.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &struct {
		UserName string
		Log      string
		Status   MybotStatus
	}{
		s.Config.Server.Name,
		s.Logger.ReadString(),
		*s.Status,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *MybotServer) setupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
			w.Header().Add("Location", "/setup/")
			w.WriteHeader(http.StatusSeeOther)
		}()

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			msg = err.Error()
			return
		}
		val := r.MultipartForm.Value

		ck := val["twitter-consumer-key"][0]
		cs := val["twitter-consumer-secret"][0]
		at := val["twitter-access-token"][0]
		ats := val["twitter-access-token-secret"][0]
		auth := &TwitterAuth{ck, cs, at, ats, s.TwitterAPI.File}

		err = auth.ToJson()
		if err != nil {
			msg = err.Error()
			return
		}

		file, _, err := r.FormFile("gcloud-credential-file")
		if err == nil {
			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				msg = err.Error()
				return
			}
			err = ioutil.WriteFile(s.VisionAPI.File, bytes, 0640)
			if err != nil {
				msg = err.Error()
				return
			}
		}
	} else if r.Method == http.MethodGet {
		tmpl, err := generateTemplate("setup", "src/pages/setup.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msg := ""
		msgCookie, err := r.Cookie("mybot.setup.message")
		if err == nil {
			msg = msgCookie.Value
		}

		data := &struct {
			UserName string
			Message  string
		}{
			s.Config.Server.Name,
			msg,
		}

		if msgCookie != nil {
			msgCookie.Value = ""
			msgCookie.Path = "/setup/"
			http.SetCookie(w, msgCookie)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func generateTemplate(name, path string) (*template.Template, error) {
	index, err := readFile(path)
	if err != nil {
		return nil, err
	}
	header, err := readFile("src/pages/header.html")
	if err != nil {
		return nil, err
	}
	navbar, err := readFile("src/pages/navbar.html")
	if err != nil {
		return nil, err
	}

	funcMap := template.FuncMap{
		"checkbox":      checkbox,
		"boolSelectbox": boolSelectbox,
		"selectbox":     selectbox,
		"listTextbox":   listTextbox,
	}

	return template.
		New("index").
		Funcs(funcMap).
		Parse(string(index) + string(header) + string(navbar))
}

func readFile(path string) ([]byte, error) {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		data, err := Asset(path)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}
