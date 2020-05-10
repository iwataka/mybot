package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/iwataka/mybot/core"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/oauth"
	"github.com/iwataka/mybot/runner"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/mybot/worker"
	"github.com/kidstuff/mongostore"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"gopkg.in/mgo.v2"
)

//go:generate mockgen -source=models/slack.go -destination=mocks/slack.go -package=mocks
//go:generate mockgen -source=models/twitter.go -destination=mocks/twitter.go -package=mocks
//go:generate mockgen -source=models/auth.go -destination=mocks/auth.go -package=mocks
//go:generate mockgen -source=core/vision.go -destination=mocks/vision.go -package=mocks
//go:generate mockgen -source=core/language.go -destination=mocks/language.go -package=mocks
//go:generate mockgen -source=utils/utils.go -destination=mocks/utils.go -package=mocks
//go:generate mockgen -source=runner/batch.go -destination=mocks/batch.go -package=mocks
//go:generate mockgen -source=models/worker.go -destination=mocks/worker.go -package=mocks
//go:generate mockgen -source=models/cli.go -destination=mocks/cli.go -package=mocks
//go:generate mockgen -source=models/mgo.go -destination=mocks/mgo.go -package=mocks
//go:generate mockgen -source=data/cache.go -destination=mocks/cache.go -package=mocks
//TODO: When mockgen Config interface, cyclic dependencies happen.

var (
	userSpecificDataMap = make(map[string]*userSpecificData)
	logger              = log.New(os.Stdout, "", log.LstdFlags)
	errLogger           = log.New(os.Stderr, "", log.LstdFlags)

	// Global-scope data
	twitterApp    oauth.OAuthApp
	slackApp      oauth.OAuthApp
	visionAPI     core.VisionMatcher
	languageAPI   core.LanguageMatcher
	cliContext    *cli.Context
	dbSession     *mgo.Session
	serverSession sessions.Store
)

const (
	configFlagName                = "config"
	cacheFlagName                 = "cache"
	gcloudFlagName                = "gcloud"
	twitterFlagName               = "twitter-auth"
	slackFlagName                 = "slack-auth"
	certFlagName                  = "cert"
	keyFlagName                   = "key"
	hostFlagName                  = "host"
	portFlagName                  = "port"
	dbAddrFlagName                = "db-addr"
	dbUserFlagName                = "db-user"
	dbPassFlagName                = "db-passwd"
	dbNameFlagName                = "db-name"
	twitterConsumerKeyFlagName    = "twitter-consumer-key"
	twitterConsumerSecretFlagName = "twitter-consumer-secret"
	twitterConsumerFileFlagName   = "twitter-app"
	slackClientIDFlagName         = "slack-client-id"
	slackClientSecretFlagName     = "slack-client-secret"
	slackClientFileFlagName       = "slack-app"
	sessionDomainFlagName         = "session-domain"
	apiFlagName                   = "api"
	workerBufSizeFlagName         = "worker-buffer-size"
	workerRestartDurationFlagName = "worker-restart-duration"
	workerRestartLimitFlagName    = "worker-restart-limit"
)

const (
	twitterDMRoutineKey = iota
	twitterUserRoutineKey
	slackRoutineKey
	twitterPeriodicRoutineKey
)

const (
	defaultWorkerBufSize = 10
)

func workerKeys() []int {
	return []int{twitterDMRoutineKey, twitterUserRoutineKey, slackRoutineKey, twitterPeriodicRoutineKey}
}

func main() {
	home, err := homedir.Dir()
	utils.ExitIfError(err)

	log.SetFlags(0)

	configDir := filepath.Join(home, ".config", "mybot")
	cacheDir := filepath.Join(home, ".cache", "mybot")

	configFlag := cli.StringFlag{
		Name:   configFlagName,
		Value:  filepath.Join(configDir, "config"),
		Usage:  "Config directory location",
		EnvVar: "MYBOT_CONFIG_PATH",
	}

	cacheFlag := cli.StringFlag{
		Name:   cacheFlagName,
		Value:  filepath.Join(cacheDir, "cache"),
		Usage:  "Cache directory location",
		EnvVar: "MYBOT_CACHE_PATH",
	}

	gcloudFlag := cli.StringFlag{
		Name:   gcloudFlagName,
		Usage:  "Credential file for Google Cloud Platform",
		EnvVar: "MYBOT_GCLOUD_CREDENTIAL",
	}

	twitterFlag := cli.StringFlag{
		Name:   twitterFlagName,
		Value:  filepath.Join(configDir, "twitter_auth"),
		Usage:  "Twitter credential directory",
		EnvVar: "MYBOT_TWITTER_CREDENTIAL",
	}

	slackFlag := cli.StringFlag{
		Name:   slackFlagName,
		Value:  filepath.Join(configDir, "slack_auth"),
		Usage:  "Slack credential directory",
		EnvVar: "MYBOT_SLACK_CREDENTIAL",
	}

	certFlag := cli.StringFlag{
		Name:   certFlagName,
		Value:  filepath.Join(configDir, "mybot.crt"),
		Usage:  "Certification file for server",
		EnvVar: "MYBOT_SSL_CERT",
	}

	keyFlag := cli.StringFlag{
		Name:   keyFlagName,
		Value:  filepath.Join(configDir, "mybot.key"),
		Usage:  "Key file for server",
		EnvVar: "MYBOT_SSL_KEY",
	}

	hostFlag := cli.StringFlag{
		Name:   strings.Join([]string{hostFlagName, "H"}, ","),
		Value:  "localhost",
		Usage:  "Host this server listen on",
		EnvVar: "MYBOT_HOST",
	}

	portFlag := cli.StringFlag{
		Name:   strings.Join([]string{portFlagName, "P"}, ","),
		Value:  "8080",
		Usage:  "Port this server listen on",
		EnvVar: "MYBOT_PORT",
	}

	dbAddrFlag := cli.StringFlag{
		Name:   dbAddrFlagName,
		Value:  "",
		Usage:  "DB address",
		EnvVar: "MYBOT_DB_ADDRESS",
	}

	dbUserFlag := cli.StringFlag{
		Name:   dbUserFlagName,
		Value:  "",
		Usage:  "DB user",
		EnvVar: "MYBOT_DB_USER",
	}

	dbPassFlag := cli.StringFlag{
		Name:   dbPassFlagName,
		Value:  "",
		Usage:  "DB password",
		EnvVar: "MYBOT_DB_PASSWD",
	}

	dbNameFlag := cli.StringFlag{
		Name:   dbNameFlagName,
		Value:  "",
		Usage:  "Target DB name",
		EnvVar: "MYBOT_DB_NAME",
	}

	twitterConsumerKeyFlag := cli.StringFlag{
		Name:   twitterConsumerKeyFlagName,
		Value:  "",
		Usage:  "Twitter consumer key",
		EnvVar: "MYBOT_TWITTER_CONSUMER_KEY",
	}

	twitterConsumerSecretFlag := cli.StringFlag{
		Name:   twitterConsumerSecretFlagName,
		Value:  "",
		Usage:  "Twitter consumer secret",
		EnvVar: "MYBOT_TWITTER_CONSUMER_SECRET",
	}

	twitterConsumerFileFlag := cli.StringFlag{
		Name:   twitterConsumerFileFlagName,
		Value:  filepath.Join(configDir, "twitter_app.toml"),
		Usage:  "Twitter application directory",
		EnvVar: "MYBOT_TWITTER_APP",
	}

	slackClientIDFlag := cli.StringFlag{
		Name:   slackClientIDFlagName,
		Value:  "",
		Usage:  "Slack client ID",
		EnvVar: "MYBOT_SLACK_APP",
	}

	slackClientSecretFlag := cli.StringFlag{
		Name:   slackClientSecretFlagName,
		Value:  "",
		Usage:  "Slack client secret",
		EnvVar: "MYBOT_SLACK_CLIENT_SECRET",
	}

	slackClientFileFlag := cli.StringFlag{
		Name:   slackClientFileFlagName,
		Value:  filepath.Join(configDir, "slack_app.toml"),
		Usage:  "Slack application directory",
		EnvVar: "MYBOT_SLACK_APP",
	}

	sessionDomainFlag := cli.StringFlag{
		Name:   sessionDomainFlagName,
		Value:  "",
		Usage:  "Session domain",
		EnvVar: "MYBOT_SESSION_DOMAIN",
	}

	apiFlag := cli.BoolFlag{
		Name:  apiFlagName,
		Usage: "Use API to validate configuration",
	}

	workerBufSizeFlag := cli.IntFlag{
		Name:   workerBufSizeFlagName,
		Value:  defaultWorkerBufSize,
		Usage:  "Worker channel buffer size",
		EnvVar: "MYBOT_WORKER_BUFFER_SIZE",
	}

	workerRestartDurationFlag := cli.DurationFlag{
		Name:   workerRestartDurationFlagName,
		Value:  15 * time.Minute,
		Usage:  "Worker restart duration",
		EnvVar: "MYBOT_WORKER_RESTART_DURATION",
	}

	workerRestartLimitFlag := cli.IntFlag{
		Name:   workerRestartLimitFlagName,
		Value:  5,
		Usage:  "Worker restart limit",
		EnvVar: "MYBOT_WORKER_RESTART_LIMIT",
	}

	commonFlags := []cli.Flag{
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

	serveFlags := []cli.Flag{
		gcloudFlag,
		certFlag,
		keyFlag,
		hostFlag,
		portFlag,
		sessionDomainFlag,
		workerBufSizeFlag,
		workerRestartDurationFlag,
		workerRestartLimitFlag,
	}
	// All `run` flags should be `serve` flag
	serveFlags = append(serveFlags, commonFlags...)

	validateFlags := []cli.Flag{apiFlag}
	// All `run` flags should be `validate` flag
	validateFlags = append(validateFlags, commonFlags...)

	app := cli.NewApp()
	app.Name = "mybot"
	app.Version = "0.1"
	app.Usage = "Automatically collect and broadcast information based on your configuration"
	app.Author = "iwataka"

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

	app.Commands = []cli.Command{serveCmd, validateCmd}
	err = app.Run(os.Args)
	utils.ExitIfError(err)
}

type userSpecificData struct {
	config      core.Config
	cache       data.Cache
	twitterAPI  *core.TwitterAPI
	twitterAuth oauth.OAuthCreds
	slackAPI    *core.SlackAPI
	slackAuth   oauth.OAuthCreds
	workerMgrs  map[int]*worker.WorkerManager
}

func (d *userSpecificData) statuses() map[int]bool {
	s := map[int]bool{}
	for _, key := range workerKeys() {
		wm := d.workerMgrs[key]
		if wm == nil {
			s[key] = false
		} else {
			s[key] = wm.Status() == worker.StatusStarted
		}
	}
	return s
}

func (d *userSpecificData) delete() error {
	for k, wm := range d.workerMgrs {
		wm.Close()
		delete(d.workerMgrs, k)
	}
	for _, del := range []utils.Deletable{d.config, d.cache, d.twitterAuth, d.slackAuth} {
		err := del.Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *userSpecificData) restart() {
	for _, ch := range d.workerMgrs {
		ch.Send(worker.RestartSignal)
	}
}

func newUserSpecificData(c models.Context, session *mgo.Session, userID string) (*userSpecificData, error) {
	var err error
	userData := &userSpecificData{}
	userData.workerMgrs = map[int]*worker.WorkerManager{}
	dbName := c.String(dbNameFlagName)

	if session == nil {
		userData.cache, err = newFileCache(c, userID)
	} else {
		col := models.NewMgoCollection(session.DB(dbName).C("cache"))
		userData.cache, err = data.NewDBCache(col, userID)
	}
	if err != nil {
		return nil, utils.WithStack(err)
	}

	if session == nil {
		userData.config, err = newFileConfig(c, userID)
	} else {
		col := models.NewMgoCollection(session.DB(dbName).C("config"))
		userData.config, err = core.NewDBConfig(col, userID)
	}
	if err != nil {
		return nil, utils.WithStack(err)
	}

	if session == nil {
		userData.twitterAuth, err = newFileOAuthCreds(c, twitterFlagName, userID)
	} else {
		col := models.NewMgoCollection(session.DB(dbName).C("twitter-user-auth"))
		userData.twitterAuth, err = oauth.NewDBOAuthCreds(col, userID)
	}
	if err != nil {
		return nil, utils.WithStack(err)
	}
	// saving twitter auth means user-register
	err = userData.twitterAuth.Save()
	if err != nil {
		return nil, utils.WithStack(err)
	}

	if session == nil {
		userData.slackAuth, err = newFileOAuthCreds(c, slackFlagName, userID)
	} else {
		col := models.NewMgoCollection(session.DB(dbName).C("slack-user-auth"))
		userData.slackAuth, err = oauth.NewDBOAuthCreds(col, userID)
	}
	if err != nil {
		return nil, utils.WithStack(err)
	}

	userData.twitterAPI = core.NewTwitterAPIWithAuth(userData.twitterAuth, userData.config, userData.cache)
	slackID, _ := userData.slackAuth.GetCreds()
	userData.slackAPI = core.NewSlackAPIWithAuth(slackID, userData.config, userData.cache)

	return userData, nil
}

func newFileCache(c models.Context, userID string) (data.Cache, error) {
	dir, err := argValueWithMkdir(c, cacheFlagName)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
	cache, err := data.NewFileCache(file)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	return cache, nil
}

func newFileConfig(c models.Context, userID string) (core.Config, error) {
	dir, err := argValueWithMkdir(c, configFlagName)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
	config, err := core.NewFileConfig(file)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func newFileOAuthCreds(c models.Context, flagName, userID string) (oauth.OAuthCreds, error) {
	dir, err := argValueWithMkdir(c, flagName)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	file := filepath.Join(dir, fmt.Sprintf("%s.toml", userID))
	creds, err := oauth.NewFileOAuthCreds(file)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	return creds, nil
}

func startUserSpecificData(c models.Context, data *userSpecificData, userID string) error {
	if len(data.workerMgrs) > 0 {
		return fmt.Errorf("%s's workers already started", userID)
	}

	restartDuration := c.Duration(workerRestartDurationFlagName)
	restartLimit := c.Int(workerRestartLimitFlagName)
	restarter := worker.NewStrategicRestarter(restartDuration, restartLimit, false)
	bufSize := c.Int(workerBufSizeFlagName)

	twitterDMWorker := newTwitterDMWorker(data.twitterAPI, userID)
	data.workerMgrs[twitterDMRoutineKey] = activateWorkerAndStart(
		twitterDMWorker,
		workerMessageLogger{twitterDMWorker.Name(), logger, errLogger},
		bufSize,
		restarter,
	)

	twitterUserWorker := newTwitterUserWorker(data.twitterAPI, data.slackAPI, visionAPI, languageAPI, data.cache, userID)
	data.workerMgrs[twitterUserRoutineKey] = activateWorkerAndStart(
		twitterUserWorker,
		workerMessageLogger{twitterUserWorker.Name(), logger, errLogger},
		bufSize,
		restarter,
	)

	r := runner.NewBatchRunnerUsedWithStream(data.twitterAPI, data.slackAPI, visionAPI, languageAPI, data.config)
	twitterPeriodicWorker := newTwitterPeriodicWorker(r, data.cache, data.config, userID)
	data.workerMgrs[twitterPeriodicRoutineKey] = activateWorkerAndStart(
		twitterPeriodicWorker,
		workerMessageLogger{twitterPeriodicWorker.Name(), logger, errLogger},
		bufSize,
		restarter,
	)

	slackWorker := newSlackWorker(data.slackAPI, data.twitterAPI, visionAPI, languageAPI, userID)
	data.workerMgrs[slackRoutineKey] = activateWorkerAndStart(
		slackWorker,
		workerMessageLogger{slackWorker.Name(), logger, errLogger},
		bufSize,
		restarter,
	)

	return nil
}

func serve(c *cli.Context) error {
	go httpServer(c)
	var ch chan interface{}
	<-ch
	return nil
}

func validate(c *cli.Context) {
	for _, data := range userSpecificDataMap {
		err := data.config.Validate()
		utils.ExitIfError(err)
		if c.Bool(apiFlagName) {
			err := data.config.ValidateWithAPI(data.twitterAPI.BaseAPI())
			utils.ExitIfError(err)
		}
	}
}

func beforeServing(c *cli.Context) error {
	// ignore error at first initialization
	visionAPI, _ = core.NewVisionMatcher(c.String(gcloudFlagName))
	languageAPI, _ = core.NewLanguageMatcher(c.String(gcloudFlagName))

	err := beforeValidating(c)
	utils.ExitIfError(err)
	for userID, userData := range userSpecificDataMap {
		err := startUserSpecificData(c, userData, userID)
		utils.ExitIfError(err)
	}

	return nil
}

func beforeValidating(c *cli.Context) error {
	cliContext = c
	dbAddress := c.String(dbAddrFlagName)
	dbUser := c.String(dbUserFlagName)
	dbPasswd := c.String(dbPassFlagName)
	dbName := c.String(dbNameFlagName)
	var err error

	if dbAddress != "" && dbName != "" {
		info := &mgo.DialInfo{}
		info.Addrs = []string{dbAddress}
		info.Username = dbUser
		info.Password = dbPasswd
		info.Database = dbName
		dbSession, err = mgo.DialWithInfo(info)
	}
	if err != nil {
		return utils.WithStack(err)
	}

	initSession(c, dbName)

	err = initTwitterApp(c, dbName)
	if err != nil {
		return err
	}

	err = initSlackApp(c, dbName)
	if err != nil {
		return err
	}

	userIDs, err := getUserIDs(c, dbSession, dbName)
	if err != nil {
		return utils.WithStack(err)
	}
	for _, userID := range userIDs {
		err := initForUser(c, dbSession, userID)
		logger.Printf("Initialize for user %s\n", userID)
		if err != nil {
			return utils.WithStack(err)
		}
	}
	return nil
}

func initSession(c *cli.Context, dbName string) {
	sessionDomain := c.String(sessionDomainFlagName)
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
}

func initTwitterApp(c models.Context, dbName string) error {
	var err error

	twitterCk := c.String(twitterConsumerKeyFlagName)
	twitterCs := c.String(twitterConsumerSecretFlagName)
	if dbSession == nil {
		twitterApp, err = oauth.NewFileTwitterOAuthApp(c.String(twitterConsumerFileFlagName))
	} else {
		col := models.NewMgoCollection(dbSession.DB(dbName).C("twitter-app-auth"))
		twitterApp, err = oauth.NewDBTwitterOAuthApp(col)
	}
	if err != nil {
		return utils.WithStack(err)
	}

	if twitterCk != "" && twitterCs != "" {
		twitterApp.SetCreds(twitterCk, twitterCs)
		err = twitterApp.Save()
		if err != nil {
			return utils.WithStack(err)
		}
	}

	return nil
}

func initSlackApp(c models.Context, dbName string) error {
	var err error

	slackCk := c.String(slackClientIDFlagName)
	slackCs := c.String(slackClientSecretFlagName)
	if dbSession == nil {
		slackApp, err = oauth.NewFileOAuthApp(c.String(slackClientFileFlagName))
	} else {
		col := models.NewMgoCollection(dbSession.DB(dbName).C("slack-app-auth"))
		slackApp, err = oauth.NewDBOAuthApp(col)
	}
	if err != nil {
		return utils.WithStack(err)
	}

	if slackCk != "" && slackCs != "" {
		slackApp.SetCreds(slackCk, slackCs)
		err = slackApp.Save()
		if err != nil {
			return utils.WithStack(err)
		}
	}

	return nil
}

func initForUser(c models.Context, session *mgo.Session, userID string) error {
	data, err := newUserSpecificData(c, session, userID)
	if err != nil {
		return utils.WithStack(err)
	}
	userSpecificDataMap[userID] = data
	return nil
}

// getUserIDs returns all user IDs by checking Twitter user athentication data.
func getUserIDs(c models.Context, session *mgo.Session, dbName string) ([]string, error) {
	if session == nil {
		dir, err := argValueWithMkdir(c, twitterFlagName)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		files, err := filepath.Glob(filepath.Join(dir, "*.toml"))
		if err != nil {
			return nil, utils.WithStack(err)
		}
		userIDs := []string{}
		for _, file := range files {
			base := filepath.Base(file)
			ext := filepath.Ext(file)
			userIDs = append(userIDs, base[0:len(base)-len(ext)])
		}
		return userIDs, nil
	}

	col := session.DB(dbName).C("twitter-user-auth")
	auths := []map[string]interface{}{}
	err := col.Find(nil).All(&auths)
	if err != nil {
		return nil, utils.WithStack(err)
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

func httpServer(c models.Context) {
	err := startServer(c.String(hostFlagName), c.String(portFlagName), c.String(certFlagName), c.String(keyFlagName))
	utils.ExitIfError(err)
}

func argValueWithMkdir(c models.Context, key string) (string, error) {
	dir := c.String(key)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return "", utils.WithStack(err)
	}
	return dir, nil
}
