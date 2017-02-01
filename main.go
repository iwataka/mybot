package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/lib"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var (
	twitterAPI  *mybot.TwitterAPI
	twitterApp  *mybot.OAuthApp
	twitterAuth *mybot.OAuthCredentials
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	config      *mybot.Config
	cache       mybot.Cache
	logger      *mybot.Logger
	status      *mybot.Status

	ctxt *cli.Context

	userListenerStream *anaconda.Stream
	dmListenerStream   *anaconda.Stream
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	configDir := filepath.Join(home, ".config", "mybot")
	cacheDir := filepath.Join(home, ".cache", "mybot")

	logFlag := cli.StringFlag{
		Name:  "log",
		Value: filepath.Join(cacheDir, "mybot.log"),
		Usage: "Log file's location",
	}

	configFlag := cli.StringFlag{
		Name:  "config",
		Value: filepath.Join(configDir, "config.toml"),
		Usage: "Config file's location",
	}

	cacheFlag := cli.StringFlag{
		Name:  "cache",
		Value: filepath.Join(cacheDir, "cache.json"),
		Usage: "Cache file's location",
	}

	gcloudFlag := cli.StringFlag{
		Name:  "gcloud",
		Value: filepath.Join(configDir, "google_application_credentials.json"),
		Usage: "Credential file for Google Cloud Platform",
	}

	twitterAppFlag := cli.StringFlag{
		Name:  "twitter-app",
		Value: filepath.Join(configDir, "twitter_application_settings.toml"),
		Usage: "Application Setting file for Twitter API",
	}

	twitterFlag := cli.StringFlag{
		Name:  "twitter",
		Value: filepath.Join(configDir, "twitter_authentication.toml"),
		Usage: "Credential file for Twitter API",
	}

	certFlag := cli.StringFlag{
		Name:  "cert",
		Value: filepath.Join(configDir, "mybot.crt"),
		Usage: "Certification file for server",
	}

	keyFlag := cli.StringFlag{
		Name:  "key",
		Value: filepath.Join(configDir, "mybot.key"),
		Usage: "Key file for server",
	}

	hostFlag := cli.StringFlag{
		Name:  "host,H",
		Value: "",
		Usage: "Host this server listen on",
	}

	portFlag := cli.StringFlag{
		Name:  "port,P",
		Value: "",
		Usage: "Port this server listen on",
	}

	runFlags := []cli.Flag{
		logFlag,
		configFlag,
		cacheFlag,
		gcloudFlag,
		twitterAppFlag,
		twitterFlag,
	}

	serveFlags := []cli.Flag{
		certFlag,
		keyFlag,
		hostFlag,
		portFlag,
	}
	// All `run` flags should be `serve` flag
	for _, f := range runFlags {
		serveFlags = append(serveFlags, f)
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

	apiFlag := cli.BoolFlag{
		Name:  "api",
		Usage: "Use API to validate configuration",
	}

	validateFlags := []cli.Flag{
		configFlag,
		cacheFlag,
		twitterFlag,
		apiFlag,
	}

	validateCmd := cli.Command{
		Name:    "validate",
		Aliases: []string{"v"},
		Usage:   "Validates the user configuration",
		Flags:   validateFlags,
		Before:  beforeValidate,
		Action:  validate,
	}

	app.Commands = []cli.Command{runCmd, serveCmd, validateCmd}
	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func beforeRunning(c *cli.Context) error {
	ctxt = c

	err := beforeValidate(c)
	if err != nil {
		panic(err)
	}

	if info, err := os.Stat(c.String("gcloud")); err == nil && !info.IsDir() {
		visionAPI, err = mybot.NewVisionAPI(c.String("gcloud"))
		if err != nil {
			panic(err)
		}
		languageAPI, err = mybot.NewLanguageAPI(c.String("gcloud"))
		if err != nil {
			panic(err)
		}
	} else {
		visionAPI = &mybot.VisionAPI{}
		visionAPI.File = c.String("gcloud")
		languageAPI = &mybot.LanguageAPI{}
		languageAPI.File = c.String("gcloud")
	}

	logger, err = mybot.NewLogger(c.String("log"), -1, twitterAPI, config)
	if err != nil {
		panic(err)
	}

	status = mybot.NewStatus()
	status.UpdateTwitterAuth(twitterAPI)

	return nil
}

func beforeValidate(c *cli.Context) error {
	var err error
	cache, err = mybot.NewFileCache(c.String("cache"))
	if err != nil {
		panic(err)
	}

	config, err = mybot.NewConfig(c.String("config"))
	if err != nil {
		panic(err)
	}

	twitterApp = &mybot.OAuthApp{}
	err = twitterApp.Decode(c.String("twitter-app"))
	if err != nil {
		panic(err)
	}

	twitterAuth = &mybot.OAuthCredentials{}
	err = twitterAuth.Decode(c.String("twitter"))
	if err != nil {
		panic(err)
	}

	anaconda.SetConsumerKey(twitterApp.ConsumerKey)
	anaconda.SetConsumerSecret(twitterApp.ConsumerSecret)
	twitterAPI = mybot.NewTwitterAPI(twitterAuth, cache, config)

	return nil
}

func run(c *cli.Context) error {
	err := runTwitterWithoutStream()
	if err != nil {
		logger.Println(err)
	}
	err = cache.Save()
	if err != nil {
		logger.Println(err)
	}
	return nil
}

func twitterListenDM() {
	if !status.PassTwitterAuth {
		return
	}

	status.LockListenDMRoutine()
	defer status.UnlockListenDMRoutine()

	r := twitterAPI.DefaultDirectMessageReceiver
	listener, err := twitterAPI.ListenMyself(nil, r)
	if err != nil {
		logger.Println(err)
		return
	}
	dmListenerStream = listener.Stream
	err = listener.Listen()
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Println("Failed to keep twitter direct message listener")
}

func twitterListenUsers() {
	if !status.PassTwitterAuth {
		return
	}

	status.LockListenUsersRoutine()
	defer status.UnlockListenUsersRoutine()

	listener, err := twitterAPI.ListenUsers(nil)
	if err != nil {
		logger.Println(err)
		return
	}
	userListenerStream = listener.Stream
	err = listener.Listen(visionAPI, languageAPI, cache)
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Println("Failed to keep twitter user listener")
}

func twitterPeriodically() {
	if !status.PassTwitterAuth {
		return
	}

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
		err = cache.Save()
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
	file := ctxt.String("config")

	if status.GetMonitorStatus(file) {
		return
	}
	status.SetMonitorStatus(file, true)
	defer func() { status.SetMonitorStatus(file, false) }()

	monitorFile(
		file,
		time.Duration(1)*time.Second,
		func() {
			status.LockMonitor(file)
			defer status.UnlockMonitor(file)
			cfg, err := mybot.NewConfig(file)
			if err == nil {
				*config = *cfg
				reloadListeners()
			}
			status.SendToMonitor(file, true)
		},
	)
}

func monitorTwitterApp() {
	file := ctxt.String("twitter-app")

	if status.GetMonitorStatus(file) {
		return
	}
	status.SetMonitorStatus(file, true)
	defer func() { status.SetMonitorStatus(file, false) }()

	monitorFile(
		file,
		time.Duration(1)*time.Second,
		func() {
			status.LockMonitor(file)
			defer status.UnlockMonitor(file)

			app := &mybot.OAuthApp{}
			err := app.Decode(file)
			if err == nil {
				*twitterApp = *app
				anaconda.SetConsumerKey(app.ConsumerKey)
				anaconda.SetConsumerSecret(app.ConsumerSecret)
				status.UpdateTwitterAuth(twitterAPI)
				reloadListeners()
			}

			status.SendToMonitor(file, true)
		},
	)
}

func monitorTwitterCred() {
	file := ctxt.String("twitter")

	if status.GetMonitorStatus(file) {
		return
	}
	status.SetMonitorStatus(file, true)
	defer func() { status.SetMonitorStatus(file, false) }()

	monitorFile(
		file,
		time.Duration(1)*time.Second,
		func() {
			status.LockMonitor(file)
			defer status.UnlockMonitor(file)

			auth := &mybot.OAuthCredentials{}
			err := auth.Decode(file)
			if err == nil {
				*twitterAuth = *auth
				*twitterAPI = *mybot.NewTwitterAPI(auth, cache, config)
				status.UpdateTwitterAuth(twitterAPI)
				reloadListeners()
			}

			status.SendToMonitor(file, true)
		},
	)
}

func monitorGCloudCred() {
	file := ctxt.String("gcloud")

	if status.GetMonitorStatus(file) {
		return
	}
	status.SetMonitorStatus(file, true)
	defer func() { status.SetMonitorStatus(file, false) }()

	monitorFile(
		file,
		time.Duration(1)*time.Second,
		func() {
			status.LockMonitor(file)
			defer status.UnlockMonitor(file)
			vis, err := mybot.NewVisionAPI(file)
			if err == nil {
				*visionAPI = *vis
				return
			}
			lang, err := mybot.NewLanguageAPI(file)
			if err == nil {
				*languageAPI = *lang
				return
			}
			reloadListeners()
			status.SendToMonitor(file, true)
		},
	)
}

func reloadListeners() {
	if userListenerStream != nil {
		userListenerStream.Stop()
	}
	go twitterListenUsers()

	if dmListenerStream != nil {
		dmListenerStream.Stop()
	}
	go twitterListenDM()

	if !status.TwitterStatus {
		go twitterPeriodically()
	}
}

func httpServer() {
	if status.ServerStatus {
		return
	}
	status.ServerStatus = true
	defer func() { status.ServerStatus = false }()
	host := ctxt.String("host")
	port := ctxt.String("port")
	cert := ctxt.String("cert")
	key := ctxt.String("key")
	err := startServer(host, port, cert, key)
	if err != nil {
		panic(err)
	}
}

func serve(c *cli.Context) error {
	go httpServer()

	go twitterListenDM()
	go twitterListenUsers()
	go twitterPeriodically()

	go monitorConfig()
	go monitorTwitterCred()
	go monitorTwitterApp()
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
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if len(a.ResultType) != 0 {
			v.Set("result_type", a.ResultType)
		}
		for _, query := range a.Queries {
			ts, err := twitterAPI.ProcessSearch(query, v, a.Filter, visionAPI, languageAPI, a.Action)
			if err != nil {
				return err
			}
			tweets = append(tweets, ts...)
		}
	}
	for _, a := range config.Twitter.Favorites {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		for _, name := range a.ScreenNames {
			ts, err := twitterAPI.ProcessFavorites(name, v, a.Filter, visionAPI, languageAPI, a.Action)
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

	for _, a := range config.Twitter.APIs {
		msg, err := a.Message()
		if err != nil {
			return err
		}
		_, err = twitterAPI.PostTweet(msg, nil)
		if mybot.CheckTwitterError(err) {
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
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if a.ExcludeReplies != nil {
			v.Set("exclude_replies", fmt.Sprintf("%v", *a.ExcludeReplies))
		}
		if a.IncludeRts != nil {
			v.Set("include_rts", fmt.Sprintf("%v", *a.IncludeRts))
		}
		for _, name := range a.ScreenNames {
			ts, err := twitterAPI.ProcessTimeline(name, v, a.Filter, visionAPI, languageAPI, a.Action)
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

func validate(c *cli.Context) error {
	err := config.Validate()
	if err != nil {
		panic(err)
	}
	if c.Bool("api") {
		err := config.ValidateWithAPI(twitterAPI)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
