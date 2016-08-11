package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type homeData struct {
	Title     string
	Header    string
	Log       string
	LogExists bool
	BotName   string
}

func initHttp() {
	http.HandleFunc("/", handler)
	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))
	http.ListenAndServe(":"+config.Option.Port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	index, err := ioutil.ReadFile("index.html")
	tmpl, err := template.New("index").Parse(string(index))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	data, err := newHomeData()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
}

func newHomeData() (*homeData, error) {
	var log string
	var logExists bool
	if info, _ := os.Stat(logFile); info == nil || info.IsDir() {
		logExists = false
	} else {
		bytes, err := ioutil.ReadFile(logFile)
		if err != nil {
			return nil, err
		}
		log = string(bytes)
		if log == "" {
			logExists = false
		} else {
			logExists = true
		}
	}
	var title string
	var header string
	name := config.Option.Name
	if name == "" {
		title = "Mybot Home"
		header = "Mybot"
	} else {
		title = name + "'s bot Home"
		header = name + "'s bot"
	}
	botName, err := getTwitterSelf()
	if err != nil {
		return nil, err
	}
	return &homeData{title, header, log, logExists, botName}, nil
}
