package payload

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
)

type Operation uint8

const (
	SVC_OPERATION_CREATE Operation = iota
	SVC_OPERATION_UPDATE
	SVC_OPERATION_DELETE
)

type Services struct {
	Db            *database.Db
	SessionsDb    *redisdb.RedisDB
	PermissionsDb *redisdb.RedisDB
	HistoryDb     *redisdb.RedisDB
	LogsDb        *redisdb.RedisDB
}

type Payload struct {
	started  Started
	Config   *configuration.Config
	Services Services
	Http     *Http
}

type Started struct {
	StartedDb            bool
	StartedSessionsDb    bool
	StartedPermissionsDb bool
	StartedHistoryDb     bool
	StartedLogsDb        bool
}

func (p *Payload) Start() error {

	if !p.hasConfig() {
		return fmt.Errorf("payload has no config")
	}

	if p.Services.Db != nil && !p.started.StartedDb {
		if err := p.Services.Db.Connect(); err != nil {
			return err
		}
		p.started.StartedDb = true
	}

	if p.Services.LogsDb != nil && !p.started.StartedLogsDb {
		if err := p.Services.LogsDb.Connect(); err != nil {
			return err
		}
		p.started.StartedLogsDb = true
	}

	if p.Services.SessionsDb != nil && !p.started.StartedSessionsDb {
		if err := p.Services.SessionsDb.Connect(); err != nil {
			return err
		}
		p.started.StartedSessionsDb = true
	}

	if p.Services.PermissionsDb != nil && !p.started.StartedPermissionsDb {
		if err := p.Services.PermissionsDb.Connect(); err != nil {
			return err
		}
		p.started.StartedPermissionsDb = true
	}

	if p.Services.HistoryDb != nil && !p.started.StartedHistoryDb {
		if err := p.Services.HistoryDb.Connect(); err != nil {
			return err
		}
		p.started.StartedHistoryDb = true
	}

	return nil
}

func (p *Payload) SetConfig(config *configuration.Config) {
	p.Config = config
}

func (p *Payload) SetDatabase(db *database.Db) error {

	if !p.hasConfig() {
		return fmt.Errorf("payload has no config")
	}

	if db == nil {
		dbConfig := p.Config.GetDatabase()
		dbOpts := &database.DbOptions{
			Ctx:      context.TODO(),
			Name:     dbConfig.Database,
			Host:     dbConfig.Host,
			Username: dbConfig.Username,
			Password: dbConfig.Password,
			Options:  dbConfig.Options,
		}
		db = database.New(dbOpts)
	}

	p.Services.Db = db

	return nil
}

func (p *Payload) SetSessionsDb(db *redisdb.RedisDB) error {

	if !p.hasConfig() {
		return fmt.Errorf("payload has no config")
	}

	if db == nil {
		dbConfig := p.Config.GetSessions()
		db = redisdb.New("sessions", dbConfig.Addr, dbConfig.Port, dbConfig.Db, dbConfig.Username, dbConfig.Password)
	}

	p.Services.SessionsDb = db

	return nil
}

func (p *Payload) SetPermissionsDb(db *redisdb.RedisDB) error {

	if !p.hasConfig() {
		return fmt.Errorf("payload has no config")
	}

	if db == nil {
		dbConfig := p.Config.GetPermissions()
		db = redisdb.New("permissions", dbConfig.Addr, dbConfig.Port, dbConfig.Db, dbConfig.Username, dbConfig.Password)
	}

	p.Services.PermissionsDb = db

	return nil
}

func (p *Payload) SetHistoryDb(db *redisdb.RedisDB) error {

	if !p.hasConfig() {
		return fmt.Errorf("payload has no config")
	}

	if db == nil {
		dbConfig := p.Config.GetHistory()
		db = redisdb.New("history", dbConfig.Addr, dbConfig.Port, dbConfig.Db, dbConfig.Username, dbConfig.Password)
	}

	p.Services.HistoryDb = db

	return nil
}

func (p *Payload) SetLogsDb(db *redisdb.RedisDB) error {

	if !p.hasConfig() {
		return fmt.Errorf("payload has no config")
	}

	if db == nil {
		dbConfig := p.Config.GetLogs()
		db = redisdb.New("logs", dbConfig.Addr, dbConfig.Port, dbConfig.Db, dbConfig.Username, dbConfig.Password)
	}

	p.Services.LogsDb = db

	return nil
}

func (p *Payload) SetHttp(engine *gin.Engine) error {

	if !p.hasConfig() {
		return fmt.Errorf("payload has no config")
	}

	p.Http = NewHttp(p.Config, engine)

	return nil
}

func (p *Payload) hasConfig() bool {
	return p.Config != nil
}
