package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
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
	"github.com/markbates/goth/gothic"
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
	logger              = log.New(os.Stdout, "[mybot] ", log.LstdFlags)
	errLogger           = log.New(os.Stderr, "[mybot] ", log.LstdFlags)

	// Global-scope data
	twitterApp    oauth.OAuthApp
	slackApp      oauth.OAuthApp
	visionAPI     core.VisionMatcher
	languageAPI   core.LanguageMatcher
	cliContext    *cli.Context
	database      *mgo.Database
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
	defaultWorkerBufSize  = 10
	twitterCollectionName = "twitter-app-auth"
	slackCollectionName   = "slack-app-auth"
)

func workerKeys() []int {
	return []int{twitterDMRoutineKey, twitterUserRoutineKey, slackRoutineKey, twitterPeriodicRoutineKey}
}

func main() {
	log.SetFlags(0)

	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(home, ".config", "mybot")
	cacheDir := filepath.Join(home, ".cache", "mybot")

	commonFlags := getCommonFlags(configDir, cacheDir)
	serveFlags := append(commonFlags, getServeSpecificFlags(configDir)...)
	validateFlags := append(commonFlags, getValidateSpecificFlags()...)

	app := cli.NewApp()
	app.Name = "mybot"
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
	if err != nil {
		panic(err)
	}
}

func getCommonFlags(configDir, cacheDir string) []cli.Flag {
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

	return []cli.Flag{
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
}

func getServeSpecificFlags(configDir string) []cli.Flag {
	gcloudFlag := cli.StringFlag{
		Name:   gcloudFlagName,
		Usage:  "Credential file for Google Cloud Platform",
		EnvVar: "MYBOT_GCLOUD_CREDENTIAL",
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

	sessionDomainFlag := cli.StringFlag{
		Name:   sessionDomainFlagName,
		Value:  "",
		Usage:  "Session domain",
		EnvVar: "MYBOT_SESSION_DOMAIN",
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

	return []cli.Flag{
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
}

func getValidateSpecificFlags() []cli.Flag {
	apiFlag := cli.BoolFlag{
		Name:  apiFlagName,
		Usage: "Use API to validate configuration",
	}

	return []cli.Flag{apiFlag}
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
		if del != nil {
			err := del.Delete()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *userSpecificData) restart() {
	for _, ch := range d.workerMgrs {
		ch.Send(worker.RestartSignal)
	}
}

func newUserSpecificData(c models.Context, database *mgo.Database, userID string) (*userSpecificData, error) {
	var err error
	userData := &userSpecificData{}
	userData.workerMgrs = map[int]*worker.WorkerManager{}

	if database == nil {
		userData.cache, err = newFileCache(c, userID)
	} else {
		col := models.NewMgoCollection(database.C("cache"))
		userData.cache, err = data.NewDBCache(col, userID)
	}
	if err != nil {
		return nil, utils.WithStack(err)
	}

	if database == nil {
		userData.config, err = newFileConfig(c, userID)
	} else {
		col := models.NewMgoCollection(database.C("config"))
		userData.config, err = core.NewDBConfig(col, userID)
	}
	if err != nil {
		return nil, utils.WithStack(err)
	}

	if database == nil {
		userData.twitterAuth, err = newFileOAuthCreds(c, twitterFlagName, userID)
	} else {
		col := models.NewMgoCollection(database.C("twitter-user-auth"))
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

	if database == nil {
		userData.slackAuth, err = newFileOAuthCreds(c, slackFlagName, userID)
	} else {
		col := models.NewMgoCollection(database.C("slack-user-auth"))
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
	gothic.Store = serverSession

	host, port := c.String(hostFlagName), c.String(portFlagName)
	addr := fmt.Sprintf("%s:%s", host, port)
	cert, key := c.String(certFlagName), c.String(keyFlagName)
	_, certErr := os.Stat(cert)
	_, keyErr := os.Stat(key)

	r := setupRouter()
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Graceful shutdown
	// https://github.com/gin-gonic/gin#graceful-shutdown-or-restart
	go func() {
		var err error
		if certErr == nil && keyErr == nil {
			err = srv.ListenAndServeTLS(cert, key)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server force to shutdown: %s", err)
	}
	return nil
}

func validate(c *cli.Context) error {
	for _, data := range userSpecificDataMap {
		err := data.config.Validate()
		if err != nil {
			return err
		}
		if c.Bool(apiFlagName) {
			err := data.config.ValidateWithAPI(data.twitterAPI.BaseAPI())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func beforeServing(c *cli.Context) error {
	// ignore error at first initialization
	visionAPI, _ = core.NewVisionMatcher(c.String(gcloudFlagName))
	languageAPI, _ = core.NewLanguageMatcher(c.String(gcloudFlagName))

	err := beforeValidating(c)
	if err != nil {
		return err
	}
	for userID, userData := range userSpecificDataMap {
		err := startUserSpecificData(c, userData, userID)
		if err != nil {
			return err
		}
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
		dbSession, err := mgo.DialWithInfo(info)
		if err != nil {
			return utils.WithStack(err)
		}
		database = dbSession.DB(dbName)
	}

	serverSession = newSessionStore(c, database)

	twitterApp, err = newTwitterApp(c, database, twitterCollectionName)
	if err != nil {
		return err
	}

	slackApp, err = newSlackApp(c, database, slackCollectionName)
	if err != nil {
		return err
	}

	userIDs, err := getUserIDs(c, database)
	if err != nil {
		return utils.WithStack(err)
	}
	for _, userID := range userIDs {
		data, err := newUserSpecificData(c, database, userID)
		if err != nil {
			return utils.WithStack(err)
		}
		userSpecificDataMap[userID] = data
	}
	return nil
}

func newSessionStore(c models.Context, database *mgo.Database) sessions.Store {
	sessionDomain := c.String(sessionDomainFlagName)
	if database == nil {
		sess := sessions.NewCookieStore(
			[]byte("mybot_session_key"),
		)
		if sessionDomain != "" {
			sess.Options.Domain = sessionDomain
		}
		return sess
	}
	sess := mongostore.NewMongoStore(
		database.C("user-session"),
		86400*30,
		true,
		[]byte("mybot_session_key"),
	)
	if sessionDomain != "" {
		sess.Options.Domain = sessionDomain
	}
	return sess
}

func newTwitterApp(c models.Context, database *mgo.Database, colName string) (oauth.OAuthApp, error) {
	ck := c.String(twitterConsumerKeyFlagName)
	cs := c.String(twitterConsumerSecretFlagName)
	fpath := c.String(twitterConsumerFileFlagName)
	col := newMgoCollection(database, colName)
	return newOAuthApp(ck, cs, fpath, col, oauth.NewFileTwitterOAuthApp, oauth.NewDBTwitterOAuthApp)
}

func newSlackApp(c models.Context, database *mgo.Database, colName string) (oauth.OAuthApp, error) {
	ck := c.String(slackClientIDFlagName)
	cs := c.String(slackClientSecretFlagName)
	fpath := c.String(slackClientFileFlagName)
	col := newMgoCollection(database, colName)
	return newOAuthApp(ck, cs, fpath, col, oauth.NewFileOAuthApp, oauth.NewDBOAuthApp)
}

func newMgoCollection(database *mgo.Database, colName string) models.MgoCollection {
	if database != nil {
		return models.NewMgoCollection(database.C(colName))
	}
	return nil
}

func newOAuthApp(
	ck, cs, fpath string,
	col models.MgoCollection,
	newFile func(string) (oauth.OAuthApp, error),
	newDB func(models.MgoCollection) (oauth.OAuthApp, error),
) (oauth.OAuthApp, error) {
	var app oauth.OAuthApp
	var err error
	if col == nil {
		app, err = newFile(fpath)
	} else {
		app, err = newDB(col)
	}
	if err != nil {
		return nil, utils.WithStack(err)
	}

	if len(ck) > 0 && len(cs) > 0 {
		app.SetCreds(ck, cs)
		err = app.Save()
		if err != nil {
			return nil, utils.WithStack(err)
		}
	}

	return app, nil
}

// getUserIDs returns all user IDs by checking Twitter user athentication data.
func getUserIDs(c models.Context, database *mgo.Database) ([]string, error) {
	if database == nil {
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

	col := database.C("twitter-user-auth")
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

func argValueWithMkdir(c models.Context, key string) (string, error) {
	dir := c.String(key)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return "", utils.WithStack(err)
	}
	return dir, nil
}
