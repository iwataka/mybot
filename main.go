package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

//go:generate go-bindata assets/... pages/...

var (
	twitterAPI *TwitterAPI
	githubAPI  *GitHubAPI
	visionAPI  *VisionAPI
	config     *MybotConfig
	cache      *MybotCache
	logger     *Logger
	status     *MybotStatus
	ctxt       *cli.Context
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	logFlag := cli.StringFlag{
		Name:  "log",
		Value: filepath.Join(home, ".mybot-debug.log"),
		Usage: "Log file's location",
	}

	configFlag := cli.StringFlag{
		Name:  "config",
		Value: "config.toml",
		Usage: "Config file's location",
	}

	cacheFlag := cli.StringFlag{
		Name:  "cache",
		Value: filepath.Join(home, ".cache/mybot/cache.json"),
		Usage: "Cache file's location",
	}

	credFlag := cli.StringFlag{
		Name:  "credential,c",
		Value: "",
		Usage: "Credential for Basic Authentication (ex: user:password)",
	}

	gcloudFlag := cli.StringFlag{
		Name:  "gcloud",
		Value: "google_application_credentials.json",
		Usage: "Credential file for Google Cloud Platform",
	}

	twitterFlag := cli.StringFlag{
		Name:  "twitter",
		Value: "twitter_authentication.json",
		Usage: "Credential file for Twitter API",
	}

	certFlag := cli.StringFlag{
		Name:  "cert",
		Value: "mybot.crt",
		Usage: "Certification file for HTTPS",
	}

	keyFlag := cli.StringFlag{
		Name:  "key",
		Value: "mybot.key",
		Usage: "Key file for HTTPS",
	}

	runFlags := []cli.Flag{
		logFlag,
		configFlag,
		cacheFlag,
		gcloudFlag,
		twitterFlag,
	}

	serveFlags := []cli.Flag{
		logFlag,
		configFlag,
		cacheFlag,
		credFlag,
		gcloudFlag,
		twitterFlag,
		certFlag,
		keyFlag,
	}

	app := cli.NewApp()
	app.Name = "mybot"
	app.Version = "0.1"
	app.Usage = "Automatically collect and broadcast information based on your configuration"
	app.Author = "iwataka"

	runCmd := cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Runs the non-interactive functions only one time (almost for test usage)",
		Flags:   runFlags,
		Before:  beforeRunning,
		Action:  run,
	}

	serveCmd := cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Runs the all functions (both interactive and non-interactive) periodically",
		Flags:   serveFlags,
		Before:  beforeRunning,
		Action:  serve,
	}

	app.Commands = []cli.Command{runCmd, serveCmd}
	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func beforeRunning(c *cli.Context) error {
	var err error
	cache, err = NewMybotCache(c.String("cache"))
	if err != nil {
		panic(err)
	}

	config, err = NewMybotConfig(c.String("config"), nil)
	if err != nil {
		panic(err)
	}

	visionAPI, err = NewVisionAPI(cache, config, c.String("gcloud"))
	if err != nil {
		panic(err)
	}

	config.SetVisionoAPI(visionAPI)

	githubAPI = NewGitHubAPI(nil, cache)

	twitterAuth := &TwitterAuth{}
	twitterAuth.fromJson(c.String("twitter"))
	twitterAPI = NewTwitterAPI(twitterAuth, cache, config)
	if config.Twitter.Debug != nil {
		twitterAPI.SetDebug(*config.Twitter.Debug)
	}

	logger, err = NewLogger(c.String("log"), -1, twitterAPI, config)
	if err != nil {
		panic(err)
	}

	status = &MybotStatus{}

	ctxt = c

	return nil
}

func run(c *cli.Context) error {
	err := runGitHub(c)
	if err != nil {
		logger.Println(err)
	}
	err = runTwitterWithoutStream(c)
	if err != nil {
		logger.Println(err)
	}
	err = cache.Save(c.String("cache"))
	if err != nil {
		logger.Println(err)
	}
	return nil
}

func keepConnection(f func() error, intervalStr string, maxCount int) error {
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return err
	}
	// Abort if there are more than 15 connections in 15 minutes
	t := time.Now()
	count := 0
	for {
		err := f()
		if err != nil {
			logger.Println(err)
		}
		if time.Now().Sub(t) >= interval {
			count = 0
			t = time.Now()
		}
		count++
		if count >= maxCount {
			break
		}
	}
	msg := "Interaction feature is now disabled. Please restart."
	logger.Println(msg)
	return errors.New(msg)
}

func twitterListenMyself(c *cli.Context) {
	if status.TwitterListenMyselfStatus {
		return
	}
	status.TwitterListenMyselfStatus = true
	defer func() { status.TwitterListenMyselfStatus = false }()
	keepConnection(func() error {
		r := twitterAPI.DefaultDirectMessageReceiver
		err := twitterAPI.ListenMyself(nil, r, c.String("cache"))
		if err != nil {
			logger.Println(err)
		}
		return err
	}, "5m", 5)
}

func twitterListenUsers(c *cli.Context) {
	if status.TwitterListenUsersStatus {
		return
	}
	status.TwitterListenUsersStatus = true
	defer func() { status.TwitterListenUsersStatus = false }()
	keepConnection(func() error {
		err := twitterAPI.ListenUsers(nil, c.String("cache"))
		if err != nil {
			logger.Println(err)
		}
		return err
	}, "5m", 5)
}

func githubPeriodically(c *cli.Context) {
	if status.GithubStatus {
		return
	}
	status.GithubStatus = true
	defer func() { status.GithubStatus = false }()
	for {
		err := runGitHub(c)
		if err != nil {
			logger.Println(err)
			return
		}
		err = cache.Save(c.String("cache"))
		if err != nil {
			logger.Println(err)
			return
		}
		d, err := time.ParseDuration(config.GitHub.Duration)
		if err != nil {
			logger.Println(err)
			return
		}
		time.Sleep(d)
	}
}

func twitterPeriodically(c *cli.Context) {
	if status.TwitterStatus {
		return
	}
	status.TwitterStatus = true
	defer func() { status.TwitterStatus = false }()
	for {
		err := runTwitterWithStream(c)
		if err != nil {
			logger.Println(err)
			return
		}
		err = cache.Save(c.String("cache"))
		if err != nil {
			logger.Println(err)
			return
		}
		d, err := time.ParseDuration(config.Twitter.Duration)
		if err != nil {
			logger.Println(err)
			return
		}
		time.Sleep(d)
	}
}

func monitorConfig(c *cli.Context) {
	if status.MonitorConfigStatus {
		return
	}
	status.MonitorConfigStatus = true
	defer func() { status.MonitorConfigStatus = false }()
	monitorFile(
		c.String("config"),
		time.Duration(1)*time.Second,
		func() {
			cfg, err := NewMybotConfig(c.String("config"), visionAPI)
			if err == nil {
				*config = *cfg
			}
		})
}

func httpServer(c *cli.Context) {
	if status.HttpStatus {
		return
	}
	status.HttpStatus = true
	defer func() { status.HttpStatus = false }()
	cred := c.String("credential")
	userAndPassword := strings.SplitN(cred, ":", 2)
	user := ""
	password := ""
	if len(userAndPassword) == 2 {
		user = userAndPassword[0]
		password = userAndPassword[1]
	}
	cert := c.String("cert")
	key := c.String("key")
	err := config.HTTP.Init(user, password, cert, key)
	if err != nil {
		panic(err)
	}
}

func serve(c *cli.Context) error {
	s := config.HTTP
	s.Logger = logger
	s.TwitterAPI = twitterAPI
	s.VisionAPI = visionAPI
	s.cache = cache
	s.config = config
	s.status = status

	go twitterListenMyself(c)
	go twitterListenUsers(c)
	go githubPeriodically(c)
	go twitterPeriodically(c)
	go monitorConfig(c)
	go httpServer(c)

	ch := make(chan bool)
	<-ch
	return nil
}

func monitorFile(file string, d time.Duration, f func()) {
	info, _ := os.Stat(file)
	modTime := time.Now()
	if info != nil {
		modTime = info.ModTime()
	}
	for {
		info, _ := os.Stat(file)
		if info != nil {
			mt := info.ModTime()
			if mt.After(modTime) {
				modTime = mt
				f()
			}
		}
		time.Sleep(d)
	}
}

func runGitHub(c *cli.Context) error {
	for _, p := range config.GitHub.Projects {
		err := githubCommitTweet(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func runTwitterWithStream(c *cli.Context) error {
	tweets := []anaconda.Tweet{}
	for _, a := range config.Twitter.Searches {
		a.Filter.visionAPI = visionAPI
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if len(a.ResultType) != 0 {
			v.Set("result_type", a.ResultType)
		}
		for _, query := range a.Queries {
			ts, err := twitterAPI.DoForSearch(query, v, a.Filter, a.Action)
			if err != nil {
				return err
			}
			tweets = append(tweets, ts...)
		}
	}
	for _, a := range config.Twitter.Favorites {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		a.Filter.visionAPI = visionAPI
		for _, name := range a.ScreenNames {
			ts, err := twitterAPI.DoForFavorites(name, v, a.Filter, a.Action)
			tweets = append(tweets, ts...)
			if err != nil {
				return err
			}
		}
	}
	for _, t := range tweets {
		err := twitterAPI.NotifyToAll(&t)
		if err != nil {
			return err
		}
	}
	return nil
}

func runTwitterWithoutStream(c *cli.Context) error {
	err := runTwitterWithStream(c)
	if err != nil {
		return err
	}
	tweets := []anaconda.Tweet{}
	for _, a := range config.Twitter.Timelines {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if a.ExcludeReplies != nil {
			v.Set("exclude_replies", fmt.Sprintf("%v", *a.ExcludeReplies))
		}
		if a.IncludeRts != nil {
			v.Set("include_rts", fmt.Sprintf("%v", *a.IncludeRts))
		}
		a.Filter.visionAPI = visionAPI
		for _, name := range a.ScreenNames {
			ts, err := twitterAPI.DoForAccount(name, v, a.Filter, a.Action)
			tweets = append(tweets, ts...)
			if err != nil {
				return err
			}
		}
	}
	for _, t := range tweets {
		err := twitterAPI.NotifyToAll(&t)
		if err != nil {
			return err
		}
	}
	return nil
}

func githubCommitTweet(p GitHubProject) error {
	commit, err := githubAPI.LatestCommit(p)
	if err != nil {
		return err
	}
	if commit != nil {
		msg := p.User + "/" + p.Repo + "\n" + *commit.HTMLURL
		_, err := twitterAPI.PostTweet(msg, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
