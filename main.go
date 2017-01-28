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
	twitterAuth *mybot.TwitterAuth
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	config      *mybot.Config
	cache       *mybot.Cache
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

	confPrefix := ".config/mybot/"
	cachePrefix := ".cache/mybot/"

	logFlag := cli.StringFlag{
		Name:  "log",
		Value: filepath.Join(home, cachePrefix+"mybot.log"),
		Usage: "Log file's location",
	}

	configFlag := cli.StringFlag{
		Name:  "config",
		Value: filepath.Join(home, confPrefix+"config.toml"),
		Usage: "Config file's location",
	}

	cacheFlag := cli.StringFlag{
		Name:  "cache",
		Value: filepath.Join(home, cachePrefix+"cache.json"),
		Usage: "Cache file's location",
	}

	gcloudFlag := cli.StringFlag{
		Name:  "gcloud",
		Value: filepath.Join(home, confPrefix+"google_application_credentials.json"),
		Usage: "Credential file for Google Cloud Platform",
	}

	twitterFlag := cli.StringFlag{
		Name:  "twitter",
		Value: filepath.Join(home, confPrefix+"twitter_authentication.toml"),
		Usage: "Credential file for Twitter API",
	}

	certFlag := cli.StringFlag{
		Name:  "cert",
		Value: filepath.Join(home, confPrefix+"mybot.crt"),
		Usage: "Certification file for server",
	}

	keyFlag := cli.StringFlag{
		Name:  "key",
		Value: filepath.Join(home, confPrefix+"mybot.key"),
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
		twitterFlag,
	}

	serveFlags := []cli.Flag{
		logFlag,
		configFlag,
		cacheFlag,
		gcloudFlag,
		twitterFlag,
		certFlag,
		keyFlag,
		hostFlag,
		portFlag,
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
	cache, err = mybot.NewCache(c.String("cache"))
	if err != nil {
		panic(err)
	}

	config, err = mybot.NewConfig(c.String("config"))
	if err != nil {
		panic(err)
	}

	twitterAuth = &mybot.TwitterAuth{}
	err = twitterAuth.Read(c.String("twitter"))
	if err != nil {
		panic(err)
	}
	anaconda.SetConsumerKey(twitterAuth.ConsumerKey)
	anaconda.SetConsumerSecret(twitterAuth.ConsumerSecret)
	twitterAPI = mybot.NewTwitterAPI(twitterAuth, cache, config)

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

func twitterListenDM() {
	if !status.PassTwitterAuth {
		return
	}

	status.LockListenDMRoutine()
	defer status.UnlockListenDMRoutine()

	r := twitterAPI.DefaultDirectMessageReceiver
	listener, err := twitterAPI.ListenMyself(nil, r, ctxt.String("cache"))
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

	listener, err := twitterAPI.ListenUsers(nil, ctxt.String("cache"))
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
			status.MonitorConfigStatusMutex.Lock()
			defer status.MonitorConfigStatusMutex.Unlock()
			cfg, err := mybot.NewConfig(ctxt.String("config"))
			if err == nil {
				*config = *cfg
				reloadListeners()
			}
			status.SendToMonitorConfigStatusChans(true)
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
			status.MonitorTwitterCredMutex.Lock()
			defer status.MonitorTwitterCredMutex.Unlock()
			auth := &mybot.TwitterAuth{}
			err := auth.Read(ctxt.String("twitter"))
			if err == nil {
				anaconda.SetConsumerKey(auth.ConsumerKey)
				anaconda.SetConsumerSecret(auth.ConsumerSecret)
				api := mybot.NewTwitterAPI(auth, cache, config)
				*twitterAPI = *api
				status.UpdateTwitterAuth(api)
				reloadListeners()
			}
			status.SendToMonitorTwitterCredChans(true)
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
			status.MonitorGCloudCredMutex.Lock()
			defer status.MonitorGCloudCredMutex.Unlock()
			vis, err := mybot.NewVisionAPI(ctxt.String("gcloud"))
			if err == nil {
				*visionAPI = *vis
				return
			}
			lang, err := mybot.NewLanguageAPI(ctxt.String("gcloud"))
			if err == nil {
				*languageAPI = *lang
				return
			}
			reloadListeners()
			status.SendToMonitorTwitterCredChans(true)
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
			ts, err := twitterAPI.DoForSearch(query, v, a.Filter, visionAPI, languageAPI, a.Action)
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
			ts, err := twitterAPI.DoForFavorites(name, v, a.Filter, visionAPI, languageAPI, a.Action)
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
			ts, err := twitterAPI.DoForAccount(name, v, a.Filter, visionAPI, languageAPI, a.Action)
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
