package database

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DbOptions struct {
	Ctx      context.Context
	Name     string
	Dsn      string
	Protocol string
	Hosts    string
	Username string
	Password string
	Options  string
}

type Db struct {
	client   *mongo.Client
	database *mongo.Database
	dsn      string
	hosts    string
	name     string
	protocol string
	username string
	password string
	options  string
	ctx      context.Context
}

func (db *Db) Connect() error {

	var err error

	slog.Info("connecting to database...")

	db.client, err = mongo.Connect(db.ctx, options.Client().ApplyURI(db.dsn))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to connect to database, Err: %s", err.Error()))
		return err
	}

	if err := db.client.Ping(db.ctx, readpref.Primary()); err != nil {
		slog.Error(fmt.Sprintf("failed to communicate with database, Err: %s", err.Error()))
		return err
	}

	db.database = db.client.Database(db.name)
	slog.Info("database connected successfully")
	return nil
}

func (db *Db) Disconnect(force bool) error {

	if err := db.client.Disconnect(db.ctx); err != nil || !force {
		return err
	}

	db.database = nil
	db.client = nil

	return nil
}

func (db *Db) GetContext() context.Context {
	return db.ctx
}

func (db *Db) GetClient() *mongo.Client {
	return db.client
}

func (db *Db) GetDatabase() *mongo.Database {
	return db.database
}
