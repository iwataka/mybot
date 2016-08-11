package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

var (
	logger  *log.Logger
	logFile string
)

var logFlag = cli.StringFlag{
	Name:  "log",
	Value: "",
}

var configFlag = cli.StringFlag{
	Name:  "config",
	Value: "",
}

var cacheFlag = cli.StringFlag{
	Name:  "cache",
	Value: "",
}

func initLogger(path string) error {
	logFlag := log.Ldate | log.Ltime | log.Lshortfile
	logFile = path
	if logFile == "" {
		logFile = ".mybot-debug.log"
	} else if info, err := os.Stat(logFile); os.IsExist(err) && info.IsDir() {
		logFile = filepath.Join(logFile, ".mybot-debug.log")
	}
	file, err := os.Create(logFile)
	if err != nil {
		return err
	}
	logger = log.New(file, "", logFlag)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "mybot"
	app.Version = "0.0.1"

	runCmd := cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "send messages once",
		Flags:   []cli.Flag{logFlag, configFlag, cacheFlag},
		Before:  beforeRunning,
		Action:  run,
	}

	serveCmd := cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "send messages periodically",
		Flags:   []cli.Flag{logFlag, configFlag, cacheFlag},
		Before:  beforeRunning,
		Action:  serve,
	}

	app.Commands = []cli.Command{runCmd, serveCmd}
	app.Run(os.Args)
}

func beforeRunning(c *cli.Context) error {
	err := initLogger(c.String("log"))
	exitIfError(err, 1)

	err = unmarshalCache(c.String("cache"))
	exitIfError(err, 1)
	err = unmarshalConfig(c.String("config"))
	exitIfError(err, 1)

	anaconda.SetConsumerKey(config.Authentication.ConsumerKey)
	anaconda.SetConsumerSecret(config.Authentication.ConsumerSecret)
	twitterApi = anaconda.NewTwitterApi(config.Authentication.AccessToken, config.Authentication.AccessTokenSecret)

	githubApi = github.NewClient(nil)

	return nil
}

func run(c *cli.Context) error {
	runOnce(c, logError)
	return nil
}

func serve(c *cli.Context) error {
	ghMutex := new(sync.Mutex)
	rtMutex := new(sync.Mutex)

	go func() {
		for {
			logError(twitterInteract())
			d, err := time.ParseDuration(config.Interaction.Duration)
			logFatal(err)
			time.Sleep(d)
		}
	}()

	go func() {
		for {
			ghMutex.Lock()
			runGitHub(c, logError)
			ghMutex.Unlock()

			d, err := time.ParseDuration(config.GitHub.Duration)
			logFatal(err)
			time.Sleep(d)
		}
	}()

	go func() {
		for {
			rtMutex.Lock()
			runRetweet(c, logError)
			rtMutex.Unlock()

			d, err := time.ParseDuration(config.Retweet.Duration)
			logFatal(err)
			time.Sleep(d)
		}
	}()

	go func() {
		for {
			ghMutex.Lock()
			rtMutex.Lock()
			unmarshalConfig(c.String("config"))
			ghMutex.Unlock()
			rtMutex.Unlock()

			d, err := getMinDuration()
			logFatal(err)
			time.Sleep(d)
		}
	}()

	initHttp()
	return nil
}

func getMinDuration() (time.Duration, error) {
	gd, err := time.ParseDuration(config.GitHub.Duration)
	if err != nil {
		return 0, err
	}
	rd, err := time.ParseDuration(config.Retweet.Duration)
	if err != nil {
		return 0, err
	}
	if gd < rd {
		return time.Duration(gd), nil
	} else {
		return time.Duration(rd), nil
	}
}

func runGitHub(c *cli.Context, handle func(error)) {
	for _, proj := range config.GitHub.Projects {
		handle(githubCommitTweet(proj.User, proj.Repo))
	}
}

func runRetweet(c *cli.Context, handle func(error)) {
	for _, account := range config.Retweet.Accounts {
		handle(retweetTarget(account))
	}
}

func runOnce(c *cli.Context, handle func(error)) {
	err := unmarshalConfig(c.String("config"))
	handle(err)
	runGitHub(c, handle)
	runRetweet(c, handle)
	handle(marshalCache(c.String("cache")))
}

func logError(err error) {
	if err != nil {
		l := config.Log
		if l != nil {
			e := twitterPost(err.Error(), l.AllowSelf, l.Users)
			if e != nil {
				logger.Println(e)
			}
		}
		logger.Println(err)
	}
}

func logFatal(err error) {
	logError(err)
	if err != nil {
		panic(err)
	}
}

func retweetTarget(target accountConfig) error {
	regexps := make([]*regexp.Regexp, len(target.Patterns), len(target.Patterns))
	for i, pat := range target.Patterns {
		r, err := regexp.Compile(pat)
		if err != nil {
			return err
		}
		regexps[i] = r
	}
	return twitterRetweet(target.Name, false, func(t anaconda.Tweet) bool {
		for _, r := range regexps {
			if !r.MatchString(t.Text) {
				return false
			}
		}
		for key, val := range target.Opts {
			if key == "hasMedia" {
				if val != (len(t.Entities.Media) != 0) {
					return false
				}
			} else if key == "hasUrl" {
				if val != (len(t.Entities.Urls) != 0) {
					return false
				}
			} else if key == "retweeted" {
				if val != t.Retweeted {
					return false
				}
			}
		}
		return true
	})
}

func githubCommitTweet(user, repo string) error {
	commit, err := githubCommit(user, repo)
	if err != nil {
		return err
	}
	if commit != nil {
		msg := user + "/" + repo + "\n" + *commit.HTMLURL
		_, err := twitterApi.PostTweet(msg, nil)
		if err != nil {
			return err
		}
		updateCommitSHA(user, repo, commit)
	}
	return nil
}

func updateCommitSHA(user, repo string, commit *github.RepositoryCommit) {
	_, userExists := cache.LatestCommitSHA[user]
	if !userExists {
		cache.LatestCommitSHA[user] = make(map[string]string)
	}
	cache.LatestCommitSHA[user][repo] = *commit.SHA
}
