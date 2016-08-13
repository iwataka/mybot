package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type HTTPServer struct {
	Name       string
	Port       string
	Logger     *MultiLogger
	TwitterAPI *TwitterAPI
}

func (s *HTTPServer) initHTTP() error {
	http.HandleFunc("/", s.handler)
	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))
	err := http.ListenAndServe(":"+s.Port, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) handler(w http.ResponseWriter, r *http.Request) {
	index, err := ioutil.ReadFile("index.html")
	tmpl, err := template.New("index").Parse(string(index))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	log := s.Logger.ReadString()
	self, err := s.TwitterAPI.GetSelfCache()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	botName := self.ScreenName

	data := &struct {
		UserName string
		Log      string
		BotName  string
	}{
		s.Name,
		log,
		botName,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
}
