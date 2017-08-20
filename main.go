package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/lib"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"gopkg.in/mgo.v2"
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
	status             *mybot.Status
	twitterApp         mybot.OAuthApp
	twitterAuth        mybot.OAuthCreds
	config             mybot.Config
	cache              mybot.Cache
)

func main() {
	home, err := homedir.Dir()
	fatalIfError(err)

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
	fatalIfError(err)
}

func beforeRunning(c *cli.Context) error {
	err := beforeValidate(c)
	fatalIfError(err)

	slackToken := c.String("slack-token")
	slackAPI = mybot.NewSlackAPI(slackToken, config, cache)

	if info, err := os.Stat(c.String("gcloud")); err == nil && !info.IsDir() {
		visionAPI, err = mybot.NewVisionAPI(c.String("gcloud"))
		fatalIfError(err)
		languageAPI, err = mybot.NewLanguageAPI(c.String("gcloud"))
		fatalIfError(err)
	} else {
		visionAPI = &mybot.VisionAPI{}
		languageAPI = &mybot.LanguageAPI{}
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

	dbAddress := c.String("db-addr")
	dbUser := c.String("db-user")
	dbPasswd := c.String("db-passwd")
	dbName := c.String("db-name")

	var session *mgo.Session
	if dbAddress != "" && dbUser != "" && dbPasswd != "" && dbName != "" {
		info := &mgo.DialInfo{}
		info.Addrs = []string{dbAddress}
		info.Username = dbUser
		info.Password = dbPasswd
		info.Database = dbName
		session, err = mgo.DialWithInfo(info)
		fatalIfError(err)
	}

	if session == nil {
		cache, err = mybot.NewFileCache(c.String("cache"))
	} else {
		col := session.DB(dbName).C("cache")
		cache, err = mybot.NewDBCache(col)
	}
	fatalIfError(err)

	if session == nil {
		config, err = mybot.NewFileConfig(c.String("config"))
	} else {
		col := session.DB(dbName).C("config")
		config, err = mybot.NewDBConfig(col)
	}
	fatalIfError(err)

	ck := c.String("twitter-consumer-key")
	cs := c.String("twitter-consumer-secret")
	cFile := c.String("twitter-consumer-file")
	if session == nil {
		twitterApp, err = mybot.NewFileTwitterOAuthApp(cFile)
	} else {
		col := session.DB(dbName).C("twitter-app-auth")
		twitterApp, err = mybot.NewDBTwitterOAuthApp(col)
	}
	fatalIfError(err)
	if ck != "" && cs != "" {
		twitterApp.SetCreds(ck, cs)
		err := twitterApp.Save()
		fatalIfError(err)
	}

	if session == nil {
		twitterAuth, err = mybot.NewFileOAuthCreds(c.String("twitter"))
	} else {
		col := session.DB(dbName).C("twitter-user-auth")
		twitterAuth, err = mybot.NewDBOAuthCreds(col)
	}
	fatalIfError(err)

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
		d, err := time.ParseDuration(config.GetTwitterDuration())
		if err != nil {
			log.WithFields(logFields).Error(err)
			return
		}
		time.Sleep(d)
	}
}

func slackListens() {
	if slackAPI == nil || !slackAPI.Enabled() {
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

	go twitterPeriodically()

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
	fatalIfError(err)
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
	fatalIfError(err)
	if c.Bool("api") {
		err := config.ValidateWithAPI(twitterAPI)
		fatalIfError(err)
	}
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
