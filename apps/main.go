package main

import (
	"context"
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
		Protocol: env.Database.Protocol,
		Hosts:    env.Database.Hosts,
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
		Application:                    SERVICE_NAME,
		Db:                             db,
		AppConfigurationCollectionName: "app_configurations",
		AppCollectionName:              "app_apps",
	}

	configuration := configuration.New(confOpts)

	if err := configuration.GetConfiguration(); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "error reading configuration from database", err, nil)
		os.Exit(1)
	}

	redisAddress, err := configuration.GetString("REDIS_DB_SESSION_ADDRESS")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "getting session database address", err, nil)
		os.Exit(1)
	}
	redisDatabase, err := configuration.GetInt("REDIS_DB_SESSION_DATABASE")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "getting session database number", err, nil)
		os.Exit(1)
	}
	redisPassword, err := configuration.GetString("REDIS_DB_SESSION_PASSWORD")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "getting session database password", err, nil)
		os.Exit(1)
	}
	redisUsername := ""

	sessionsDB := redisdb.New(redisAddress, int(redisDatabase), redisUsername, redisPassword)
	if appErr := sessionsDB.Connect(); appErr != nil {
		os.Exit(1)

	}

	redisAddress, err = configuration.GetString("REDIS_DB_SESSION_ADDRESS")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "getting session database address", err, nil)
		os.Exit(1)
	}
	redisDatabase, err = configuration.GetInt("REDIS_DB_SESSION_DATABASE")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "getting session database number", err, nil)
		os.Exit(1)
	}
	redisPassword, err = configuration.GetString("REDIS_DB_SESSION_PASSWORD")
	if err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "getting session database password", err, nil)
		os.Exit(1)
	}
	redisUsername = ""

	permissionsDB := redisdb.New(redisAddress, int(redisDatabase), redisUsername, redisPassword)
	if appErr := permissionsDB.Connect(); appErr != nil {
		os.Exit(1)
	}

	svc := server.New(db, &env.Application)

	services := &structures.Services{
		Gin:           svc.Service,
		SessionsDB:    sessionsDB,
		PermissionsDB: permissionsDB,
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
