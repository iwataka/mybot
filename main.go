package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	fateUrl              = "http://www.fate-sn.com/"
	tsuredurechildrenUrl = "http://tsuredurechildren.com/"
	period               = time.Duration(10) * time.Minute
)

var (
	cacheFile = os.ExpandEnv("$HOME/.cache/mybot/cache.json")
	ghClient  *github.Client
	twApi     *anaconda.TwitterApi
	cache     = &cacheData{
		make(map[string]map[string]string),
		make(map[string]int64),
	}
	projects = map[string]string{
		"vim":    "vim",
		"neovim": "neovim",
		"golang": "go",
	}
)

type cacheData struct {
	LatestCommitSHA map[string]map[string]string
	LatestTweetId   map[string]int64
}

func main() {
	err := os.MkdirAll(filepath.Dir(cacheFile), 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := cli.NewApp()
	app.Name = "mybot"
	app.Version = "0.0.1"

	before := func(c *cli.Context) error {
		err := beforeRunning()
		if err != nil {
			return err
		}
		return nil
	}

	runCmd := cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "send messages once",
		Before:  before,
		Action: func(c *cli.Context) error {
			err = run(c)
			if err != nil {
				return err
			}
			return nil
		},
	}

	serverCmd := cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "send messages periodically",
		Flags:   []cli.Flag{cli.StringFlag{Name: "log-file"}},
		Before:  before,
		Action: func(c *cli.Context) error {
			err = server(c)
			if err != nil {
				return err
			}
			return nil
		},
	}

	app.Commands = []cli.Command{runCmd, serverCmd}
	app.Run(os.Args)
}

func exitIfError(err error, code int) {
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(code)
	}
}

func beforeRunning() error {
	err := unmarshalCache()
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

func getenv(key string) (string, error) {
	result := os.Getenv(key)
	if result == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s=", key)
		text, _ := reader.ReadString('\n')
		err := os.Setenv(key, text)
		if err != nil {
			return "", err
		}
		result = strings.TrimSpace(text)
	}
	return result, nil
}

func run(c *cli.Context) error {
	runOnce(func(err error) { exitIfError(err, 1) })
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

	handleError(marshalCache())
}

func server(c *cli.Context) error {
	logFile := c.String("log-file")
	output := os.Stdout
	if logFile != "" {
		info, err := os.Stat(logFile)
		if err == nil || info.IsDir() {
			fmt.Printf("Invalid log file: %s", logFile)
			os.Exit(1)
		} else {
			output, err = os.Create(logFile)
			exitIfError(err, 1)
		}
	}
	logger := log.New(output, "", log.Ldate|log.Ltime|log.Lshortfile)
	for {
		runOnce(func(err error) { logger.Println(err) })
		time.Sleep(period)
	}
	return nil
}

func marshalCache() error {
	if cache != nil {
		data, err := json.MarshalIndent(cache, "", "\t")
		if err != nil {
			return err
		}
		ioutil.WriteFile(cacheFile, data, 0600)
	}
	return nil
}

func unmarshalCache() error {
	info, err := os.Stat(cacheFile)
	if err == nil && info != nil && !info.IsDir() {
		data, err := ioutil.ReadFile(cacheFile)
		if err != nil {
			return err
		}
		json.Unmarshal(data, cache)
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
		if err != nil && !ignoreTwitterError(err) {
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

func formatUrl(src, dest string) (string, error) {
	domainPat, err := regexp.Compile("[^:/]+://[^/]+")
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(dest, "/") {
		return domainPat.FindString(src) + dest, nil
	} else if strings.Index(dest, "://") == -1 {
		if strings.HasSuffix(src, "/") {
			return src + dest, nil
		} else {
			return src + "/" + dest, nil
		}
	} else {
		return dest, nil
	}
}

func ignoreTwitterError(err error) bool {
	if strings.Contains(err.Error(), "Status is a duplicate") {
		return true
	} else {
		return false
	}
}
