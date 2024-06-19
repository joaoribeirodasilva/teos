package structures

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/environment"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
)

type Services struct {
	Gin           *gin.Engine
	Env           *environment.Environment
	Db            *database.Db
	Configuration *configuration.Configuration
	SessionsDB    *redisdb.RedisDB
	PermissionsDB *redisdb.RedisDB
	HistoryDB     *redisdb.RedisDB
	LogsDB        *redisdb.RedisDB
}

type RequestValues struct {
	Services *Services
	User     *token.User
	Query    requests.QueryString
}
