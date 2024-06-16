package main

import (
	"context"
	"os"

	"github.com/joaoribeirodasilva/teos/dbtest/configuration"
	"github.com/joaoribeirodasilva/teos/dbtest/database"
	"github.com/joaoribeirodasilva/teos/dbtest/environment"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
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
		Username: env.Database.Username,
		Password: env.Database.Password,
		Options:  env.Database.Options,
	}

	db := database.New("", dbOpts)

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
		logger.Error(logger.LogStatusNotFound, nil, "error reading configuration from database", err, nil)
		os.Exit(1)
	}

	logger.Info("service '%s' started", SERVICE_NAME)

	/* 	t_model := &models.TestModel{}

	   	t_coll_options := &models.BaseCollectionOptions{
	   		Ctx:        db.GetContext(),
	   		UseMetas:   true,
	   		UseUserID:  true,
	   		UseDates:   true,
	   		SoftDelete: true,
	   		Debug:      true,
	   	}

	   	logger.Info("will insert test on the database")

	   	t_coll := models.NewBaseCollection(t_model.GetCollectionName(), db, t_coll_options)
	   	userId := primitive.NewObjectID()
	   	t_coll.SetUserID(&userId)
	   	t_model.BaseModel.ID = primitive.NewObjectID()
	   	t_model.Name = "Joao"
	   	t_model.Age = 53

	   	if err := t_coll.Create(nil, t_model); err != nil {
	   		logger.Error(logger.LogStatusInternalServerError, nil, "failed to insert test in the database", err, nil)
	   		os.Exit(1)
	   	} */

	logger.Info("service '%s' streminated", SERVICE_NAME)

	os.Exit(0)
}
