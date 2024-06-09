package database

import (
	"fmt"
	"log/slog"

	"github.com/joaoribeirodasilva/teos/common/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Db struct {
	conf *conf.Conf
	Conn *gorm.DB
}

func New(conf *conf.Conf) *Db {
	d := new(Db)
	d.conf = conf
	return d
}

func (d *Db) Connect() error {

	var err error

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		d.conf.Database.Username,
		d.conf.Database.Password,
		d.conf.Database.Host,
		d.conf.Database.Port,
		d.conf.Database.Database,
		d.conf.Database.Options,
	)

	slog.Info(fmt.Sprintf("[DATABASE] connecting to database %s at %s:%d...\n", d.conf.Database.Database, d.conf.Database.Host, d.conf.Database.Port))
	d.Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error(fmt.Sprintf("[DATABASE] failed to connect to database %s at %s:%d. ERR: %s", d.conf.Database.Database, d.conf.Database.Host, d.conf.Database.Port, err.Error()))
		return err
	}
	slog.Info("[DATABASE] connected")
	return err
}
