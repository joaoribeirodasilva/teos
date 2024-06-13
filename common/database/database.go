package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	conf *conf.Conf
	Conn *mongo.Client
	Db   *mongo.Database
}

func New(conf *conf.Conf) *Db {
	d := new(Db)
	d.conf = conf
	return d
}

func (d *Db) Connect() error {

	var err error

	dsn := ""
	if d.conf.Database.Username == "" && d.conf.Database.Password == "" {

		dsn = fmt.Sprintf(
			"mongodb://%s:%d/?%s",
			d.conf.Database.Host,
			d.conf.Database.Port,
			d.conf.Database.Options,
		)
	} else {
		dsn = fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/?%s",
			d.conf.Database.Username,
			d.conf.Database.Password,
			d.conf.Database.Host,
			d.conf.Database.Port,
			d.conf.Database.Options,
		)
	}

	slog.Info(fmt.Sprintf("[COMMON::DATABASE::Connect] connecting to database %s at %s:%d...\n", d.conf.Database.Database, d.conf.Database.Host, d.conf.Database.Port))

	d.Conn, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
	if err != nil {
		return service_log.Error(0, 0, "COMMON::DATABASE::Connect", "", "failed to connect to database %s at %s:%d. ERR: %s", d.conf.Database.Database, d.conf.Database.Host, d.conf.Database.Port, err.Error())
	}

	d.Db = d.Conn.Database(d.conf.Database.Database)
	if err := d.Conn.Ping(context.TODO(), nil); err != nil {
		return service_log.Error(0, 0, "COMMON::DATABASE::Connect", "", "failed to communicate with database %s at %s:%d. ERR: %s", d.conf.Database.Database, d.conf.Database.Host, d.conf.Database.Port, err.Error())
	}
	slog.Info("[COMMON::DATABASE::Connect] connected")
	return err
}
