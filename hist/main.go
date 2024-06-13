package main

import (
	"log/slog"
	"os"

	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/info"
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/hist/routes"
)

const (
	SERVICE_NAME = "hists"
	VERSION      = "0.0.1"
)

func main() {

	info.Print(SERVICE_NAME, VERSION)

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

	tempPort := serviceConfiguration.GetKey("NET_PORT")
	if tempPort == nil || tempPort.Int == nil {
		slog.Error("invalid network port to listen to")
		os.Exit(1)
	}
	conf.Service.BindPort = *tempPort.Int

	svc := server.New(db, conf)
	router := server.NewRouter(svc.Service, conf, db, serviceConfiguration)
	routes.RegisterRoutes(router)
	if err := svc.Listen(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
