package main

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

const (
	fateUrl              = "http://www.fate-sn.com/"
	tsuredurechildrenUrl = "http://tsuredurechildren.com/"
)

var (
	cacheFile io.ReadWriteCloser
	logFile   io.ReadWriteCloser
	ghClient  *github.Client
	twApi     *anaconda.TwitterApi
)

var cache = &cacheData{
	make(map[string]map[string]string),
	make(map[string]int64),
	make(map[string]int64),
}

var projects = map[string]string{
	"vim":    "vim",
	"neovim": "neovim",
	"golang": "go",
}

type cacheData struct {
	LatestCommitSHA map[string]map[string]string
	LatestTweetId   map[string]int64
	LatestDM        map[string]int64
}

func init() {
	var err error
	cachePath := os.ExpandEnv("$HOME/.cache/mybot/cache.json")
	cacheFile, err = os.Create(cachePath)
	err = os.MkdirAll(filepath.Dir(cachePath), 0600)
	exitIfError(err, 1)
}

func main() {
	app := cli.NewApp()
	app.Name = "mybot"
	app.Version = "0.0.1"

	runCmd := cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "send messages once",
		Flags:   []cli.Flag{cli.StringFlag{Name: "log-file"}},
		Before:  beforeRunning,
		Action:  run,
	}

	serverCmd := cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "send messages periodically",
		Flags:   []cli.Flag{cli.StringFlag{Name: "log-file"}},
		Before:  beforeRunning,
		Action:  server,
	}

	app.Commands = []cli.Command{runCmd, serverCmd}
	app.Run(os.Args)
}

func beforeRunning(c *cli.Context) error {
	dec := json.NewDecoder(cacheFile)
	err := dec.Decode(cache)
	exitIfError(err, 1)

	ghClient = github.NewClient(nil)

	consumerKey, err := getenv("MYBOT_TWITTER_CONSUMER_KEY")
	exitIfError(err, 1)
	consumerSecret, err := getenv("MYBOT_TWITTER_CONSUMER_SECRET")
	exitIfError(err, 1)
	accessToken, err := getenv("MYBOT_TWITTER_ACCESS_TOKEN")
	exitIfError(err, 1)
	accessTokenSecret, err := getenv("MYBOT_TWITTER_ACCESS_TOKEN_SECRET")
	exitIfError(err, 1)

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	twApi = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	return nil
}

func run(c *cli.Context) error {
	logger, err := newLogger(c)
	exitIfError(err, 1)

	runOnce(func(err error) {
		if err != nil {
			logger.Println(err)
		}
	})
	return nil
}

func runOnce(handleError func(error)) {
	var err error
	for user, repo := range projects {
		handleError(githubCommit(user, repo))
	}
	err = retweet("Fate_SN_Anime", false, func(t anaconda.Tweet) bool {
		text := strings.ToLower(t.Text)
		return strings.Contains(text, "heaven's feel") && strings.Contains(text, "劇場版")
	})
	handleError(err)
	err = retweet("sankakujougi", false, func(t anaconda.Tweet) bool {
		return strings.Contains(t.Text, "https://t.co/p3Zy7VoPcg")
	})
	handleError(err)

	enc := json.NewEncoder(cacheFile)
	handleError(enc.Encode(cache))
}

func newLogger(c *cli.Context) (*log.Logger, error) {
	logFlag := log.Ldate | log.Ltime | log.Lshortfile
	if c.IsSet("log-file") {
		logFile, err := os.Create(c.String(("log-file")))
		if err != nil {
			return nil, err
		}
		return log.New(logFile, "", logFlag), nil
	} else {
		return log.New(os.Stdout, "", logFlag), nil
	}
}

func server(c *cli.Context) error {
	logger, err := newLogger(c)
	exitIfError(err, 1)

	go func() {
		for {
			go func() {
				err := talk()
				if err != nil {
					logger.Println(err)
				}
			}()
			time.Sleep(time.Second * time.Duration(30))
		}
	}()

	for {
		runOnce(func(err error) {
			if err != nil {
				logger.Println(err)
			}
		})
		time.Sleep(time.Minute * time.Duration(10))
	}

	return nil
}

func githubCommit(user, repo string) error {
	commits, _, err := ghClient.Repositories.ListCommits(user, repo, nil)
	if err != nil {
		return err
	}
	latest := commits[0]
	userMap, userExists := cache.LatestCommitSHA[user]
	sha, repoExists := userMap[repo]
	if !userExists || !repoExists || sha != *latest.SHA {
		msg := user + "/" + repo + "\n" + *latest.HTMLURL
		_, err := twApi.PostTweet(msg, nil)
		if err != nil {
			return err
		}
		if !userExists {
			cache.LatestCommitSHA[user] = make(map[string]string)
		}
		cache.LatestCommitSHA[user][repo] = *latest.SHA
	}
	return nil
}

func retweet(screenName string, trimUser bool, checker func(anaconda.Tweet) bool) error {
	v := url.Values{}
	v.Set("screen_name", screenName)
	tweets, err := twApi.GetUserTimeline(v)
	if err != nil {
		return err
	}
	latestId, exists := cache.LatestTweetId[screenName]
	finds := false
	updates := false
	for i := len(tweets) - 1; i >= 0; i-- {
		tweet := tweets[i]
		if checker(tweet) {
			if exists && latestId == tweet.Id {
				finds = true
			} else {
				updates = true
				cache.LatestTweetId[screenName] = tweet.Id
				if finds {
					_, err := twApi.Retweet(tweet.Id, trimUser)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	if !exists && updates {
		_, err := twApi.Retweet(cache.LatestTweetId[screenName], trimUser)
		if err != nil {
			return err
		}
	}
	return nil
}

func talk() error {
	dms, err := twApi.GetDirectMessages(nil)
	if err != nil {
		return err
	}
	userToDM := make(map[string]anaconda.DirectMessage)
	for _, dm := range dms {
		sender := dm.SenderScreenName
		_, exists := userToDM[sender]
		if !exists {
			userToDM[sender] = dm
		}
	}
	for user, dm := range userToDM {
		latest, exists := cache.LatestDM[user]
		if !exists || latest != dm.Id {
			res, err := twApi.PostDMToScreenName(dm.Text, user)
			if err != nil {
				return err
			}
			cache.LatestDM[user] = res.Id
		}
	}
	return nil
}
