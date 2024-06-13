package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
)

type Variables struct {
	Conf          *conf.Conf
	Db            *database.Db
	Configuration *configuration.Configuration
	User          *token.User
}

func MustGetAll(c *gin.Context) (*Variables, *service_errors.Error) {

	v := Variables{}
	ok := false

	co := c.MustGet("conf")
	v.Conf, ok = co.(*conf.Conf)
	if !ok {
		return nil, service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "", "Read", "invalid conf pointer").LogError()
	}

	d := c.MustGet("db")
	v.Db, ok = d.(*database.Db)
	if !ok {
		return nil, service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "", "Read", "invalid database pointer").LogError()
	}

	cf := c.MustGet("configuration")
	v.Configuration, ok = cf.(*configuration.Configuration)
	if !ok {
		return nil, service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "", "Read", "invalid configuration pointer").LogError()
	}

	v.User = nil
	a, exists := c.Get("auth")
	if exists {
		v.User, ok = a.(*token.User)
		if !ok {
			return nil, service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "Read", "invalid user pointer", "").LogError()
		}
	}

	return &v, nil
}
