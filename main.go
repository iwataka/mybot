package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

var logger *log.Logger

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
	err := unmarshalCache(c.String("cache"))
	exitIfError(err, 1)
	err = unmarshalConfig(c.String("config"))
	exitIfError(err, 1)

	anaconda.SetConsumerKey(config.Tweet.ConsumerKey)
	anaconda.SetConsumerSecret(config.Tweet.ConsumerSecret)
	twitterApi = anaconda.NewTwitterApi(config.Tweet.AccessToken, config.Tweet.AccessTokenSecret)

	githubApi = github.NewClient(nil)

	return nil
}

func run(c *cli.Context) error {
	var err error
	logger, err = newLogger(c.String("log"))
	exitIfError(err, 1)

	runOnce(c, handleError)
	return nil
}

func serve(c *cli.Context) error {
	var err error
	logger, err = newLogger(c.String("log"))
	exitIfError(err, 1)

	go func() {
		for {
			handleError(twitterTalk())
			time.Sleep(time.Second * time.Duration(config.Talk.Interval))
		}
	}()

	go func() {
		for {
			runOnce(c, handleError)
			time.Sleep(time.Minute * time.Duration(config.Tweet.Interval))
		}
	}()

	initHttp()
	return nil
}

func runOnce(c *cli.Context, handle func(error)) {
	err := unmarshalConfig(c.String("config"))
	handle(err)
	for _, proj := range config.Tweet.Github {
		handle(githubCommitTweet(proj.User, proj.Repo))
	}
	for _, target := range config.Tweet.Retweet {
		handle(retweetTarget(target))
	}
	handle(marshalCache(c.String("cache")))
}

func handleError(err error) {
	if err != nil {
		e := twitterPost(err.Error())
		if e != nil {
			logger.Println(e)
		}
		logger.Println(err)
	}
}

func retweetTarget(target retweetConfig) error {
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

func newLogger(path string) (*log.Logger, error) {
	logFlag := log.Ldate | log.Ltime | log.Lshortfile
	if path == "" {
		path = ".mybot-debug.log"
	} else if info, err := os.Stat(path); os.IsExist(err) && info.IsDir() {
		path = filepath.Join(path, ".mybot-debug.log")
	}
	logFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return log.New(logFile, "", logFlag), nil
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
