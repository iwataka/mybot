package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func initHttp() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/log", logHandler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadFile("config.html")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	tmpl, err := template.New("config").Parse(string(bytes))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	err = tmpl.Execute(w, config)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadFile(logFile)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
}
