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

// HTTPServer contains values for providing various pieces of information, such
// as the error log and Google Vision API result, via HTTP to users.
type HTTPServer struct {
	Name       string       `toml:"name"`
	Host       string       `toml:"host"`
	Port       string       `toml:"port"`
	Enabled    bool         `toml:"enabled"`
	LogLines   *int         `toml:log_lines`
	Logger     *Logger      `toml:"-"`
	TwitterAPI *TwitterAPI  `toml:"-"`
	VisionAPI  *VisionAPI   `toml:"-"`
	cache      *MybotCache  `toml:"-"`
	config     *MybotConfig `toml:"-"`
	status     *MybotStatus `toml:"-"`
}

type httpHandler func(http.ResponseWriter, *http.Request)

func checkAuth(r *http.Request, user, password string) bool {
	u, p, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return u == user && p == password
}

func wrapHandlerWithBasicAuth(f httpHandler, user, password string) httpHandler {
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

// Init initializes HTTP server if HTTPServer#Enabled is true.
func (s *HTTPServer) Init(user, password, cert, key string) error {
	if s.Enabled {
		http.HandleFunc(
			"/",
			wrapHandlerWithBasicAuth(s.handler, user, password),
		)
		http.HandleFunc(
			"/config/",
			wrapHandlerWithBasicAuth(s.configHandler, user, password),
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
		addr := s.Host + ":" + s.Port
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

func (s *HTTPServer) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		tmpl, err := generateTemplate("index", "pages/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log := s.Logger.ReadString()
		lines := strings.Split(log, "\n")
		lineNum := 10
		if s.LogLines != nil {
			lineNum = *s.LogLines
		}
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
			s.Name,
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

func (c *checkboxCounter) returnValue(index int, val map[string][]string) *bool {
	if val[c.name][index+c.extraCount] == "true" {
		c.extraCount++
		b := true
		return &b
	} else {
		b := false
		return &b
	}
}

func atoiOrNil(str string) *int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return nil
	} else {
		return &i
	}
}

func (s *HTTPServer) configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		val := r.MultipartForm.Value

		excludeRepliesCounter := checkboxCounter{"twitter.timelines.exclude_replies", 0}
		includeRtsCounter := checkboxCounter{"twitter.timelines.include_rts", 0}
		for i, _ := range s.config.Twitter.Timelines {
			timeline := s.config.Twitter.Timelines[i]
			timeline.ScreenNames = strings.Split(val["twitter.timelines.screen_names"][i], ",")
			timeline.ExcludeReplies = excludeRepliesCounter.returnValue(i, val)
			timeline.IncludeRts = includeRtsCounter.returnValue(i, val)
			timeline.Count = atoiOrNil(val["twitter.timelines.count"][i])
			s.config.Twitter.Timelines[i] = timeline
		}

		for i, _ := range s.config.Twitter.Favorites {
			favorite := s.config.Twitter.Favorites[i]
			favorite.ScreenNames = strings.Split(val["twitter.favorites.screen_names"][i], ",")
			favorite.Count = atoiOrNil(val["twitter.favorites.count"][i])
			s.config.Twitter.Favorites[i] = favorite
		}

		for i, _ := range s.config.Twitter.Searches {
			search := s.config.Twitter.Searches[i]
			search.Queries = strings.Split(val["twitter.searches.queries"][i], ",")
			search.ResultType = val["twitter.searches.result_type"][i]
			search.Count = atoiOrNil(val["twitter.searches.count"][i])
			s.config.Twitter.Searches[i] = search
		}

		s.config.DB.Driver = val["db.driver"][0]
		s.config.DB.DataSource = val["db.data_source"][0]
		s.config.DB.VisionTable = val["db.vision_table"][0]

		s.config.Interaction.Duration = val["interaction.duration"][0]
		s.config.Interaction.AllowSelf = len(val["interaction.allow_self"]) > 1
		s.config.Interaction.Users = strings.Split(val["interaction.users"][0], ",")
		s.config.Interaction.Count = atoiOrNil(val["interaction.count"][0])

		s.config.Log.AllowSelf = len(val["log.allow_self"]) > 1
		s.config.Log.Users = strings.Split(val["log.users"][0], ",")

		s.config.HTTP.Name = val["http.name"][0]
		s.config.HTTP.Host = val["http.host"][0]
		s.config.HTTP.Port = val["http.port"][0]
		s.config.HTTP.Enabled = len(val["http.enabled"]) > 1
		s.config.HTTP.LogLines = atoiOrNil(val["http.log_lines"][0])

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
	}
	tmpl, err := generateTemplate("config", "pages/config.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &struct {
		UserName string
		Config   MybotConfig
	}{
		s.Name,
		*s.config,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) assetHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/"):]
	data, err := readFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	w.Write(data)
}

func (s *HTTPServer) logHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := generateTemplate("log", "pages/log.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &struct {
		UserName string
		Log      string
	}{
		s.Name,
		logger.ReadString(),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) statusHandler(w http.ResponseWriter, r *http.Request) {
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
		s.Name,
		logger.ReadString(),
		*s.status,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) apiConfigHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *HTTPServer) apiTwitterListenMyselfHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.TwitterListenMyselfStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { twitterListenMyself(ctxt) })
	}
}

func (s *HTTPServer) apiTwitterListenUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.TwitterListenUsersStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { twitterListenUsers(ctxt) })
	}
}

func (s *HTTPServer) apiGitHubHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.GithubStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { githubPeriodically(ctxt) })
	}
}

func (s *HTTPServer) apiTwitterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.TwitterStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { twitterPeriodically(ctxt) })
	}
}

func (s *HTTPServer) apiMonitorConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.apiGetStatusHandler(w, r, s.status.MonitorConfigStatus)
	} else if r.Method == http.MethodPost {
		s.apiPostStatusHandler(w, r, func() { monitorConfig(ctxt) })
	}
}

func (s *HTTPServer) apiGetStatusHandler(w http.ResponseWriter, r *http.Request, status bool) {
	data := &APIFeatureStatus{status}
	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func (s *HTTPServer) apiPostStatusHandler(w http.ResponseWriter, r *http.Request, f func()) {
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
		"convertListToShow": convertListToShow,
		"checkBoolRef":      checkBoolRef,
		"derefString":       derefString,
		"checkbox":          checkbox,
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
