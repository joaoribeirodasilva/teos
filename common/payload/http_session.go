package payload

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
)

type HttpSession struct {
	ctx  *gin.Context
	conf *configuration.Config
	Auth *HttpCookie
}

func NewHttpSession(conf *configuration.Config, ctx *gin.Context, permissionsDb *redisdb.RedisDB, sessionsDb *redisdb.RedisDB, method string, route string) *HttpSession {
	s := &HttpSession{
		conf: conf,
		ctx:  ctx,
		Auth: NewHttpCookie(&conf.GetServices().Cookie, ctx, permissionsDb, sessionsDb, method, route),
	}
	return s
}

func (h *HttpSession) Parse() error {
	if err := h.Auth.Parse(); err != nil {
		return err
	}
	return nil
}
