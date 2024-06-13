package main

import (
	"log/slog"
	"os"

	"github.com/joaoribeirodasilva/teos/auth/routes"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/info"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/common/service_log"
)

const (
	SERVICE_NAME = "auth"
	VERSION      = "0.0.1"
)

type RTest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {

	info.Print(SERVICE_NAME, VERSION)

	service_log.ApplicationName = SERVICE_NAME
	service_log.IsDatabase = false
	service_log.IsStdout = true

	conf := conf.New(SERVICE_NAME)
	if !conf.Read() {
		os.Exit(1)
	}

	db := database.New(conf)
	if err := db.Connect(); err != nil {
		os.Exit(1)
	}

	serviceConfiguration := configuration.New(db, conf)
	if err := serviceConfiguration.GetAppId(); err != nil {
		os.Exit(1)
	}
	if err := serviceConfiguration.Read(); err != nil {
		os.Exit(1)
	}

	sessionsDB := redisdb.New("Sessions Database", serviceConfiguration.DbSessions.Addresses, serviceConfiguration.DbSessions.Db, serviceConfiguration.DbSessions.Username, serviceConfiguration.DbSessions.Password)
	if appErr := sessionsDB.Connect(); appErr != nil {
		os.Exit(1)
	}

	permissionsDB := redisdb.New("Permissions Database", serviceConfiguration.DbPermissions.Addresses, serviceConfiguration.DbPermissions.Db, serviceConfiguration.DbPermissions.Username, serviceConfiguration.DbPermissions.Password)
	if appErr := permissionsDB.Connect(); appErr != nil {
		os.Exit(1)
	}

	if err := service_log.InitServiceLog(serviceConfiguration.DbLog.Addresses, serviceConfiguration.DbLog.Db, serviceConfiguration.DbLog.Username, serviceConfiguration.DbLog.Password); err != nil {
		os.Exit(1)
	}

	service_log.IsDatabase = true
	service_log.IsStdout = true

	tempPort := serviceConfiguration.GetKey("NET_PORT")
	if tempPort == nil || tempPort.Int == nil {
		slog.Error("invalid network port to listen to")
		os.Exit(1)
	}
	conf.Service.BindPort = *tempPort.Int

	svc := server.New(db, conf)
	router := server.NewRouter(svc.Service, conf, db, serviceConfiguration, sessionsDB, permissionsDB)
	routes.RegisterRoutes(router)
	if err := svc.Listen(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
