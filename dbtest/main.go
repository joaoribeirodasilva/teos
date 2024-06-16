package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joaoribeirodasilva/teos/dbtest/database"
	"github.com/joaoribeirodasilva/teos/dbtest/environment"
	"github.com/joaoribeirodasilva/teos/dbtest/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	env := environment.New()
	if err := env.Read(); err != nil {
		os.Exit(1)
	}

	db := database.New(
		context.TODO(),
		env.Database.Name,
		env.Database.Dsn,
		env.Database.Protocol,
		env.Database.Hosts,
		env.Database.Username,
		env.Database.Password,
		env.Database.Options,
	)

	if err := db.Connect(); err != nil {
		os.Exit(1)
	}

	t_model := &models.TestModel{}

	t_coll_options := &models.BaseCollectionOptions{
		Ctx:        db.GetContext(),
		UseMetas:   true,
		UseUserID:  true,
		UseDates:   true,
		SoftDelete: true,
		Debug:      true,
	}

	t_coll := models.NewBaseCollection(t_model.GetCollectionName(), db, t_coll_options)
	userId := primitive.NewObjectID()
	t_coll.SetUserID(&userId)
	t_model.BaseModel.ID = primitive.NewObjectID()
	t_model.Name = "Joao"
	t_model.Age = 53

	if err := t_coll.Create(nil, t_model); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	/* 	meta := &models.MetaModel{}
	   	if err := t_coll.Create(nil, meta); err != nil {
	   		slog.Error(err.Error())
	   		os.Exit(1)
	   	} */
	os.Exit(0)
}
