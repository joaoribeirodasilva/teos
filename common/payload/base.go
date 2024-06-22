package payload

import (
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
	Config   *configuration.Config
	Services Services
	Http     *Http
}
