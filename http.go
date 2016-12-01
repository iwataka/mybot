package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// HTTPServer contains values for providing various pieces of information, such
// as the error log and Google Vision API result, via HTTP to users.
type HTTPServer struct {
	Name       string       `toml:"name"`
	Host       string       `toml:"host"`
	Port       string       `toml:"port"`
	Enabled    bool         `toml:"enabled"`
	LogLines   *int         `toml:logLines`
	Logger     *Logger      `toml:"-"`
	TwitterAPI *TwitterAPI  `toml:"-"`
	VisionAPI  *VisionAPI   `toml:"-"`
	cache      *MybotCache  `toml:"-"`
	config     *MybotConfig `toml:"-"`
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
		http.HandleFunc("/", wrapHandlerWithBasicAuth(s.handler, user, password))
		http.HandleFunc("/config/", wrapHandlerWithBasicAuth(s.configHandler, user, password))
		http.HandleFunc("/assets/", wrapHandlerWithBasicAuth(s.assetHandler, user, password))
		http.HandleFunc("/log/", wrapHandlerWithBasicAuth(s.logHandler, user, password))
		http.HandleFunc("/api/config/", wrapHandlerWithBasicAuth(s.apiConfigHandler, user, password))
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

func (s *HTTPServer) configHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := generateTemplate("config", "pages/config.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &struct {
		UserName string
	}{
		s.Name,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) assetHandler(w http.ResponseWriter, r *http.Request) {
	data, err := Asset(r.URL.Path[len("/"):])
	if err != nil {
		s.Logger.Println(err)
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
	}
}

func generateTemplate(name, path string) (*template.Template, error) {
	index, err := Asset(path)
	if err != nil {
		return nil, err
	}
	header, err := Asset("pages/header.html")
	if err != nil {
		return nil, err
	}
	navbar, err := Asset("pages/navbar.html")
	if err != nil {
		return nil, err
	}
	return template.New("index").Parse(string(index) + string(header) + string(navbar))
}
