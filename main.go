package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/lib"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"gopkg.in/mgo.v2"
	"log"
)

//go:generate go-bindata assets/...
//go:generate mockgen -source=models/slack.go -destination=mocks/slack.go -package=mocks
//go:generate mockgen -source=models/twitter.go -destination=mocks/twitter.go -package=mocks
//go:generate mockgen -source=lib/vision.go -destination=mocks/vision.go -package=mocks
//go:generate mockgen -source=lib/language.go -destination=mocks/language.go -package=mocks

var (
	twitterAPI         *mybot.TwitterAPI
	userListenerStream *anaconda.Stream
	dmListenerStream   *anaconda.Stream
	visionAPI          *mybot.VisionAPI
	languageAPI        *mybot.LanguageAPI
	slackAPI           *mybot.SlackAPI
	slackListener      *mybot.SlackListener
	twitterApp         mybot.OAuthApp
	twitterAuth        mybot.OAuthCreds
	config             mybot.Config
	cache              mybot.Cache

	executers map[int]*sync.Once = make(map[int]*sync.Once)
	statuses  map[int]bool       = make(map[int]bool)
)

const (
	twitterDMRoutineKey = iota
	twitterUserRoutineKey
	slackRoutineKey
	twitterPeriodicRoutineKey
	httpRoutineKey
)

func init() {
	executers[twitterDMRoutineKey] = new(sync.Once)
	executers[twitterUserRoutineKey] = new(sync.Once)
	executers[slackRoutineKey] = new(sync.Once)
	executers[twitterPeriodicRoutineKey] = new(sync.Once)
	executers[httpRoutineKey] = new(sync.Once)

	statuses[twitterDMRoutineKey] = false
	statuses[twitterUserRoutineKey] = false
	statuses[slackRoutineKey] = false
	statuses[twitterPeriodicRoutineKey] = false
	statuses[httpRoutineKey] = false
}

func main() {
	home, err := homedir.Dir()
	exitIfError(err)

	configDir := filepath.Join(home, ".config", "mybot")
	cacheDir := filepath.Join(home, ".cache", "mybot")

	envFlag := cli.StringFlag{
		Name:   "env",
		Value:  "",
		Usage:  `Assign "production" for production environment`,
		EnvVar: "MYBOT_ENV",
	}

	logFlag := cli.StringFlag{
		Name:   "log",
		Value:  filepath.Join(cacheDir, "mybot.log"),
		Usage:  "Log file's location",
		EnvVar: "MYBOT_LOG_PATH",
	}

	configFlag := cli.StringFlag{
		Name:   "config",
		Value:  filepath.Join(configDir, "config.toml"),
		Usage:  "Config file's location",
		EnvVar: "MYBOT_CONFIG_PATH",
	}

	cacheFlag := cli.StringFlag{
		Name:   "cache",
		Value:  filepath.Join(cacheDir, "cache.json"),
		Usage:  "Cache file's location",
		EnvVar: "MYBOT_CACHE_PATH",
	}

	gcloudFlag := cli.StringFlag{
		Name:   "gcloud",
		Value:  filepath.Join(configDir, "google_application_credentials.json"),
		Usage:  "Credential file for Google Cloud Platform",
		EnvVar: "MYBOT_GCLOUD_CREDENTIAL",
	}

	twitterFlag := cli.StringFlag{
		Name:   "twitter",
		Value:  filepath.Join(configDir, "twitter_authentication.toml"),
		Usage:  "Credential file for Twitter API",
		EnvVar: "MYBOT_TWITTER_CREDENTIAL",
	}

	certFlag := cli.StringFlag{
		Name:   "cert",
		Value:  filepath.Join(configDir, "mybot.crt"),
		Usage:  "Certification file for server",
		EnvVar: "MYBOT_SSL_CERT",
	}

	keyFlag := cli.StringFlag{
		Name:   "key",
		Value:  filepath.Join(configDir, "mybot.key"),
		Usage:  "Key file for server",
		EnvVar: "MYBOT_SSL_KEY",
	}

	hostFlag := cli.StringFlag{
		Name:   "host,H",
		Value:  "localhost",
		Usage:  "Host this server listen on",
		EnvVar: "MYBOT_HOST",
	}

	portFlag := cli.StringFlag{
		Name:   "port,P",
		Value:  "8080",
		Usage:  "Port this server listen on",
		EnvVar: "MYBOT_PORT",
	}

	dbAddrFlag := cli.StringFlag{
		Name:   "db-addr",
		Value:  "",
		Usage:  "DB address",
		EnvVar: "MYBOT_DB_ADDRESS",
	}

	dbUserFlag := cli.StringFlag{
		Name:   "db-user",
		Value:  "",
		Usage:  "DB user",
		EnvVar: "MYBOT_DB_USER",
	}

	dbPassFlag := cli.StringFlag{
		Name:   "db-passwd",
		Value:  "",
		Usage:  "DB password",
		EnvVar: "MYBOT_DB_PASSWD",
	}

	dbNameFlag := cli.StringFlag{
		Name:   "db-name",
		Value:  "",
		Usage:  "Target DB name",
		EnvVar: "MYBOT_DB_NAME",
	}

	slackTokenFlag := cli.StringFlag{
		Name:   "slack-token",
		Value:  "",
		Usage:  "Slack bot Token",
		EnvVar: "MYBOT_SLACK_TOKEN",
	}

	twitterConsumerKeyFlag := cli.StringFlag{
		Name:   "twitter-consumer-key",
		Value:  "",
		Usage:  "Twitter consumer key",
		EnvVar: "MYBOT_TWITTER_CONSUMER_KEY",
	}

	twitterConsumerSecretFlag := cli.StringFlag{
		Name:   "twitter-consumer-secret",
		Value:  "",
		Usage:  "Twitter consumer secret",
		EnvVar: "MYBOT_TWITTER_CONSUMER_SECRET",
	}

	twitterConsumerFileFlag := cli.StringFlag{
		Name:   "twitter-consumer-file",
		Value:  filepath.Join(configDir, "twitter_consumer_credentials.toml"),
		Usage:  "Twitter consumer file",
		EnvVar: "MYBOT_TWITTER_CONSUMER_FILE",
	}

	runFlags := []cli.Flag{
		envFlag,
		logFlag,
		configFlag,
		cacheFlag,
		gcloudFlag,
		twitterFlag,
		dbAddrFlag,
		dbUserFlag,
		dbPassFlag,
		dbNameFlag,
		slackTokenFlag,
		twitterConsumerKeyFlag,
		twitterConsumerSecretFlag,
		twitterConsumerFileFlag,
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
	exitIfError(err)
}

func beforeRunning(c *cli.Context) error {
	err := beforeValidate(c)
	exitIfError(err)

	slackToken := c.String("slack-token")
	slackAPI = mybot.NewSlackAPI(slackToken, config, cache)

	if info, err := os.Stat(c.String("gcloud")); err == nil && !info.IsDir() {
		visionAPI, err = mybot.NewVisionAPI(c.String("gcloud"))
		exitIfError(err)
		languageAPI, err = mybot.NewLanguageAPI(c.String("gcloud"))
		exitIfError(err)
	} else {
		visionAPI = &mybot.VisionAPI{}
		languageAPI = &mybot.LanguageAPI{}
	}

	return nil
}

func beforeValidate(c *cli.Context) error {
	var err error
	dbAddress := c.String("db-addr")
	dbUser := c.String("db-user")
	dbPasswd := c.String("db-passwd")
	dbName := c.String("db-name")

	var session *mgo.Session
	if dbAddress != "" && dbName != "" {
		info := &mgo.DialInfo{}
		info.Addrs = []string{dbAddress}
		info.Username = dbUser
		info.Password = dbPasswd
		info.Database = dbName
		session, err = mgo.DialWithInfo(info)
		exitIfError(err)
	}

	if session == nil {
		cache, err = mybot.NewFileCache(c.String("cache"))
	} else {
		col := session.DB(dbName).C("cache")
		cache, err = mybot.NewDBCache(col)
	}
	exitIfError(err)

	if session == nil {
		config, err = mybot.NewFileConfig(c.String("config"))
	} else {
		col := session.DB(dbName).C("config")
		config, err = mybot.NewDBConfig(col)
	}
	exitIfError(err)

	ck := c.String("twitter-consumer-key")
	cs := c.String("twitter-consumer-secret")
	cFile := c.String("twitter-consumer-file")
	if session == nil {
		twitterApp, err = mybot.NewFileTwitterOAuthApp(cFile)
	} else {
		col := session.DB(dbName).C("twitter-app-auth")
		twitterApp, err = mybot.NewDBTwitterOAuthApp(col)
	}
	exitIfError(err)
	if ck != "" && cs != "" {
		twitterApp.SetCreds(ck, cs)
		err := twitterApp.Save()
		exitIfError(err)
	}

	if session == nil {
		twitterAuth, err = mybot.NewFileOAuthCreds(c.String("twitter"))
	} else {
		col := session.DB(dbName).C("twitter-user-auth")
		twitterAuth, err = mybot.NewDBOAuthCreds(col)
	}
	exitIfError(err)

	twitterAPI = mybot.NewTwitterAPI(twitterAuth, cache, config)

	return nil
}

func run(c *cli.Context) {
	if err := runTwitterWithoutStream(); err != nil {
		log.Print(err)
		return
	}
	if err := cache.Save(); err != nil {
		log.Print(err)
		return
	}
}

func twitterListenDM() {
	key := twitterDMRoutineKey
	defer func() { executers[key] = new(sync.Once) }()
	f := func() {
		if !twitterAPIIsAvailable() {
			return
		}

		statuses[key] = true
		defer func() { statuses[key] = false }()

		r := twitterAPI.DefaultDirectMessageReceiver
		listener, err := twitterAPI.ListenMyself(nil, r)
		if err != nil {
			log.Print(err)
			return
		}
		dmListenerStream = listener.Stream
		if err := listener.Listen(); err != nil {
			log.Print(err)
			return
		}
		log.Print("Failed to keep connection")
	}
	executers[key].Do(f)
}

func twitterListenUsers() {
	key := twitterUserRoutineKey
	defer func() { executers[key] = new(sync.Once) }()
	f := func() {
		if !twitterAPIIsAvailable() {
			return
		}

		statuses[key] = true
		defer func() { statuses[key] = false }()

		listener, err := twitterAPI.ListenUsers(nil)
		if err != nil {
			log.Print(err)
			return
		}
		userListenerStream = listener.Stream
		if err := listener.Listen(visionAPI, languageAPI, slackAPI, cache); err != nil {
			log.Print(err)
			return
		}
		log.Print("Failed to keep connection")
	}
	executers[key].Do(f)
}

func twitterPeriodically() {
	key := twitterPeriodicRoutineKey
	defer func() { executers[key] = new(sync.Once) }()
	f := func() {
		if !twitterAPIIsAvailable() {
			return
		}

		statuses[key] = true
		defer func() { statuses[key] = false }()

		for {
			if err := runTwitterWithStream(); err != nil {
				log.Print(err)
				return
			}
			if err := cache.Save(); err != nil {
				log.Print(err)
				return
			}
			d, err := time.ParseDuration(config.GetTwitterDuration())
			if err != nil {
				log.Print(err)
				return
			}
			time.Sleep(d)
		}
	}
	executers[key].Do(f)
}

func twitterAPIIsAvailable() bool {
	if twitterAPI == nil {
		return false
	} else if success, err := twitterAPI.VerifyCredentials(); !success || err != nil {
		return false
	}
	return true
}

func slackListens() {
	key := slackRoutineKey
	defer func() { executers[key] = new(sync.Once) }()
	f := func() {
		if slackAPI == nil || !slackAPI.Enabled() {
			return
		}

		statuses[key] = true
		defer func() { statuses[key] = false }()

		slackListener = slackAPI.Listen()
		if err := slackListener.Start(visionAPI, languageAPI, twitterAPI); err != nil {
			log.Print(err)
			return
		}
		log.Print("Failed to keep connection")
	}
	executers[key].Do(f)
}

func reloadListeners() {
	go twitterListenUsers()
	go twitterListenDM()
	go twitterPeriodically()
	go slackListens()
}

func httpServer(c *cli.Context) {
	key := httpRoutineKey
	defer func() { executers[key] = new(sync.Once) }()
	f := func() {
		statuses[key] = true
		defer func() { statuses[key] = false }()
		err := startServer(c.String("host"), c.String("port"), c.String("cert"), c.String("key"))
		exitIfError(err)
	}
	executers[httpRoutineKey].Do(f)
}

func serve(c *cli.Context) error {
	go httpServer(c)
	reloadListeners()

	ch := make(chan bool)
	<-ch
	return nil
}

func runTwitterWithStream() error {
	tweets := []anaconda.Tweet{}
	for _, a := range config.GetTwitterSearches() {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if len(a.ResultType) != 0 {
			v.Set("result_type", a.ResultType)
		}
		for _, query := range a.Queries {
			ts, err := twitterAPI.ProcessSearch(
				query,
				v,
				a.Filter,
				visionAPI,
				languageAPI,
				slackAPI,
				a.Action,
			)
			if err != nil {
				return err
			}
			tweets = append(tweets, ts...)
		}
	}
	for _, a := range config.GetTwitterFavorites() {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		for _, name := range a.ScreenNames {
			ts, err := twitterAPI.ProcessFavorites(
				name,
				v,
				a.Filter,
				visionAPI,
				languageAPI,
				slackAPI,
				a.Action,
			)
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
	for _, a := range config.GetTwitterTimelines() {
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
			ts, err := twitterAPI.ProcessTimeline(
				name,
				v,
				a.Filter,
				visionAPI,
				languageAPI,
				slackAPI,
				a.Action,
			)
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

func validate(c *cli.Context) {
	err := config.Validate()
	exitIfError(err)
	if c.Bool("api") {
		err := config.ValidateWithAPI(twitterAPI)
		exitIfError(err)
	}
}

func exitIfError(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
