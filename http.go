package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
)

type HTTPServer struct {
	Name       string
	Port       string
	Logger     *Logger
	TwitterAPI *TwitterAPI
	VisionAPI  *VisionAPI
	cache      *MybotCache
}

func (s *HTTPServer) Init() error {
	http.HandleFunc("/", s.handler)
	http.HandleFunc("/assets/css/custom.css", s.customCSSHandler)
	http.HandleFunc("/404", s.notFoundHandler)
	err := http.ListenAndServe(":"+s.Port, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		index, err := Asset("index.html")
		tmpl, err := template.New("index").Parse(string(index))
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		log := s.Logger.ReadString()
		var botName string
		self, err := s.TwitterAPI.GetSelfCache()
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

		data := &struct {
			UserName            string
			Log                 string
			BotName             string
			GCloudProjectID     string
			ImageURL            string
			ImageAnalysisResult string
			ImageAnalysisDate   string
		}{
			s.Name,
			log,
			botName,
			pid,
			s.cache.ImageURL,
			imageAnalysisResult,
			s.cache.ImageAnalysisDate,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}
	} else {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}
}

func (s *HTTPServer) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	data, err := Asset("404.html")
	if err != nil {
		s.Logger.InfoIfError(err)
		return
	}
	w.Write(data)
}

func (s *HTTPServer) customCSSHandler(w http.ResponseWriter, r *http.Request) {
	data, err := Asset("assets/css/custom.css")
	if err != nil {
		s.Logger.InfoIfError(err)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	w.Write(data)
}
