package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
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
		make(map[string]string),
	}
	projects = map[string]string{
		"vim":    "vim",
		"neovim": "neovim",
		"golang": "go",
	}
)

type cacheData struct {
	LatestCommitSHA map[string]map[string]string
	LatestID        map[string]string
}

func main() {
	os.Mkdir(filepath.Dir(cacheFile), 0600)

	app := cli.NewApp()
	app.Name = "mybot"
	app.Version = "0.0.1"

	runCmd := cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "send messages once",
		Action: func(c *cli.Context) error {
			err := setup()
			if err != nil {
				return err
			}
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
		Action: func(c *cli.Context) error {
			err := setup()
			if err != nil {
				return err
			}
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

func setup() error {
	err := unmarshalCache()
	if err != nil {
		return err
	}

	ghClient = github.NewClient(nil)

	consumerKey, err := getenv("MYBOT_TWITTER_CONSUMER_KEY")
	if err != nil {
		return err
	}
	consumerSecret, err := getenv("MYBOT_TWITTER_CONSUMER_SECRET")
	if err != nil {
		return err
	}
	accessToken, err := getenv("MYBOT_TWITTER_ACCESS_TOKEN")
	if err != nil {
		return err
	}
	accessTokenSecret, err := getenv("MYBOT_TWITTER_ACCESS_TOKEN_SECRET")
	if err != nil {
		return err
	}

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
	var err error
	for user, repo := range projects {
		err = latestCommit(user, repo)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}
	err = scrapeFateNews()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = scrapeTsuredurechildren()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	err = marshalCache()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func server(c *cli.Context) error {
	for {
		err := run(c)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
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

func latestCommit(user, repo string) error {
	commits, _, err := ghClient.Repositories.ListCommits(user, repo, nil)
	if err != nil {
		return err
	}
	latest := commits[0]
	if *latest.SHA != cache.LatestCommitSHA[user][repo] {
		msg := user + "/" + repo + "\n" + *latest.HTMLURL
		_, err := twApi.PostTweet(msg, nil)
		if err != nil && !ignoreTwitterError(err) {
			return err
		}
		if cache.LatestCommitSHA[user] == nil {
			cache.LatestCommitSHA[user] = make(map[string]string)
		}
		cache.LatestCommitSHA[user][repo] = *latest.SHA
	}
	return nil
}

func scrapeFateNews() error {
	doc, err := goquery.NewDocument(fateUrl)
	if err != nil {
		return err
	}

	latestDate := ""
	err = nil
	result := make(map[string]string)
	doc.Find(".news-list li").Each(func(i int, s *goquery.Selection) {
		dateBlock := s.Find(".day")
		a := s.Find("a")
		if dateBlock != nil && a != nil {
			date := dateBlock.Text()
			if latestDate == "" || latestDate == date {
				latestDate = date
				url, exists := a.Attr("href")
				if !exists {
					err = errors.New(fmt.Sprintf("%s is not found in %s", "href", a.Html))
					return
				}
				url, err = formatUrl(fateUrl, url)
				if err != nil {
					return
				}
				result[date+" "+a.Text()] = url
			}
		}
	})
	latestID, exists := cache.LatestID[fateUrl]
	if !exists || latestID != latestDate {
		cache.LatestID[fateUrl] = latestDate
		for title, url := range result {
			_, err := twApi.PostTweet(title+"\n"+url, nil)
			if err != nil && !ignoreTwitterError(err) {
				return err
			}
		}
	}
	return err
}

func scrapeTsuredurechildren() error {
	doc, err := goquery.NewDocument(tsuredurechildrenUrl)
	if err != nil {
		return err
	}
	article := doc.Find("div .column-one article").First()
	titleBlock := article.Find(".post-title a")
	title := titleBlock.Text()
	contentUrl, exists := titleBlock.Attr("href")
	if !exists {
		return errors.New(fmt.Sprintf("%s is not found in %s", "href", titleBlock.Html))
	}
	img := article.Find("div .blog-entry div img")
	imgUrl, exists := img.Attr("src")
	if !exists {
		return errors.New(fmt.Sprintf("%s is not found in %s", "src", img.Html))
	}
	latestID, exists := cache.LatestID[tsuredurechildrenUrl]
	if !exists || latestID != title {
		cache.LatestID[tsuredurechildrenUrl] = title
		base64Img, err := base64ImageUrl(imgUrl)
		if err != nil {
			return err
		}
		media, err := twApi.UploadMedia(base64Img)
		if err != nil {
			return err
		}
		v := url.Values{}
		v.Add("media_ids", media.MediaIDString)
		_, err = twApi.PostTweet(title+"\n"+contentUrl, v)
		if err != nil && !ignoreTwitterError(err) {
			return err
		}
	}

	return nil
}

func base64ImageUrl(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(body), nil
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
