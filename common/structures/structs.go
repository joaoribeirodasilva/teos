package structures

import (
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
)

type Variables struct {
	Conf          *conf.Conf
	Db            *database.Db
	Configuration *configuration.Configuration
	User          *token.User
	SessionsDB    *redisdb.RedisDB
	PermissionsDB *redisdb.RedisDB
}

type RequestValues struct {
	Variables *Variables
	Query     requests.QueryString
}
