package database

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbOptions struct {
	Ctx      context.Context
	Name     string
	Dsn      string
	Host     string
	Username string
	Password string
	Options  string
}

type Db struct {
	database *gorm.DB
	dsn      string
	host     string
	name     string
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
		db.host = options.Host
		db.username = options.Username
		db.password = options.Password
		db.options = options.Options
		if db.username != "" && db.password != "" {
			db.dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", db.username, db.password, db.host, db.name, db.options)
		} else {
			db.dsn = fmt.Sprintf("@tcp(%s)/%s?%s", db.host, db.name, db.options)
		}
	}

	return db
}

func (db *Db) Connect() error {

	var err error

	slog.Info("connecting to database...")

	dsn := mysql.Config{
		DSN: db.dsn,
	}

	db.database, err = gorm.Open(mysql.New(dsn), &gorm.Config{})
	if err != nil {
		slog.Error(fmt.Sprintf("failed to connect to database, Err: %s", err.Error()))
		return err
	}

	slog.Info("database connected successfully")

	return nil
}

func (db *Db) GetContext() context.Context {
	return db.ctx
}

func (db *Db) GetDatabase() *gorm.DB {
	return db.database
}
