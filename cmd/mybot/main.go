package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/src"
	"github.com/urfave/cli"
)

var (
	twitterAPI         *mybot.TwitterAPI
	visionAPI          *mybot.VisionAPI
	server             *mybot.MybotServer
	config             *mybot.MybotConfig
	cache              *mybot.MybotCache
	logger             *mybot.Logger
	status             *mybot.MybotStatus
	ctxt               *cli.Context
	userListenerChan   chan interface{}
	myselfListenerChan chan interface{}
)

func main() {
	var err error

	logFlag := cli.StringFlag{
		Name:  "log",
		Value: "mybot.log",
		Usage: "Log file's location",
	}

	configFlag := cli.StringFlag{
		Name:  "config",
		Value: "config.toml",
		Usage: "Config file's location",
	}

	cacheFlag := cli.StringFlag{
		Name:  "cache",
		Value: ".cache/mybot/cache.json",
		Usage: "Cache file's location",
	}

	credFlag := cli.StringFlag{
		Name:  "credential,c",
		Value: "",
		Usage: "Credential for Basic Authentication (ex: user:password)",
	}

	gcloudFlag := cli.StringFlag{
		Name:  "gcloud",
		Value: "google_application_credentials.json",
		Usage: "Credential file for Google Cloud Platform",
	}

	twitterFlag := cli.StringFlag{
		Name:  "twitter",
		Value: "twitter_authentication.json",
		Usage: "Credential file for Twitter API",
	}

	certFlag := cli.StringFlag{
		Name:  "cert",
		Value: "mybot.crt",
		Usage: "Certification file for HTTPS",
	}

	keyFlag := cli.StringFlag{
		Name:  "key",
		Value: "mybot.key",
		Usage: "Key file for HTTPS",
	}

	runFlags := []cli.Flag{
		logFlag,
		configFlag,
		cacheFlag,
		gcloudFlag,
		twitterFlag,
	}

	serveFlags := []cli.Flag{
		logFlag,
		configFlag,
		cacheFlag,
		credFlag,
		gcloudFlag,
		twitterFlag,
		certFlag,
		keyFlag,
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
	ctxt = c

	var err error
	cache, err = mybot.NewMybotCache(c.String("cache"))
	if err != nil {
		panic(err)
	}

	config, err = mybot.NewMybotConfig(c.String("config"), nil)
	if err != nil {
		panic(err)
	}

	visionAPI, err = mybot.NewVisionAPI(cache, config, c.String("gcloud"))
	if err != nil {
		panic(err)
	}

	config.SetVisionoAPI(visionAPI)

	twitterAuth := &mybot.TwitterAuth{}
	twitterAuth.FromJson(c.String("twitter"))
	mybot.SetConsumer(twitterAuth.ConsumerKey, twitterAuth.ConsumerSecret)
	twitterAPI = mybot.NewTwitterAPI(twitterAuth.AccessToken, twitterAuth.AccessTokenSecret, cache, config)

	logger, err = mybot.NewLogger(c.String("log"), -1, twitterAPI, config)
	if err != nil {
		panic(err)
	}

	status = &mybot.MybotStatus{}

	server = &mybot.MybotServer{
		Logger:     logger,
		TwitterAPI: twitterAPI,
		VisionAPI:  visionAPI,
		Cache:      cache,
		Config:     config,
		Status:     status,
	}

	return nil
}

func run(c *cli.Context) error {
	err := runTwitterWithoutStream()
	if err != nil {
		logger.Println(err)
	}
	err = cache.Save(ctxt.String("cache"))
	if err != nil {
		logger.Println(err)
	}
	return nil
}

func keepConnection(f func() error, intervalStr string, maxCount int) error {
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return err
	}
	t := time.Now()
	count := 0
	for {
		err := f()
		if time.Now().Sub(t) >= interval {
			count = 0
			t = time.Now()
		}
		if err != nil {
			logger.Println(err)
			switch err.(type) {
			case mybot.KillError:
				break
			default:
				count++
			}
		}
		if count >= maxCount {
			break
		}
	}
	msg := "Failed to keep connection"
	logger.Println(msg)
	return errors.New(msg)
}

func twitterListenMyself() {
	if status.TwitterListenMyselfStatus {
		return
	}
	status.TwitterListenMyselfStatus = true
	defer func() { status.TwitterListenMyselfStatus = false }()
	keepConnection(func() error {
		r := twitterAPI.DefaultDirectMessageReceiver
		listener, err := twitterAPI.ListenMyself(nil, r, ctxt.String("cache"))
		if err != nil {
			logger.Println(err)
			return err
		}
		myselfListenerChan = listener.C
		err = listener.Listen()
		if err != nil {
			logger.Println(err)
			return err
		}
		return nil
	}, "5m", 5)
}

func twitterListenUsers() {
	if status.TwitterListenUsersStatus {
		return
	}
	status.TwitterListenUsersStatus = true
	defer func() { status.TwitterListenUsersStatus = false }()
	keepConnection(func() error {
		listener, err := twitterAPI.ListenUsers(nil, ctxt.String("cache"))
		if err != nil {
			logger.Println(err)
			return err
		}
		userListenerChan = listener.C
		err = listener.Listen()
		if err != nil {
			logger.Println(err)
			return err
		}
		return nil
	}, "5m", 5)
}

func twitterPeriodically() {
	if status.TwitterStatus {
		return
	}
	status.TwitterStatus = true
	defer func() { status.TwitterStatus = false }()
	for {
		err := runTwitterWithStream()
		if err != nil {
			logger.Println(err)
			return
		}
		err = cache.Save(ctxt.String("cache"))
		if err != nil {
			logger.Println(err)
			return
		}
		d, err := time.ParseDuration(config.Twitter.Duration)
		if err != nil {
			logger.Println(err)
			return
		}
		time.Sleep(d)
	}
}

func monitorConfig() {
	if status.MonitorConfigStatus {
		return
	}
	status.MonitorConfigStatus = true
	defer func() { status.MonitorConfigStatus = false }()
	monitorFile(
		ctxt.String("config"),
		time.Duration(1)*time.Second,
		func() {
			cfg, err := mybot.NewMybotConfig(ctxt.String("config"), visionAPI)
			if err == nil {
				*config = *cfg
				reloadListeners()
			}
		},
	)
}

func monitorTwitterCred() {
	if status.MonitorTwitterCred {
		return
	}
	status.MonitorTwitterCred = true
	defer func() { status.MonitorTwitterCred = false }()
	monitorFile(
		ctxt.String("twitter"),
		time.Duration(1)*time.Second,
		func() {
			auth := &mybot.TwitterAuth{}
			err := auth.FromJson(ctxt.String("twitter"))
			if err != nil {
				mybot.SetConsumer(auth.ConsumerKey, auth.ConsumerSecret)
				api := mybot.NewTwitterAPI(auth.AccessToken, auth.AccessTokenSecret, cache, config)
				*twitterAPI = *api
				reloadListeners()
			}
		},
	)
}

func monitorGCloudCred() {
	if status.MonitorGCloudCred {
		return
	}
	status.MonitorGCloudCred = true
	defer func() { status.MonitorGCloudCred = false }()
	monitorFile(
		ctxt.String("gcloud"),
		time.Duration(1)*time.Second,
		func() {
			api, err := mybot.NewVisionAPI(cache, config, ctxt.String("gcloud"))
			if err != nil {
				*visionAPI = *api
				reloadListeners()
			}
		},
	)
}

func reloadListeners() {
	if userListenerChan != nil {
		userListenerChan <- os.Interrupt
	} else {
		twitterListenUsers()
	}
	if myselfListenerChan != nil {
		myselfListenerChan <- os.Interrupt
	} else {
		twitterListenMyself()
	}
}

func httpServer() {
	if status.HttpStatus {
		return
	}
	status.HttpStatus = true
	defer func() { status.HttpStatus = false }()
	cred := ctxt.String("credential")
	userAndPassword := strings.SplitN(cred, ":", 2)
	user := ""
	password := ""
	if len(userAndPassword) == 2 {
		user = userAndPassword[0]
		password = userAndPassword[1]
	}
	cert := ctxt.String("cert")
	key := ctxt.String("key")
	err := server.Init(user, password, cert, key)
	if err != nil {
		panic(err)
	}
}

func serve(c *cli.Context) error {
	go httpServer()

	go twitterListenMyself()
	go twitterListenUsers()
	go twitterPeriodically()

	go monitorConfig()
	go monitorTwitterCred()
	go monitorGCloudCred()

	ch := make(chan bool)
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

func runTwitterWithStream() error {
	tweets := []anaconda.Tweet{}
	for _, a := range config.Twitter.Searches {
		a.Filter.VisionAPI = visionAPI
		v := url.Values{}
		if a.Count > 0 {
			v.Set("count", fmt.Sprintf("%d", a.Count))
		}
		if len(a.ResultType) != 0 {
			v.Set("result_type", a.ResultType)
		}
		for _, query := range a.Queries {
			ts, err := twitterAPI.DoForSearch(query, v, a.Filter, a.Action)
			if err != nil {
				return err
			}
			tweets = append(tweets, ts...)
		}
	}
	for _, a := range config.Twitter.Favorites {
		v := url.Values{}
		if a.Count > 0 {
			v.Set("count", fmt.Sprintf("%d", a.Count))
		}
		a.Filter.VisionAPI = visionAPI
		for _, name := range a.ScreenNames {
			ts, err := twitterAPI.DoForFavorites(name, v, a.Filter, a.Action)
			tweets = append(tweets, ts...)
			if err != nil {
				return err
			}
		}
	}
	for _, t := range tweets {
		err := twitterAPI.NotifyToAll(&t)
		if err != nil {
			return err
		}
	}
	return nil
}

func runTwitterWithoutStream() error {
	err := runTwitterWithStream()
	if err != nil {
		return err
	}
	tweets := []anaconda.Tweet{}
	for _, a := range config.Twitter.Timelines {
		v := url.Values{}
		if a.Count > 0 {
			v.Set("count", fmt.Sprintf("%d", a.Count))
		}
		if a.ExcludeReplies != nil {
			v.Set("exclude_replies", fmt.Sprintf("%v", *a.ExcludeReplies))
		}
		if a.IncludeRts != nil {
			v.Set("include_rts", fmt.Sprintf("%v", *a.IncludeRts))
		}
		a.Filter.VisionAPI = visionAPI
		for _, name := range a.ScreenNames {
			ts, err := twitterAPI.DoForAccount(name, v, a.Filter, a.Action)
			tweets = append(tweets, ts...)
			if err != nil {
				return err
			}
		}
	}
	for _, t := range tweets {
		err := twitterAPI.NotifyToAll(&t)
		if err != nil {
			return err
		}
	}
	return nil
}