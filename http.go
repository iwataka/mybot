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
)

// MybotServer contains values for providing various pieces of information, such
// as the error log and Google Vision API result, via HTTP to users.
type MybotServer struct {
	Logger     *Logger      `toml:"-"`
	TwitterAPI *TwitterAPI  `toml:"-"`
	VisionAPI  *VisionAPI   `toml:"-"`
	cache      *MybotCache  `toml:"-"`
	config     *MybotConfig `toml:"-"`
	status     *MybotStatus `toml:"-"`
}

func checkAuth(r *http.Request, user, password string) bool {
	u, p, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return u == user && p == password
}

func wrapHandlerWithBasicAuth(f http.HandlerFunc, user, password string) http.HandlerFunc {
	if len(user) != 0 && len(password) != 0 {
		return func(w http.ResponseWriter, r *http.Request) {
			if !checkAuth(r, user, password) {
				w.Header().Set("WWW-Authenticate", `Basic realm="Enter username and password"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				return
			}
			f(w, r)
		}
	} else {
		return f
	}
}

// Init initializes HTTP server if MybotServer#Enabled is true.
func (s *MybotServer) Init(user, password, cert, key string) error {
	if s.config.HTTP.Enabled {
		http.HandleFunc(
			"/",
			wrapHandlerWithBasicAuth(s.handler, user, password),
		)
		http.HandleFunc(
			"/config/",
			wrapHandlerWithBasicAuth(s.configHandler, user, password),
		)
		http.HandleFunc(
			"/config/timelines/add",
			wrapHandlerWithBasicAuth(s.configTimelineAddHandler, user, password),
		)
		http.HandleFunc(
			"/config/favorites/add",
			wrapHandlerWithBasicAuth(s.configFavoriteAddHandler, user, password),
		)
		http.HandleFunc(
			"/config/searches/add",
			wrapHandlerWithBasicAuth(s.configSearchAddHandler, user, password),
		)
		http.HandleFunc(
			"/assets/",
			wrapHandlerWithBasicAuth(s.assetHandler, user, password),
		)
		http.HandleFunc(
			"/log/",
			wrapHandlerWithBasicAuth(s.logHandler, user, password),
		)
		http.HandleFunc(
			"/status/",
			wrapHandlerWithBasicAuth(s.statusHandler, user, password),
		)

		// API handlers
		http.HandleFunc(
			"/api/config/",
			wrapHandlerWithBasicAuth(s.apiConfigHandler, user, password),
		)
		http.HandleFunc(
			"/api/features/listen/myself/",
			wrapHandlerWithBasicAuth(s.apiTwitterListenMyselfHandler, user, password),
		)
		http.HandleFunc(
			"/api/features/listen/users/",
			wrapHandlerWithBasicAuth(s.apiTwitterListenUsersHandler, user, password),
		)
		http.HandleFunc(
			"/api/features/github/periodic/",
			wrapHandlerWithBasicAuth(s.apiGitHubHandler, user, password),
		)
		http.HandleFunc(
			"/api/features/twitter/periodic/",
			wrapHandlerWithBasicAuth(s.apiTwitterHandler, user, password),
		)
		http.HandleFunc(
			"/api/features/monitor/config/",
			wrapHandlerWithBasicAuth(s.apiMonitorConfigHandler, user, password),
		)

		var err error
		addr := s.config.HTTP.Host + ":" + s.config.HTTP.Port
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
	}
	return nil
}

func (s *MybotServer) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		tmpl, err := generateTemplate("index", "pages/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log := s.Logger.ReadString()
		lines := strings.Split(log, "\n")
		lineNum := s.config.HTTP.LogLines
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
		if s.cache != nil {
			buf := new(bytes.Buffer)
			err := json.Indent(buf, []byte(s.cache.ImageAnalysisResult), "", "  ")
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
			s.config.HTTP.Name,
			log,
			botName,
			s.cache.ImageURL,
			s.cache.ImageSource,
			imageAnalysisResult,
			s.cache.ImageAnalysisDate,
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
		if len(deletedFlags) != len(s.config.Twitter.Timelines) {
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
		s.config.Twitter.Timelines = timelines

		deletedFlags = val["twitter.favorites.deleted"]
		if len(deletedFlags) != len(s.config.Twitter.Favorites) {
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
			s.config.Twitter.Favorites[i] = favorite
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
		s.config.Twitter.Favorites = favorites

		deletedFlags = val["twitter.searches.deleted"]
		if len(deletedFlags) != len(s.config.Twitter.Searches) {
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
			s.config.Twitter.Searches[i] = search
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
		s.config.Twitter.Searches = searches

		s.config.Twitter.Notification.Place.AllowSelf = len(val["twitter.notification.place.allow_self"]) > 1
		s.config.Twitter.Notification.Place.Users = getListTextboxValue(val, 0, "twitter.notification.place.users")

		s.config.DB.Driver = val["db.driver"][0]
		s.config.DB.DataSource = val["db.data_source"][0]
		s.config.DB.VisionTable = val["db.vision_table"][0]

		s.config.Interaction.Duration = val["interaction.duration"][0]
		s.config.Interaction.AllowSelf = len(val["interaction.allow_self"]) > 1
		s.config.Interaction.Users = getListTextboxValue(val, 0, "interaction.users")
		s.config.Interaction.Count = atoiOrDefault(val["interaction.count"][0], s.config.Interaction.Count)

		s.config.Log.AllowSelf = len(val["log.allow_self"]) > 1
		s.config.Log.Users = getListTextboxValue(val, 0, "log.users")

		s.config.HTTP.Name = val["http.name"][0]
		s.config.HTTP.Host = val["http.host"][0]
		s.config.HTTP.Port = val["http.port"][0]
		s.config.HTTP.Enabled = len(val["http.enabled"]) > 1
		s.config.HTTP.LogLines = atoiOrDefault(val["http.log_lines"][0], s.config.HTTP.LogLines)

		err = ValidateConfig(s.config)
		if err != nil {
			s.config.Reload()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.config.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// This two lines must be in this order, I don't know the reason.
		w.Header().Add("Location", "/config/")
		w.WriteHeader(http.StatusSeeOther)
	} else if r.Method == http.MethodGet {
		tmpl, err := generateTemplate("config", "pages/config.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := &struct {
			UserName string
			Config   MybotConfig
		}{
			s.config.HTTP.Name,
			*s.config,
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
		timelines := s.config.Twitter.Timelines
		timelines = append(timelines, *NewTimelineConfig())
		s.config.Twitter.Timelines = timelines
		w.Header().Add("Location", "/config/")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *MybotServer) configFavoriteAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		favorites := s.config.Twitter.Favorites
		favorites = append(favorites, *NewFavoriteConfig())
		s.config.Twitter.Favorites = favorites
		w.Header().Add("Location", "/config/")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *MybotServer) configSearchAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		searches := s.config.Twitter.Searches
		searches = append(searches, *NewSearchConfig())
		s.config.Twitter.Searches = searches
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
	tmpl, err := generateTemplate("log", "pages/log.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &struct {
		UserName string
		Log      string
	}{
		s.config.HTTP.Name,
		logger.ReadString(),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *MybotServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := generateTemplate("status", "pages/status.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &struct {
		UserName string
		Log      string
		Status   MybotStatus
	}{
		s.config.HTTP.Name,
		logger.ReadString(),
		*s.status,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *MybotServer) apiConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/json")
		bytes, err := json.Marshal(s.config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(bytes, s.config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = s.config.Save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type APIFeatureStatus struct {
	Status bool `json:status`
}

func (s *MybotServer) apiTwitterListenMyselfHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.TwitterListenMyselfStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { twitterListenMyself() })
	}
}

func (s *MybotServer) apiTwitterListenUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.TwitterListenUsersStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { twitterListenUsers() })
	}
}

func (s *MybotServer) apiGitHubHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.GithubStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { githubPeriodically() })
	}
}

func (s *MybotServer) apiTwitterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.TwitterStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { twitterPeriodically() })
	}
}

func (s *MybotServer) apiMonitorConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.MonitorConfigStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { monitorConfig() })
	}
}

func (s *MybotServer) apiGetStatusHandler(w http.ResponseWriter, r *http.Request, status bool) {
	data := &APIFeatureStatus{status}
	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func (s *MybotServer) apiPostStatusHandler(w http.ResponseWriter, r *http.Request, f func()) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &APIFeatureStatus{}
	err = json.Unmarshal(bytes, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data.Status {
		f()
	}
}

func generateTemplate(name, path string) (*template.Template, error) {
	index, err := readFile(path)
	if err != nil {
		return nil, err
	}
	header, err := readFile("pages/header.html")
	if err != nil {
		return nil, err
	}
	navbar, err := readFile("pages/navbar.html")
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
