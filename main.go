package main

import (
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
)

//go:generate go-bindata assets/... index.html 404.html

var (
	twitterAPI *TwitterAPI
	githubAPI  *GitHubAPI
	visionAPI  *VisionAPI
	config     *MybotConfig
	cache      *MybotCache
	logger     *MultiLogger
)

var logFlag = cli.StringFlag{
	Name:  "log",
	Value: ".mybot-debug.log",
	Usage: "Log file's location",
}

var configFlag = cli.StringFlag{
	Name:  "config",
	Value: "config.yml",
	Usage: "Config file's location",
}

var cacheFlag = cli.StringFlag{
	Name:  "cache",
	Value: os.ExpandEnv("$HOME/.cache/mybot/cache.json"),
	Usage: "Cache file's location",
}

var visionCredentialFlag = cli.StringFlag{
	Name:  "gcp-credential",
	Value: "credential.json",
	Usage: "Location of Google Cloud Platform credential file",
}

var flags = []cli.Flag{
	logFlag,
	configFlag,
	cacheFlag,
	visionCredentialFlag,
}

func main() {
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
	err := app.Run(os.Args)
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

	twitterLogger := NewTwitterLogger(twitterAPI, config.Log)
	logger, err = NewLogger(c.String("log"), "", -1, []Logger{twitterLogger})
	if err != nil {
		panic(err)
	}

	// visionAPI is nil if there exists no credential file
	visionAPI, err = NewVisionAPI(c.String("gcp-credential"))
	logger.InfoIfError(err)

	return nil
}

func run(c *cli.Context) error {
	runGitHub(c, logger.InfoIfError)
	runRetweet(c, logger.InfoIfError)
	logger.InfoIfError(cache.Save(c.String("cache")))
	return nil
}

func serve(c *cli.Context) error {
	s := config.Option
	s.Logger = logger
	s.TwitterAPI = twitterAPI
	s.VisionAPI = visionAPI

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
			d, err := time.ParseDuration(config.GitHub.Duration)
			logger.FatalIfError(err)
			time.Sleep(d)
		}
	}()

	go func() {
		for {
			runRetweet(c, logger.InfoIfError)
			d, err := time.ParseDuration(config.Retweet.Duration)
			logger.FatalIfError(err)
			time.Sleep(d)
		}
	}()

	go func() {
		w, err := fsnotify.NewWatcher()
		logger.InfoIfError(err)
		defer w.Close()
		for {
			select {
			case event := <-w.Events:
				if event.Op&fsnotify.Write == fsnotify.Write &&
					event.Op&fsnotify.Create == fsnotify.Create {
					err := config.ReadFile(c.String("config"))
					logger.InfoIfError(err)
				}
			case err := <-w.Errors:
				logger.InfoIfError(err)
			}
		}
	}()

	go func() {
		w, err := fsnotify.NewWatcher()
		logger.InfoIfError(err)
		defer w.Close()
		for {
			select {
			case event := <-w.Events:
				if event.Op&fsnotify.Write == fsnotify.Write &&
					event.Op&fsnotify.Create == fsnotify.Create {
					a, err := NewVisionAPI(c.String("gcp-credential"))
					logger.InfoIfError(err)
					if err == nil {
						visionAPI = a
						s.VisionAPI = a
					}
				}
			case err := <-w.Errors:
				logger.InfoIfError(err)
			}
		}
	}()

	err := s.Init()
	if err != nil {
		panic(err)
	}

	return nil
}

func runGitHub(c *cli.Context, handle func(error)) {
	for _, p := range config.GitHub.Projects {
		handle(githubCommitTweet(p))
	}
}

func runRetweet(c *cli.Context, handle func(error)) {
	for _, a := range config.Retweet.Accounts {
		ts, err := twitterAPI.RetweetWithChecker(a.Name, false, a.GetChecker(visionAPI))
		handle(err)
		if ts != nil {
			for _, t := range ts {
				err := twitterAPI.NotifyToAll(t.RetweetedStatus, config.Retweet.Notification)
				handle(err)
			}
		}
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
		_, userExists := cache.LatestCommitSHA[p.User]
		if !userExists {
			cache.LatestCommitSHA[p.User] = make(map[string]string)
		}
		cache.LatestCommitSHA[p.User][p.Repo] = *commit.SHA
	}
	return nil
}
