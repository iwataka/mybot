package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/lib"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var (
	twitterAPI         *mybot.TwitterAPI
	userListenerStream *anaconda.Stream
	dmListenerStream   *anaconda.Stream
	twitterApp         *mybot.OAuthApp
	twitterAuth        *mybot.OAuthCredentials
	visionAPI          *mybot.VisionAPI
	languageAPI        *mybot.LanguageAPI
	slackAPI           *mybot.SlackAPI
	slackListener      *mybot.SlackListener
	config             *mybot.FileConfig
	cache              mybot.Cache
	status             *mybot.Status
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

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
		Value:  "",
		Usage:  "Host this server listen on",
		EnvVar: "MYBOT_HOST",
	}

	portFlag := cli.StringFlag{
		Name:   "port,P",
		Value:  "",
		Usage:  "Port this server listen on",
		EnvVar: "MYBOT_PORT",
	}

	driverNameFlag := cli.StringFlag{
		Name:   "driver-name",
		Value:  "",
		Usage:  "Driver name for DB",
		EnvVar: "MYBOT_DRIVER_NAME",
	}

	dataSourceFlag := cli.StringFlag{
		Name:   "data-source",
		Value:  "",
		Usage:  "Data Source for DB",
		EnvVar: "MYBOT_DATA_SOURCE",
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

	runFlags := []cli.Flag{
		envFlag,
		logFlag,
		configFlag,
		cacheFlag,
		gcloudFlag,
		twitterFlag,
		driverNameFlag,
		dataSourceFlag,
		slackTokenFlag,
		twitterConsumerKeyFlag,
		twitterConsumerSecretFlag,
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
	err := beforeValidate(c)
	if err != nil {
		panic(err)
	}

	slackToken := c.String("slack-token")
	slackAPI = mybot.NewSlackAPI(slackToken, config, cache)

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

	status = mybot.NewStatus()
	status.UpdateTwitterAuth(twitterAPI)

	return nil
}

func beforeValidate(c *cli.Context) error {
	var err error
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	if c.String("env") == "production" {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	driverName := c.String("driver-name")
	dataSource := c.String("data-source")

	if driverName == "" || dataSource == "" {
		cache, err = mybot.NewFileCache(c.String("cache"))
	} else {
		cache, err = mybot.NewDBCache(driverName, dataSource)
	}
	if err != nil {
		panic(err)
	}

	config, err = mybot.NewFileConfig(c.String("config"))
	if err != nil {
		panic(err)
	}

	twitterApp = &mybot.OAuthApp{
		c.String("twitter-consumer-key"),
		c.String("twitter-consumer-secret"),
		"",
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

func run(c *cli.Context) {
	logFields := log.Fields{
		"type": "twitter",
	}

	if err := runTwitterWithoutStream(); err != nil {
		log.WithFields(logFields).Error(err)
		return
	}
	if err := cache.Save(); err != nil {
		log.WithFields(logFields).Error(err)
		return
	}
}

func twitterListenDM() {
	if !status.PassTwitterAuth {
		return
	}

	status.LockListenDMRoutine()
	defer status.UnlockListenDMRoutine()

	logFields := log.Fields{
		"type": "twitter",
	}

	r := twitterAPI.DefaultDirectMessageReceiver
	listener, err := twitterAPI.ListenMyself(nil, r)
	if err != nil {
		log.WithFields(logFields).Error(err)
		return
	}
	dmListenerStream = listener.Stream
	if err := listener.Listen(); err != nil {
		log.WithFields(logFields).Error(err)
		return
	}
	log.WithFields(logFields).Error("Failed to keep connection")
}

func twitterListenUsers() {
	if !status.PassTwitterAuth {
		return
	}

	status.LockListenUsersRoutine()
	defer status.UnlockListenUsersRoutine()

	logFields := log.Fields{
		"type": "twitter",
	}

	listener, err := twitterAPI.ListenUsers(nil)
	if err != nil {
		log.WithFields(logFields).Error(err)
		return
	}
	userListenerStream = listener.Stream
	if err := listener.Listen(visionAPI, languageAPI, slackAPI, cache); err != nil {
		log.WithFields(logFields).Error(err)
		return
	}
	log.WithFields(logFields).Error("Failed to keep connection")
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

	logFields := log.Fields{
		"type": "twitter",
	}

	for {
		if err := runTwitterWithStream(); err != nil {
			log.WithFields(logFields).Error(err)
			return
		}
		if err := cache.Save(); err != nil {
			log.WithFields(logFields).Error(err)
			return
		}
		d, err := time.ParseDuration(config.Twitter.Duration)
		if err != nil {
			log.WithFields(logFields).Error(err)
			return
		}
		time.Sleep(d)
	}
}

func slackListens() {
	if !slackAPI.Enabled() {
		return
	}

	status.LockSlackListenRoutine()
	defer status.UnlockSlackListenRoutine()

	logFields := log.Fields{
		"type": "slack",
	}

	slackListener = slackAPI.Listen()
	if err := slackListener.Start(visionAPI, languageAPI, twitterAPI); err != nil {
		log.WithFields(logFields).Error(err)
		return
	}
	log.WithFields(logFields).Error("Failed to keep connection")
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

	if dmListenerStream != nil {
		dmListenerStream.Stop()
	}
	go twitterListenDM()

	if !status.TwitterStatus {
		go twitterPeriodically()
	}

	if slackListener != nil {
		slackListener.Stop()
	}
	go slackListens()
}

func httpServer(c *cli.Context) {
	if status.ServerStatus {
		return
	}
	status.ServerStatus = true
	defer func() { status.ServerStatus = false }()
	host := c.String("host")
	port := c.String("port")
	cert := c.String("cert")
	key := c.String("key")
	err := startServer(host, port, cert, key)
	if err != nil {
		panic(err)
	}
}

func serve(c *cli.Context) error {
	go httpServer(c)
	go twitterListenDM()
	go twitterListenUsers()
	go twitterPeriodically()
	go slackListens()

	ch := make(chan bool)
	<-ch
	return nil
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
	for _, a := range config.Twitter.Favorites {
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
	if err != nil {
		panic(err)
	}
	if c.Bool("api") {
		err := config.ValidateWithAPI(twitterAPI)
		if err != nil {
			panic(err)
		}
	}
}
