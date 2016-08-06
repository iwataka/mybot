package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

var (
	cachePath = os.ExpandEnv("$HOME/.cache/mybot/cache.json")
	cache     *cacheData
)

var githubProjects = map[string]string{
	"vim":    "vim",
	"neovim": "neovim",
	"golang": "go",
}

type cacheData struct {
	LatestCommitSHA map[string]map[string]string
	LatestTweetId   map[string]int64
	LatestDM        map[string]int64
}

func unmarshalCache(path string) error {
	if cache == nil {
		cache = &cacheData{
			make(map[string]map[string]string),
			make(map[string]int64),
			make(map[string]int64),
		}
	}

	info, _ := os.Stat(path)
	if info != nil && !info.IsDir() {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, cache)
		if err != nil {
			return err
		}
	}
	return nil
}

func marshalCache(path string) error {
	var err error
	err = os.MkdirAll(filepath.Dir(path), 0600)
	if err != nil {
		return err
	}
	if cache != nil {
		data, err := json.Marshal(cache)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, data, 0600)
		if err != nil {
			return err
		}
	}
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
		Flags:   []cli.Flag{cli.StringFlag{Name: "log-file", Value: ""}},
		Before:  beforeRunning,
		Action:  run,
	}

	serveCmd := cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "send messages periodically",
		Flags:   []cli.Flag{cli.StringFlag{Name: "log-file", Value: ""}},
		Before:  beforeRunning,
		Action:  serve,
	}

	app.Commands = []cli.Command{runCmd, serveCmd}
	app.Run(os.Args)
}

func beforeRunning(c *cli.Context) error {
	err := unmarshalCache(cachePath)
	exitIfError(err, 1)

	twitterConsumerKey, err := getenv("MYBOT_TWITTER_CONSUMER_KEY")
	exitIfError(err, 1)
	twitterConsumerSecret, err := getenv("MYBOT_TWITTER_CONSUMER_SECRET")
	exitIfError(err, 1)
	twitterAccessToken, err := getenv("MYBOT_TWITTER_ACCESS_TOKEN")
	exitIfError(err, 1)
	twitterAccessTokenSecret, err := getenv("MYBOT_TWITTER_ACCESS_TOKEN_SECRET")
	exitIfError(err, 1)
	anaconda.SetConsumerKey(twitterConsumerKey)
	anaconda.SetConsumerSecret(twitterConsumerSecret)
	twitterApi = anaconda.NewTwitterApi(twitterAccessToken, twitterAccessTokenSecret)

	githubApi = github.NewClient(nil)

	return nil
}

func run(c *cli.Context) error {
	logger, err := newLogger(c.String("log-file"))
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
	for user, repo := range githubProjects {
		handleError(githubCommitTweet(user, repo))
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

	handleError(marshalCache(cachePath))
}

func newLogger(path string) (*log.Logger, error) {
	logFlag := log.Ldate | log.Ltime | log.Lshortfile
	if path == "" {
		p, err := filepath.Abs(filepath.Dir(os.Args[0]))
		p = filepath.Join(path, ".mybot-debug.log")
		if err != nil {
			return nil, err
		}
		path = p
	}
	logFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return log.New(logFile, "", logFlag), nil
}

func serve(c *cli.Context) error {
	logger, err := newLogger(c.String("log-file"))
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

	go func() {
		for {
			runOnce(func(err error) {
				if err != nil {
					logger.Println(err)
				}
			})
			time.Sleep(time.Minute * time.Duration(10))
		}
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to mybot root!")
	fmt.Fprintln(w, "This is being under development.")
	fmt.Fprintln(w, "This domain will be used to provide API.")
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
