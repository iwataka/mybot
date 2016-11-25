package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

//go:generate go-bindata assets/... index.html

var (
	twitterAPI *TwitterAPI
	githubAPI  *GitHubAPI
	visionAPI  *VisionAPI
	config     *MybotConfig
	cache      *MybotCache
	logger     *Logger
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

	runFlags := []cli.Flag{
		logFlag,
		configFlag,
		cacheFlag,
	}

	serveFlags := []cli.Flag{
		logFlag,
		configFlag,
		cacheFlag,
		credFlag,
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

	visionAPI, err = NewVisionAPI(cache)
	if err != nil {
		panic(err)
	}

	config, err = NewMybotConfig(c.String("config"), visionAPI)
	if err != nil {
		panic(err)
	}

	githubAPI = NewGitHubAPI(nil, cache)
	twitterAPI = NewTwitterAPI(config.Authentication, cache, config)
	ok, err := twitterAPI.api.VerifyCredentials()
	if err != nil {
		panic(err)
	}
	if !ok {
		panic("Twitter authentication failed")
	}
	if config.Twitter.Debug != nil {
		twitterAPI.SetDebug(*config.Twitter.Debug)
	}

	logger, err = NewLogger(c.String("log"), -1, twitterAPI, config)
	if err != nil {
		panic(err)
	}

	return nil
}

func run(c *cli.Context) error {
	runGitHub(c, logger.HandleError)
	runTwitterWithoutStream(c, logger.HandleError)
	err := cache.Save(c.String("cache"))
	if err != nil {
		logger.Println(err)
	}
	return nil
}

func keepConnection(f func() error) {
	// Abort if there are more than 15 connections in 15 minutes
	t := time.Now()
	count := 0
	for {
		err := f()
		if err != nil {
			logger.Println(err)
		}
		if time.Now().Sub(t) >= 15*time.Minute {
			count = 0
			t = time.Now()
		}
		count++
		if count >= 15 {
			break
		}
	}
	logger.Println("Interaction feature is now disabled. Please restart.")
}

func serve(c *cli.Context) error {
	s := config.HTTP
	s.Logger = logger
	s.TwitterAPI = twitterAPI
	s.VisionAPI = visionAPI
	s.cache = cache
	ch := make(chan bool)

	go keepConnection(func() error {
		r := twitterAPI.DefaultDirectMessageReceiver
		return twitterAPI.ListenMyself(nil, r, c.String("cache"))
	})

	go keepConnection(func() error {
		return twitterAPI.ListenUsers(nil, c.String("cache"))
	})

	go func() {
		for {
			runGitHub(c, logger.HandleError)
			err := cache.Save(c.String("cache"))
			if err != nil {
				logger.Println(err)
			}
			d, err := time.ParseDuration(config.GitHub.Duration)
			if err != nil {
				logger.Println(err)
				panic(err)
			}
			time.Sleep(d)
		}
	}()

	go func() {
		for {
			runTwitterWithStream(c, logger.HandleError)
			err := cache.Save(c.String("cache"))
			if err != nil {
				logger.Println(err)
			}
			d, err := time.ParseDuration(config.Twitter.Duration)
			if err != nil {
				logger.Println(err)
				panic(err)
			}
			time.Sleep(d)
		}
	}()

	go monitorFile(
		c.String("config"),
		time.Duration(1)*time.Second,
		func() {
			cfg, err := NewMybotConfig(c.String("config"), visionAPI)
			if err == nil {
				if !reflect.DeepEqual(cfg.Authentication, config.Authentication) {
					*twitterAPI = *NewTwitterAPI(cfg.Authentication, cache, config)
				}
				*config = *cfg
			}
		})

	go func() {
		cred := c.String("credential")
		userAndPassword := strings.SplitN(cred, ":", 2)
		user := ""
		password := ""
		if len(userAndPassword) == 2 {
			user = userAndPassword[0]
			password = userAndPassword[1]
		}
		err := s.Init(user, password)
		if err != nil {
			panic(err)
		}
	}()

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

func runGitHub(c *cli.Context, handle func(error)) {
	for _, p := range config.GitHub.Projects {
		handle(githubCommitTweet(p))
	}
}

func runTwitterWithStream(c *cli.Context, handle func(error)) {
	tweets := []anaconda.Tweet{}
	for _, a := range config.Twitter.Searches {
		a.Filter.visionAPI = visionAPI
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if a.ResultType != nil {
			v.Set("result_type", *a.ResultType)
		}
		if a.Query != nil {
			ts, err := twitterAPI.DoForSearch(*a.Query, v, a.Filter, a.Action)
			handle(err)
			tweets = append(tweets, ts...)
		} else {
			for _, query := range a.Queries {
				ts, err := twitterAPI.DoForSearch(query, v, a.Filter, a.Action)
				handle(err)
				tweets = append(tweets, ts...)
			}
		}
	}
	for _, t := range tweets {
		err := twitterAPI.NotifyToAll(&t)
		handle(err)
	}
}

func runTwitterWithoutStream(c *cli.Context, handle func(error)) {
	runTwitterWithStream(c, handle)
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
		if a.ScreenName != nil {
			ts, err := twitterAPI.DoForAccount(*a.ScreenName, v, a.Filter, a.Action)
			tweets = append(tweets, ts...)
			handle(err)
		} else {
			for _, name := range a.ScreenNames {
				ts, err := twitterAPI.DoForAccount(name, v, a.Filter, a.Action)
				tweets = append(tweets, ts...)
				handle(err)
			}
		}
	}
	for _, a := range config.Twitter.Favorites {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		a.Filter.visionAPI = visionAPI
		if a.ScreenName != nil {
			ts, err := twitterAPI.DoForFavorites(*a.ScreenName, v, a.Filter, a.Action)
			tweets = append(tweets, ts...)
			handle(err)
		} else {
			for _, name := range a.ScreenNames {
				ts, err := twitterAPI.DoForFavorites(name, v, a.Filter, a.Action)
				tweets = append(tweets, ts...)
				handle(err)
			}
		}
	}
	for _, t := range tweets {
		err := twitterAPI.NotifyToAll(&t)
		handle(err)
	}
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
