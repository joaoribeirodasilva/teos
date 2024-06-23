package main

import (
	"os"

	"github.com/joaoribeirodasilva/teos/apps/routes"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/payload"
	"github.com/joaoribeirodasilva/teos/common/server"
	controllerServices "github.com/joaoribeirodasilva/teos/consumers/logs/services"
)

const (
	SERVICE_NAME = "clogs"
	VERSION      = "0.0.1"
)

func main() {
	logger.SetApplication(SERVICE_NAME)

	conf := configuration.New(SERVICE_NAME)
	if err := conf.GetEnv(); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to read environment", err, nil)
		os.Exit(1)
	}

	logger.SetLogLevel(logger.LogLevel(conf.GetLog().Level))

	services := payload.Payload{}
	services.SetConfig(conf)
	if err := services.SetDatabase(nil); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to create database configuration", err, nil)
		os.Exit(1)
	}

	if err := services.Start(); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to start enabled services", err, nil)
		os.Exit(1)
	}

	if err := conf.GetConfiguration(services.Services.Db); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to red configuration from database", err, nil)
		os.Exit(1)
	}

	if err := services.SetLogsDb(nil); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to start log database", err, nil)
		os.Exit(1)
	}

	logger.SetDatabase(services.Services.LogsDb)

	if err := services.SetPermissionsDb(nil); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to start permissions database", err, nil)
		os.Exit(1)
	}

	if err := services.SetSessionsDb(nil); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to start sessions database", err, nil)
		os.Exit(1)
	}

	if err := services.SetHistoryDb(nil); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to start history database", err, nil)
		os.Exit(1)
	}

	svc := server.New(conf.GetServices().Http.Addr, conf.GetServices().Http.Port)

	if err := services.SetHttp(svc.Service); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to create Http service", err, nil)
		os.Exit(1)
	}

	if err := services.Start(); err != nil {
		logger.Error(logger.LogStatusInternalServerError, nil, "failed to start services", err, nil)
		os.Exit(1)
	}

	loader := controllerServices.NewRoutesService(&services)
	if err := loader.Refresh(); err != nil {
		os.Exit(1)
	}

	router := server.NewRouter(&services)
	routes.RegisterRoutes(router)
	if err := svc.Listen(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)

}
