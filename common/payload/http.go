package payload

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
)

type Http struct {
	conf    *configuration.Config
	Engine  *gin.Engine
	Request *HttpRequest
}

func NewHttp(conf *configuration.Config, engine *gin.Engine) *Http {

	return &Http{
		conf:    conf,
		Engine:  engine,
		Request: nil,
	}
}

func (h *Http) Parse(g *gin.Context) error {

	h.Request = NewRequest(h.conf, g)

	return h.Request.Parse()
}
