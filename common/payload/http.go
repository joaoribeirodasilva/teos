package payload

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
)

type Http struct {
	conf          *configuration.Config
	Engine        *gin.Engine
	Request       *HttpRequest
	permissionsDb *redisdb.RedisDB
	sessionsDb    *redisdb.RedisDB
}

func NewHttp(conf *configuration.Config, engine *gin.Engine, permissionsDb *redisdb.RedisDB, sessionsDb *redisdb.RedisDB) *Http {

	return &Http{
		conf:          conf,
		Engine:        engine,
		Request:       nil,
		permissionsDb: permissionsDb,
		sessionsDb:    sessionsDb,
	}
}

func (h *Http) Parse(g *gin.Context) error {

	h.Request = NewRequest(h.conf, g, h.permissionsDb, h.sessionsDb)

	return h.Request.Parse()
}
