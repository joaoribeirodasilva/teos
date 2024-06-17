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

func New(options *DbOptions) *Db {

	if options == nil {
		slog.Error("failed to create database connection instance, no options specified")
		return nil
	}

	db := &Db{}

	db.ctx = context.TODO()
	if options.Ctx != nil {
		db.ctx = options.Ctx
	}

	db.dsn = options.Dsn
	db.name = options.Name
	if db.dsn == "" {
		db.protocol = options.Protocol
		db.hosts = options.Hosts
		db.username = options.Username
		db.password = options.Password
		db.options = options.Options
		if db.username != "" && db.password != "" {
			db.dsn = fmt.Sprintf("%s://%s:%s@%s/?%s", db.protocol, db.username, db.password, db.hosts, db.options)
		} else {
			db.dsn = fmt.Sprintf("%s://%s/?%s", db.protocol, db.hosts, db.options)
		}
	}

	return db
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
