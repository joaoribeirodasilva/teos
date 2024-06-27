package payload

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
)

type HttpRequest struct {
	ctx     *gin.Context
	conf    *configuration.Config
	Session *HttpSession
	Query   *HttpQuery
	ID      uint
}

func NewRequest(conf *configuration.Config, ctx *gin.Context, permissionsDb *redisdb.RedisDB, sessionsDb *redisdb.RedisDB) *HttpRequest {

	method := ctx.Request.Method
	route := ctx.FullPath()

	req := &HttpRequest{
		ctx:     ctx,
		conf:    conf,
		Session: NewHttpSession(conf, ctx, permissionsDb, sessionsDb, method, route),
		Query:   NewHttpQuery(ctx),
	}

	return req
}

func (h *HttpRequest) Parse() error {

	if err := h.parseID(); err != nil {
		return err
	}

	if err := h.Session.Parse(); err != nil {
		return err
	}

	if err := h.Query.Parse(); err != nil {
		return err
	}

	return nil
}

func (h *HttpRequest) Bind(obj any) error {

	if err := h.ctx.BindJSON(obj); err != nil {
		return err
	}

	return nil
}

func (h *HttpRequest) parseID() error {

	strId := h.ctx.Params.ByName("id")
	if strId != "" {
		tempId, err := strconv.Atoi(strId)
		if err != nil {
			return fmt.Errorf("invalid id")
		}
		h.ID = uint(tempId)
	}

	return nil

}
