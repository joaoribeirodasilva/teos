package main

import (
	"os"

	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/info"
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/users/routes"
)

const (
	SERVICE_NAME = "users"
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

	svc := server.New(db, conf)
	router := server.NewRouter(svc.Service, conf, db)
	routes.RegisterRoutes(router)
	if err := svc.Listen(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
