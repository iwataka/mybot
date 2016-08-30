package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type HTTPServer struct {
	Name       string      `toml:"name"`
	Host       string      `toml:"host"`
	Port       string      `toml:"port"`
	Enabled    bool        `toml:"enabled"`
	Logger     *Logger     `toml:"-"`
	TwitterAPI *TwitterAPI `toml:"-"`
	VisionAPI  *VisionAPI  `toml:"-"`
	cache      *MybotCache `toml:"-"`
}

func (s *HTTPServer) Init() error {
	if s.Enabled {
		fmt.Printf("Open %s:%s for more details\n", s.Host, s.Port)
		http.HandleFunc("/", s.handler)
		http.HandleFunc("/assets/", s.customCSSHandler)
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
		var botName string
		self, err := s.TwitterAPI.GetSelf()
		if err == nil {
			botName = self.ScreenName
		} else {
			botName = ""
		}

		pid := ""
		if s.VisionAPI != nil {
			pid = s.VisionAPI.ProjectID
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
		colList, err := s.TwitterAPI.GetCollectionListByUserId(self.Id, nil)
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
			GCloudProjectID     string
			ImageURL            string
			ImageAnalysisResult string
			ImageAnalysisDate   string
			CollectionMap       map[string]string
		}{
			s.Name,
			log,
			botName,
			pid,
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

func (s *HTTPServer) customCSSHandler(w http.ResponseWriter, r *http.Request) {
	data, err := Asset(r.URL.Path[1:])
	if err != nil {
		s.Logger.Println(err)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	w.Write(data)
}
