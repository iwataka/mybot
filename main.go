package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/sessions"
	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/worker"
	"github.com/kidstuff/mongostore"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"gopkg.in/mgo.v2"
)

//go:generate go-bindata assets/...
//go:generate mockgen -source=models/slack.go -destination=mocks/slack.go -package=mocks
//go:generate mockgen -source=models/twitter.go -destination=mocks/twitter.go -package=mocks
//go:generate mockgen -source=models/auth.go -destination=mocks/auth.go -package=mocks
//go:generate mockgen -source=lib/vision.go -destination=mocks/vision.go -package=mocks
//go:generate mockgen -source=lib/language.go -destination=mocks/language.go -package=mocks
//go:generate mockgen -source=lib/batch.go -destination=mocks/batch.go -package=mocks
//go:generate mockgen -source=lib/utils.go -destination=mocks/utils.go -package=mocks

var (
	userSpecificDataMap map[string]*userSpecificData = make(map[string]*userSpecificData)

	// Global-scope data
	twitterApp               mybot.OAuthApp
	slackApp                 mybot.OAuthApp
	visionAPI                mybot.VisionMatcher
	languageAPI              mybot.LanguageMatcher
	cliContext               *cli.Context
	dbSession                *mgo.Session
	serverSession            sessions.Store
	sessionDomain            string
	accessControlAllowOrigin string
)

const (
	twitterDMRoutineKey = iota
	twitterUserRoutineKey
	slackRoutineKey
	twitterPeriodicRoutineKey
)

type userSpecificData struct {
	config      mybot.Config
	cache       mybot.Cache
	twitterAPI  *mybot.TwitterAPI
	twitterAuth mybot.OAuthCreds
	slackAPI    *mybot.SlackAPI
	slackAuth   mybot.OAuthCreds
	workerChans map[int]chan *worker.WorkerSignal
	statuses    map[int]*bool
}

func newUserSpecificData(c *cli.Context, session *mgo.Session, userID string) (*userSpecificData, error) {
	var err error
	data := &userSpecificData{}
	data.workerChans = map[int]chan *worker.WorkerSignal{}
	data.statuses = map[int]*bool{}
	initStatuses(data.statuses)
	dbName := c.String("db-name")

	if session == nil {
		dir, err := argValueWithMkdir(c, "cache")
		if err != nil {
			return nil, mybot.WithStack(err)
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		data.cache, err = mybot.NewFileCache(file)
	} else {
		col := session.DB(dbName).C("cache")
		data.cache, err = mybot.NewDBCache(col, userID)
	}
	if err != nil {
		return nil, mybot.WithStack(err)
	}

	if session == nil {
		dir, err := argValueWithMkdir(c, "config")
		if err != nil {
			return nil, mybot.WithStack(err)
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		data.config, err = mybot.NewFileConfig(file)
	} else {
		col := session.DB(dbName).C("config")
		data.config, err = mybot.NewDBConfig(col, userID)
	}
	if err != nil {
		return nil, mybot.WithStack(err)
	}

	if session == nil {
		dir, err := argValueWithMkdir(c, "twitter-auth")
		if err != nil {
			return nil, mybot.WithStack(err)
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		data.twitterAuth, err = mybot.NewFileOAuthCreds(file)
	} else {
		col := session.DB(dbName).C("twitter-user-auth")
		data.twitterAuth, err = mybot.NewDBOAuthCreds(col, userID)
	}
	if err != nil {
		return nil, mybot.WithStack(err)
	}

	if session == nil {
		dir, err := argValueWithMkdir(c, "slack-auth")
		if err != nil {
			return nil, mybot.WithStack(err)
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
		data.slackAuth, err = mybot.NewFileOAuthCreds(file)
	} else {
		col := session.DB(dbName).C("slack-user-auth")
		data.slackAuth, err = mybot.NewDBOAuthCreds(col, userID)
	}
	if err != nil {
		return nil, mybot.WithStack(err)
	}

	data.twitterAPI = mybot.NewTwitterAPI(data.twitterAuth, data.cache, data.config)

	slackId, _ := data.slackAuth.GetCreds()
	data.slackAPI = mybot.NewSlackAPI(slackId, data.config, data.cache)

	return data, nil
}

func startUserSpecificData(userID string, data *userSpecificData) {
	manageWorkerWithStart(
		twitterDMRoutineKey,
		data.workerChans,
		data.statuses,
		newTwitterDMWorker(data.twitterAPI, userID, time.Minute),
	)
	manageWorkerWithStart(
		twitterUserRoutineKey,
		data.workerChans,
		data.statuses,
		newTwitterUserWorker(data.twitterAPI, data.slackAPI, visionAPI, languageAPI, data.cache, userID, time.Minute),
	)
	runner := mybot.NewBatchRunnerWithStream(data.twitterAPI, data.slackAPI, visionAPI, languageAPI, data.config)
	manageWorkerWithStart(
		twitterPeriodicRoutineKey,
		data.workerChans,
		data.statuses,
		newTwitterPeriodicWorker(runner, data.cache, data.config.GetTwitterDuration(), time.Minute, userID),
	)
	manageWorkerWithStart(
		slackRoutineKey,
		data.workerChans,
		data.statuses,
		newSlackWorker(data.slackAPI, data.twitterAPI, visionAPI, languageAPI, userID),
	)
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

	log.SetFlags(0)

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
		Usage:  "Slack application directory",
		EnvVar: "MYBOT_SLACK_APP",
	}

	sessionDomainFlag := cli.StringFlag{
		Name:        "session-domain",
		Value:       "",
		Usage:       "Session domain",
		EnvVar:      "MYBOT_SESSION_DOMAIN",
		Destination: &sessionDomain,
	}

	apiFlag := cli.BoolFlag{
		Name:  "api",
		Usage: "Use API to validate configuration",
	}

	accessControlAllowOriginFlag := cli.StringFlag{
		Name:        "access-control-allow-origin",
		Value:       "",
		Usage:       "Access Control Allow Origin value for API endpoints",
		EnvVar:      "MYBOT_ACCESS_CONTROL_ALLOW_ORIGIN",
		Destination: &accessControlAllowOrigin,
	}

	commonFlags := []cli.Flag{
		envFlag,
		configFlag,
		cacheFlag,
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

	runFlags := []cli.Flag{
		gcloudFlag,
	}

	for _, f := range commonFlags {
		runFlags = append(runFlags, f)
	}

	serveFlags := []cli.Flag{
		certFlag,
		keyFlag,
		hostFlag,
		portFlag,
		sessionDomainFlag,
		accessControlAllowOriginFlag,
	}

	// All `run` flags should be `serve` flag
	for _, f := range runFlags {
		serveFlags = append(serveFlags, f)
	}

	validateFlags := []cli.Flag{apiFlag}

	// All `run` flags should be `validate` flag
	for _, f := range commonFlags {
		validateFlags = append(validateFlags, f)
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
		Before:  beforeServing,
		Action:  serve,
	}

	validateCmd := cli.Command{
		Name:    "validate",
		Aliases: []string{"v"},
		Usage:   "Validates the user configuration",
		Flags:   validateFlags,
		Before:  beforeValidating,
		Action:  validate,
	}

	app.Commands = []cli.Command{runCmd, serveCmd, validateCmd}
	err = app.Run(os.Args)
	exitIfError(err)
}

func run(c *cli.Context) {
	for _, data := range userSpecificDataMap {
		baseRunner := mybot.NewBatchRunnerWithStream(data.twitterAPI, data.slackAPI, visionAPI, languageAPI, data.config)
		runner := mybot.NewBatchRunnerWithoutStream(baseRunner)
		if err := runner.Run(); err != nil {
			log.Printf("%+v\n", err)
			return
		}
		if err := data.cache.Save(); err != nil {
			log.Printf("%+v\n", err)
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

func beforeRunning(c *cli.Context) error {
	var err error
	visionAPI, err = mybot.NewVisionMatcher(c.String("gcloud"))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	languageAPI, err = mybot.NewLanguageMatcher(c.String("gcloud"))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = beforeValidating(c)
	exitIfError(err)
	return nil
}

func beforeServing(c *cli.Context) error {
	err := beforeRunning(c)
	exitIfError(err)
	for userID, data := range userSpecificDataMap {
		startUserSpecificData(userID, data)
	}
	return nil
}

func beforeValidating(c *cli.Context) error {
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
		dbSession, err = mgo.DialWithInfo(info)
		if err != nil {
			return mybot.WithStack(err)
		}
	}

	if dbSession == nil {
		sess := sessions.NewCookieStore(
			[]byte("mybot_session_key"),
		)
		if sessionDomain != "" {
			sess.Options.Domain = sessionDomain
		}
		serverSession = sess
	} else {
		sess := mongostore.NewMongoStore(
			dbSession.DB(dbName).C("user-session"),
			86400*30,
			true,
			[]byte("mybot_session_key"),
		)
		if sessionDomain != "" {
			sess.Options.Domain = sessionDomain
		}
		serverSession = sess
	}

	twitterCk := c.String("twitter-consumer-key")
	twitterCs := c.String("twitter-consumer-secret")
	if dbSession == nil {
		twitterApp, err = mybot.NewFileTwitterOAuthApp(c.String("twitter-app"))
	} else {
		col := dbSession.DB(dbName).C("twitter-app-auth")
		twitterApp, err = mybot.NewDBTwitterOAuthApp(col)
	}
	if err != nil {
		return mybot.WithStack(err)
	}
	if twitterCk != "" && twitterCs != "" {
		twitterApp.SetCreds(twitterCk, twitterCs)
		err := twitterApp.Save()
		if err != nil {
			return mybot.WithStack(err)
		}
	}

	slackCk := c.String("slack-client-id")
	slackCs := c.String("slack-client-secret")
	if dbSession == nil {
		slackApp, err = mybot.NewFileOAuthApp(c.String("slack-app"))
	} else {
		col := dbSession.DB(dbName).C("slack-app-auth")
		slackApp, err = mybot.NewDBOAuthApp(col)
	}
	if err != nil {
		return mybot.WithStack(err)
	}
	if slackCk != "" && slackCs != "" {
		slackApp.SetCreds(slackCk, slackCs)
		err := slackApp.Save()
		if err != nil {
			return mybot.WithStack(err)
		}
	}

	userIDs, err := getUserIDs(c, dbSession, dbName)
	if err != nil {
		return mybot.WithStack(err)
	}
	for _, userID := range userIDs {
		err := initForUser(c, dbSession, dbName, userID)
		fmt.Printf("Initialize for user %s\n", userID)
		if err != nil {
			return mybot.WithStack(err)
		}
	}
	return nil
}

func initForUser(c *cli.Context, session *mgo.Session, dbName, userID string) error {
	data, err := newUserSpecificData(c, session, userID)
	if err != nil {
		return mybot.WithStack(err)
	}
	userSpecificDataMap[userID] = data
	return nil
}

func getUserIDs(c *cli.Context, session *mgo.Session, dbName string) ([]string, error) {
	if session == nil {
		dir, err := argValueWithMkdir(c, "twitter-auth")
		if err != nil {
			return nil, mybot.WithStack(err)
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
			return nil, mybot.WithStack(err)
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

func httpServer(c *cli.Context) {
	err := startServer(c.String("host"), c.String("port"), c.String("cert"), c.String("key"))
	exitIfError(err)
}

func exitIfError(err error) {
	if err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func argValueWithMkdir(c *cli.Context, key string) (string, error) {
	dir := c.String(key)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return "", mybot.WithStack(err)
	}
	return dir, nil
}
