package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// HTTPServer contains values for providing various pieces of information, such
// as the error log and Google Vision API result, via HTTP to users.
type HTTPServer struct {
	Name       string      `toml:"name"`
	Host       string      `toml:"host"`
	Port       string      `toml:"port"`
	Enabled    bool        `toml:"enabled"`
	LogLines   *int        `toml:logLines`
	Logger     *Logger     `toml:"-"`
	TwitterAPI *TwitterAPI `toml:"-"`
	VisionAPI  *VisionAPI  `toml:"-"`
	cache      *MybotCache `toml:"-"`
}

// Init initializes HTTP server if HTTPServer#Enabled is true.
func (s *HTTPServer) Init() error {
	if s.Enabled {
		fmt.Printf("Open %s:%s for more details\n", s.Host, s.Port)
		http.HandleFunc("/", s.handler)
		http.HandleFunc("/assets/", s.assetHandler)
		http.HandleFunc("/log/", s.logHandler)
		err := http.ListenAndServe(s.Host+":"+s.Port, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *HTTPServer) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		index, err := Asset("index.html")
		tmpl, err := template.New("index").Parse(string(index))
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
			ImageAnalysisResult string
			ImageAnalysisDate   string
			CollectionMap       map[string]string
		}{
			s.Name,
			log,
			botName,
			s.cache.ImageURL,
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

func (s *HTTPServer) assetHandler(w http.ResponseWriter, r *http.Request) {
	data, err := Asset(r.URL.Path[1:])
	if err != nil {
		s.Logger.Println(err)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	w.Write(data)
}

func (s *HTTPServer) logHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.ReadString()
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(log))
}
