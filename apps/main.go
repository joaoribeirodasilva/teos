package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joaoribeirodasilva/teos/apps/routes"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/environment"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/common/structures"
)

const (
	SERVICE_NAME = "apps"
	VERSION      = "0.0.1"
)

func main() {

	logger.SetApplication(SERVICE_NAME)
	logger.SetCollectionName("log_logs")

	env := environment.New()
	if err := env.Read(); err != nil {
		os.Exit(1)
	}

	logger.SetLogLevel(logger.LogLevel(env.Application.LogLevel))

	dbOpts := &database.DbOptions{
		Ctx:      context.TODO(),
		Dsn:      env.Database.Dsn,
		Host:     env.Database.Host,
		Name:     env.Database.Name,
		Username: env.Database.Username,
		Password: env.Database.Password,
		Options:  env.Database.Options,
	}

	db := database.New(dbOpts)

	if err := db.Connect(); err != nil {
		os.Exit(1)
	}

	logger.SetDatabase(db)

	confOpts := &configuration.ConfigurationOptions{
		Application: SERVICE_NAME,
		Db:          db,
		Environment: env.Application.EnvironmentName,
	}

	configuration := configuration.New(confOpts)

	if err := configuration.GetConfiguration(); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "error reading configuration from database", err, nil)
		os.Exit(1)
	}

	// Redis databases
	redisConf, err := configuration.GetRedisConf("DB_PERMISSIONS")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "bad DB_PERMISSIONS configuration", err, nil)
		os.Exit(1)
	}

	permissionsDB := redisdb.New("PERMISSIONS", fmt.Sprintf("%s:%d", redisConf.Addr, redisConf.Port), redisConf.Db, redisConf.Username, redisConf.Password)
	if appErr := permissionsDB.Connect(); appErr != nil {
		os.Exit(1)
	}

	redisConf, err = configuration.GetRedisConf("DB_SESSIONS")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "bad DB_SESSIONS configuration", err, nil)
		os.Exit(1)
	}

	sessionsDB := redisdb.New("SESSIONS", fmt.Sprintf("%s:%d", redisConf.Addr, redisConf.Port), redisConf.Db, redisConf.Username, redisConf.Password)
	if appErr := sessionsDB.Connect(); appErr != nil {
		os.Exit(1)
	}

	redisConf, err = configuration.GetRedisConf("DB_HISTORY")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "bad DB_HISTORY configuration", err, nil)
		os.Exit(1)
	}

	historyDB := redisdb.New("HISTORY", fmt.Sprintf("%s:%d", redisConf.Addr, redisConf.Port), redisConf.Db, redisConf.Username, redisConf.Password)
	if appErr := historyDB.Connect(); appErr != nil {
		os.Exit(1)
	}

	redisConf, err = configuration.GetRedisConf("DB_LOGS")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "bad DB_LOGS configuration", err, nil)
		os.Exit(1)
	}

	logsDB := redisdb.New("LOGS", fmt.Sprintf("%s:%d", redisConf.Addr, redisConf.Port), redisConf.Db, redisConf.Username, redisConf.Password)
	if appErr := logsDB.Connect(); appErr != nil {
		os.Exit(1)
	}

	// Http server
	httpAddr, err := configuration.GetString("HTTP_BIND_ADDR")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "invalid http listen address", err, nil)
		os.Exit(1)
	}

	httpPort, err := configuration.GetInt("HTTP_BIND_PORT")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "invalid http listen address", err, nil)
		os.Exit(1)
	}

	svc := server.New(db, httpAddr, int(httpPort))

	services := &structures.Services{
		Gin:           svc.Service,
		SessionsDB:    sessionsDB,
		PermissionsDB: permissionsDB,
		HistoryDB:     historyDB,
		LogsDB:        logsDB,
		Env:           env,
		Db:            db,
		Configuration: configuration,
	}

	router := server.NewRouter(services)
	routes.RegisterRoutes(router)
	if err := svc.Listen(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
