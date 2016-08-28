package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

//go:generate go-bindata assets/... index.html 404.html

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

	visionCredentialFlag := cli.StringFlag{
		Name:  "gcp-credential",
		Value: "credential.json",
		Usage: "Location of Google Cloud Platform credential file",
	}

	flags := []cli.Flag{
		logFlag,
		configFlag,
		cacheFlag,
		visionCredentialFlag,
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
		Flags:   flags,
		Before:  beforeRunning,
		Action:  run,
	}

	serveCmd := cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Runs the all functions (both interactive and non-interactive) periodically",
		Flags:   flags,
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
	config, err = NewMybotConfig(c.String("config"))
	if err != nil {
		panic(err)
	}

	githubAPI = NewGitHubAPI(nil, cache)
	twitterAPI = NewTwitterAPI(config.Authentication, cache)

	logger, err = NewLogger(c.String("log"), -1, twitterAPI, config)
	if err != nil {
		panic(err)
	}

	// visionAPI is nil if there exists no credential file
	visionAPI, err = NewVisionAPI(c.String("gcp-credential"), cache)
	if err != nil {
		logger.InfoIfError(err)
		visionAPI = new(VisionAPI)
	}

	return nil
}

func run(c *cli.Context) error {
	runGitHub(c, logger.InfoIfError)
	runRetweet(c, logger.InfoIfError)
	logger.InfoIfError(cache.Save(c.String("cache")))
	return nil
}

func serve(c *cli.Context) error {
	s := config.HTTP
	s.Logger = logger
	s.TwitterAPI = twitterAPI
	s.VisionAPI = visionAPI
	s.cache = cache

	go func() {
		for {
			if config.Interaction != nil {
				logger.InfoIfError(twitterAPI.Response(config.Interaction.Users))
			}
			d, err := time.ParseDuration(config.Interaction.Duration)
			logger.FatalIfError(err)
			time.Sleep(d)
		}
	}()

	go func() {
		for {
			runGitHub(c, logger.InfoIfError)
			logger.InfoIfError(cache.Save(c.String("cache")))
			d, err := time.ParseDuration(config.GitHub.Duration)
			logger.FatalIfError(err)
			time.Sleep(d)
		}
	}()

	go func() {
		for {
			runRetweet(c, logger.InfoIfError)
			logger.InfoIfError(cache.Save(c.String("cache")))
			d, err := time.ParseDuration(config.Twitter.Duration)
			logger.FatalIfError(err)
			time.Sleep(d)
		}
	}()

	go monitorFile(
		c.String("config"),
		time.Duration(1)*time.Second,
		func() {
			cfg, err := NewMybotConfig(c.String("config"))
			if err == nil {
				if !reflect.DeepEqual(cfg.Authentication, config.Authentication) {
					*twitterAPI = *NewTwitterAPI(cfg.Authentication, cache)
				}
				*config = *cfg
			}
		})

	go monitorFile(
		c.String("gcp-credential"),
		time.Duration(1)*time.Second,
		// If it fails to read a credential file, it may be better to
		// execute `*visionAPI = new(VisionAPI)`
		func() {
			a, err := NewVisionAPI(c.String("gcp-credential"), cache)
			logger.InfoIfError(err)
			if err == nil {
				*visionAPI = *a
			}

		})

	fmt.Printf("Open %s:%s for detailed information\n", s.Host, s.Port)
	err := s.Init()
	if err != nil {
		panic(err)
	}

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
			}
			f()
		}
		time.Sleep(d)
	}
}

func runGitHub(c *cli.Context, handle func(error)) {
	for _, p := range config.GitHub.Projects {
		handle(githubCommitTweet(p))
	}
}

func runRetweet(c *cli.Context, handle func(error)) {
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
		cs := []TweetChecker{a.Filter.GetChecker(visionAPI)}
		if a.ScreenName != nil {
			ts, err := twitterAPI.RetweetAccount(*a.ScreenName, v, cs, a.Action)
			tweets = append(tweets, ts...)
			handle(err)
		} else {
			for _, name := range a.ScreenNames {
				ts, err := twitterAPI.RetweetAccount(name, v, cs, a.Action)
				tweets = append(tweets, ts...)
				handle(err)
			}
		}
	}
	for _, a := range config.Twitter.Searches {
		cs := []TweetChecker{a.Filter.GetChecker(visionAPI)}
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if a.ResultType != nil {
			v.Set("result_type", *a.ResultType)
		}
		if a.Query != nil {
			ts, err := twitterAPI.RetweetSearch(*a.Query, v, cs, a.Action)
			handle(err)
			tweets = append(tweets, ts...)
		} else {
			for _, query := range a.Queries {
				ts, err := twitterAPI.RetweetSearch(query, v, cs, a.Action)
				handle(err)
				tweets = append(tweets, ts...)
			}
		}
	}
	for _, t := range tweets {
		err := twitterAPI.NotifyToAll(&t, config.Twitter.Notification)
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
