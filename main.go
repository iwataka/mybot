package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"log"

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
	userSpecificDataMap map[string]*userSpecificData = make(map[string]*userSpecificData)

	// Global-scope data
	twitterApp  mybot.OAuthApp
	slackApp    mybot.OAuthApp
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	cliContext  *cli.Context
	session     *mgo.Session
)

const (
	twitterDMRoutineKey = iota
	twitterUserRoutineKey
	slackRoutineKey
	twitterPeriodicRoutineKey
)

const (
	startSignal = iota
	restartSignal
	stopSignal
	killSignal
)

type userSpecificData struct {
	config      mybot.Config
	cache       mybot.Cache
	twitterAPI  *mybot.TwitterAPI
	twitterAuth mybot.OAuthCreds
	slackAPI    *mybot.SlackAPI
	slackAuth   mybot.OAuthCreds
	workerChans map[int]chan int
	statuses    map[int]*bool
}

func newUserSpecificData(c *cli.Context, session *mgo.Session, userID string) (*userSpecificData, error) {
	var err error
	data := &userSpecificData{}
	data.workerChans = map[int]chan int{}
	data.statuses = map[int]*bool{}
	initStatuses(data.statuses)
	dbName := c.String("db-name")

	if session == nil {
		dir, err := argValueWithMkdir(c, "cache")
		if err != nil {
			return nil, err
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		data.cache, err = mybot.NewFileCache(file)
	} else {
		col := session.DB(dbName).C("cache")
		data.cache, err = mybot.NewDBCache(col, userID)
	}
	if err != nil {
		return nil, err
	}

	if session == nil {
		dir, err := argValueWithMkdir(c, "config")
		if err != nil {
			return nil, err
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		data.config, err = mybot.NewFileConfig(file)
	} else {
		col := session.DB(dbName).C("config")
		data.config, err = mybot.NewDBConfig(col, userID)
	}
	if err != nil {
		return nil, err
	}

	if session == nil {
		dir, err := argValueWithMkdir(c, "twitter-auth")
		if err != nil {
			return nil, err
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		data.twitterAuth, err = mybot.NewFileOAuthCreds(file)
	} else {
		col := session.DB(dbName).C("twitter-user-auth")
		data.twitterAuth, err = mybot.NewDBOAuthCreds(col, userID)
	}
	if err != nil {
		return nil, err
	}

	if session == nil {
		dir, err := argValueWithMkdir(c, "slack-auth")
		if err != nil {
			return nil, err
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		data.slackAuth, err = mybot.NewFileOAuthCreds(file)
	} else {
		col := session.DB(dbName).C("slack-user-auth")
		data.slackAuth, err = mybot.NewDBOAuthCreds(col, userID)
	}
	if err != nil {
		return nil, err
	}

	data.twitterAPI = mybot.NewTwitterAPI(data.twitterAuth, data.cache, data.config)

	slackId, _ := data.slackAuth.GetCreds()
	data.slackAPI = mybot.NewSlackAPI(slackId, data.config, data.cache)

	go manageWorkerWithStart(
		twitterDMRoutineKey,
		userID,
		newTwitterDMWorker(data.twitterAPI, data.statuses[twitterDMRoutineKey]),
	)
	go manageWorkerWithStart(
		twitterUserRoutineKey,
		userID,
		newTwitterUserWorker(data.twitterAPI, data.slackAPI, visionAPI, languageAPI, data.cache, data.statuses[twitterUserRoutineKey]),
	)
	go manageWorkerWithStart(
		twitterPeriodicRoutineKey,
		userID,
		newTwitterPeriodicWorker(data.twitterAPI, data.slackAPI, visionAPI, languageAPI, data.cache, data.config, data.statuses[twitterPeriodicRoutineKey]),
	)
	go manageWorkerWithStart(
		slackRoutineKey,
		userID,
		newSlackWorker(data.slackAPI, data.twitterAPI, visionAPI, languageAPI, data.statuses[slackRoutineKey]),
	)

	return data, nil
}

func initStatuses(statuses map[int]*bool) {
	keys := []int{twitterDMRoutineKey, twitterUserRoutineKey, twitterPeriodicRoutineKey, slackRoutineKey}
	for _, key := range keys {
		val := false
		statuses[key] = &val
	}
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
		Before:  beforeAll,
		Action:  validate,
	}

	app.Commands = []cli.Command{runCmd, serveCmd, validateCmd}
	err = app.Run(os.Args)
	exitIfError(err)
}

func beforeAll(c *cli.Context) error {
	cliContext = c
	dbAddress := c.String("db-addr")
	dbUser := c.String("db-user")
	dbPasswd := c.String("db-passwd")
	dbName := c.String("db-name")
	var err error

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

	userIDs, err := getUserIDs(c, session, dbName)
	if err != nil {
		return err
	}
	for _, userID := range userIDs {
		err := initForUser(c, session, dbName, userID)
		log.Printf("Initialize for user %s", userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func initForUser(c *cli.Context, session *mgo.Session, dbName, userID string) error {
	data, err := newUserSpecificData(c, session, userID)
	if err != nil {
		return err
	}
	userSpecificDataMap[userID] = data
	return nil
}

func beforeRunning(c *cli.Context) error {
	err := beforeAll(c)
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

func getUserIDs(c *cli.Context, session *mgo.Session, dbName string) ([]string, error) {
	if session == nil {
		dir, err := argValueWithMkdir(c, "twitter-auth")
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

func run(c *cli.Context) {
	for _, data := range userSpecificDataMap {
		if err := runTwitterWithoutStream(data.twitterAPI, data.slackAPI, data.config); err != nil {
			log.Print(err)
			return
		}
		if err := data.cache.Save(); err != nil {
			log.Print(err)
			return
		}
	}
}

func serve(c *cli.Context) error {
	go httpServer(c)
	ch := make(chan bool)
	<-ch
	return nil
}

func httpServer(c *cli.Context) {
	err := startServer(c.String("host"), c.String("port"), c.String("cert"), c.String("key"))
	exitIfError(err)
}

func runTwitterWithStream(
	twitterAPI *mybot.TwitterAPI,
	slackAPI *mybot.SlackAPI,
	visionAPI *mybot.VisionAPI,
	languageAPI *mybot.LanguageAPI,
	config mybot.Config,
) error {
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
	err := runTwitterWithStream(twitterAPI, slackAPI, visionAPI, languageAPI, config)
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
	for _, data := range userSpecificDataMap {
		err := data.config.Validate()
		exitIfError(err)
		if c.Bool("api") {
			err := data.config.ValidateWithAPI(data.twitterAPI)
			exitIfError(err)
		}
	}
}

func exitIfError(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func argValueWithMkdir(c *cli.Context, key string) (string, error) {
	dir := c.String(key)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func twitterAPIIsAvailable(twitterAPI *mybot.TwitterAPI) bool {
	if twitterAPI == nil {
		return false
	} else if success, err := twitterAPI.VerifyCredentials(); !success || err != nil {
		return false
	}
	return true
}
