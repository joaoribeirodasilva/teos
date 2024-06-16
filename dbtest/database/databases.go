package database

import (
	"context"
	"fmt"
	"log/slog"
)

var (
	dbs map[string]*Db
)

func New(name string, options *DbOptions) *Db {

	if dbs == nil {
		dbs = make(map[string]*Db)
	}

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

	dbs[name] = db

	return dbs[name]
}

func GetDatabase(name string) *Db {
	db, ok := dbs[name]
	if !ok {
		return nil
	}
	return db
}

func GetAvailable() *[]string {

	if len(dbs) == 0 {
		return nil
	}

	keys := make([]string, 0)
	for k := range dbs {
		keys = append(keys, k)
	}

	return &keys
}

func GetServices() *map[string]*Db {
	return &dbs
}
