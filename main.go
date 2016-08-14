package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	"github.com/urfave/cli"
)

//go:generate go-bindata assets/... index.html 404.html config.template.yml

var (
	twitterAPI *TwitterAPI
	githubAPI  *GitHubAPI
	visionAPI  *VisionAPI
	config     *MybotConfig
	cache      *MybotCache
	logger     *Logger
)

var logFlag = cli.StringFlag{
	Name:  "log",
	Value: os.ExpandEnv("$HOME/.mybot-debug.log"),
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
		fmt.Println("No configuration file detected...")
		fmt.Printf("Create new sample configuration file: %s\n", c.String("config"))
		fmt.Println("Edit this file as you want")
		data, err := Asset("config.template.yml")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(c.String("config"), data, 0664)
		if err != nil {
			panic(err)
		}
		config, err = NewMybotConfig(c.String("config"))
	}

	githubAPI = NewGitHubAPI(nil, cache)
	twitterAPI = NewTwitterAPI(config.Authentication, cache)

	logger, err = NewLogger(c.String("log"), -1, twitterAPI, config)
	if err != nil {
		panic(err)
	}

	// visionAPI is nil if there exists no credential file
	visionAPI, err = NewVisionAPI(c.String("gcp-credential"))
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
			a, err := NewVisionAPI(c.String("gcp-credential"))
			logger.InfoIfError(err)
			if err == nil {
				*visionAPI = *a
			}

		})

	fmt.Printf("Open 127.0.0.1:%s to see the detail information\n", s.Port)
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
