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
	config             mybot.Config
	cache              mybot.Cache
	twitterAPI         *mybot.TwitterAPI
	twitterAuth        mybot.OAuthCreds
	slackAPI           *mybot.SlackAPI
	slackAuth          mybot.OAuthCreds
	slackListener      *mybot.SlackListener
	userListenerStream *anaconda.Stream
	dmListenerStream   *anaconda.Stream

	// Global-scope data
	twitterApp  mybot.OAuthApp
	slackApp    mybot.OAuthApp
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	executers   map[int]*sync.Once = make(map[int]*sync.Once)
	statuses    map[int]bool       = make(map[int]bool)
)

const defaultUserID = "mybot-myself"

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

	configFlag := cli.StringFlag{
		Name:   "config",
		Value:  filepath.Join(configDir, "config"),
		Usage:  "Config directory location",
		EnvVar: "MYBOT_CONFIG_PATH",
	}

	cacheFlag := cli.StringFlag{
		Name:   "cache",
		Value:  filepath.Join(cacheDir, "cache"),
		Usage:  "Cache directory location",
		EnvVar: "MYBOT_CACHE_PATH",
	}

	gcloudFlag := cli.StringFlag{
		Name:   "gcloud",
		Value:  filepath.Join(configDir, "google_application_credentials.json"),
		Usage:  "Credential file for Google Cloud Platform",
		EnvVar: "MYBOT_GCLOUD_CREDENTIAL",
	}

	twitterFlag := cli.StringFlag{
		Name:   "twitter-auth",
		Value:  filepath.Join(configDir, "twitter_auth"),
		Usage:  "Twitter credential directory",
		EnvVar: "MYBOT_TWITTER_CREDENTIAL",
	}

	slackFlag := cli.StringFlag{
		Name:   "slack-auth",
		Value:  filepath.Join(configDir, "slack_auth"),
		Usage:  "Slack credential directory",
		EnvVar: "MYBOT_SLACK_CREDENTIAL",
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
		Name:   "twitter-app",
		Value:  filepath.Join(configDir, "twitter_app.toml"),
		Usage:  "Twitter application directory",
		EnvVar: "MYBOT_TWITTER_APP",
	}

	slackClientIDFlag := cli.StringFlag{
		Name:   "slack-client-id",
		Value:  "",
		Usage:  "Slack client ID",
		EnvVar: "MYBOT_SLACK_APP",
	}

	slackClientSecretFlag := cli.StringFlag{
		Name:   "slack-client-secret",
		Value:  "",
		Usage:  "Slack client secret",
		EnvVar: "MYBOT_SLACK_CLIENT_SECRET",
	}

	slackClientFileFlag := cli.StringFlag{
		Name:   "slack-app",
		Value:  filepath.Join(configDir, "slack_app.toml"),
		Usage:  "slack application directory",
		EnvVar: "MYBOT_SLACK_APP",
	}

	runFlags := []cli.Flag{
		envFlag,
		configFlag,
		cacheFlag,
		gcloudFlag,
		twitterFlag,
		slackFlag,
		dbAddrFlag,
		dbUserFlag,
		dbPassFlag,
		dbNameFlag,
		twitterConsumerKeyFlag,
		twitterConsumerSecretFlag,
		twitterConsumerFileFlag,
		slackClientIDFlag,
		slackClientSecretFlag,
		slackClientFileFlag,
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
		var err error
		session, err = mgo.DialWithInfo(info)
		if err != nil {
			return err
		}
	}

	// userIDs, err := getUserIDs(c, session, dbName)
	// if err != nil {
	// 	return err
	// }
	// for _, userID := range userIDs {
	// 	err := initForUser(c, session, dbName, userID)
	// 	log.Printf("Initialize for user %s", userID)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// return nil
	return initForUser(c, session, dbName, defaultUserID)
}

func getUserIDs(c *cli.Context, session *mgo.Session, dbName string) ([]string, error) {
	if session == nil {
		dir, err := getArgDir(c, "twitter-auth")
		if err != nil {
			return nil, err
		}
		files, err := filepath.Glob(filepath.Join(dir, "*.toml"))
		userIDs := []string{}
		for _, file := range files {
			base := filepath.Base(file)
			ext := filepath.Ext(file)
			userIDs = append(userIDs, base[0:len(base)-len(ext)])
		}
		return userIDs, nil
	} else {
		col := session.DB(dbName).C("twitter-user-auth")
		auths := []map[string]interface{}{}
		err := col.Find(nil).All(&auths)
		if err != nil {
			return nil, err
		}
		userIDs := []string{}
		for _, auth := range auths {
			id, ok := auth["id"].(string)
			if ok && id != "" {
				userIDs = append(userIDs, id)
			}
		}
		return userIDs, nil
	}
}

func initForUser(c *cli.Context, session *mgo.Session, dbName, userID string) error {
	var err error

	if session == nil {
		dir, err := getArgDir(c, "cache")
		if err != nil {
			return err
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		cache, err = mybot.NewFileCache(file)
	} else {
		col := session.DB(dbName).C("cache")
		cache, err = mybot.NewDBCache(col, userID)
	}
	if err != nil {
		return err
	}

	if session == nil {
		dir, err := getArgDir(c, "config")
		if err != nil {
			return err
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		config, err = mybot.NewFileConfig(file)
	} else {
		col := session.DB(dbName).C("config")
		config, err = mybot.NewDBConfig(col, userID)
	}
	if err != nil {
		return err
	}

	twitterCk := c.String("twitter-consumer-key")
	twitterCs := c.String("twitter-consumer-secret")
	if session == nil {
		twitterApp, err = mybot.NewFileTwitterOAuthApp(c.String("twitter-app"))
	} else {
		col := session.DB(dbName).C("twitter-app-auth")
		twitterApp, err = mybot.NewDBTwitterOAuthApp(col)
	}
	if err != nil {
		return err
	}
	if twitterCk != "" && twitterCs != "" {
		twitterApp.SetCreds(twitterCk, twitterCs)
		err := twitterApp.Save()
		if err != nil {
			return err
		}
	}

	slackCk := c.String("slack-client-id")
	slackCs := c.String("slack-client-secret")
	if session == nil {
		slackApp, err = mybot.NewFileOAuthApp(c.String("slack-app"))
	} else {
		col := session.DB(dbName).C("slack-app-auth")
		slackApp, err = mybot.NewDBOAuthApp(col)
	}
	if err != nil {
		return err
	}
	if slackCk != "" && slackCs != "" {
		slackApp.SetCreds(slackCk, slackCs)
		err := slackApp.Save()
		if err != nil {
			return err
		}
	}

	if session == nil {
		dir, err := getArgDir(c, "twitter-auth")
		if err != nil {
			return err
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		twitterAuth, err = mybot.NewFileOAuthCreds(file)
	} else {
		col := session.DB(dbName).C("twitter-user-auth")
		twitterAuth, err = mybot.NewDBOAuthCreds(col, userID)
	}
	if err != nil {
		return err
	}

	if session == nil {
		dir, err := getArgDir(c, "slack-auth")
		if err != nil {
			return err
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		slackAuth, err = mybot.NewFileOAuthCreds(file)
	} else {
		col := session.DB(dbName).C("slack-user-auth")
		slackAuth, err = mybot.NewDBOAuthCreds(col, userID)
	}
	if err != nil {
		return err
	}

	twitterAPI = mybot.NewTwitterAPI(twitterAuth, cache, config)

	slackId, _ := slackAuth.GetCreds()
	slackAPI = mybot.NewSlackAPI(slackId, config, cache)

	return nil
}

func run(c *cli.Context) {
	if err := runTwitterWithoutStream(twitterAPI, slackAPI, config); err != nil {
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
		if !twitterAPIIsAvailable(twitterAPI) {
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
		if !twitterAPIIsAvailable(twitterAPI) {
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
		if !twitterAPIIsAvailable(twitterAPI) {
			return
		}

		statuses[key] = true
		defer func() { statuses[key] = false }()

		for {
			if err := runTwitterWithStream(twitterAPI, slackAPI, config); err != nil {
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

func twitterAPIIsAvailable(twitterAPI *mybot.TwitterAPI) bool {
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

func runTwitterWithStream(twitterAPI *mybot.TwitterAPI, slackAPI *mybot.SlackAPI, config mybot.Config) error {
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

func runTwitterWithoutStream(twitterAPI *mybot.TwitterAPI, slackAPI *mybot.SlackAPI, config mybot.Config) error {
	err := runTwitterWithStream(twitterAPI, slackAPI, config)
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

func getArgDir(c *cli.Context, key string) (string, error) {
	dir := c.String(key)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return "", err
	}
	return dir, nil
}
