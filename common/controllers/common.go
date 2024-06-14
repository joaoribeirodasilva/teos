package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
)

func GetValues(c *gin.Context) (*structures.RequestValues, *service_errors.Error) {

	vars, appErr := MustGetAll(c)
	if appErr != nil {
		return nil, appErr
	}

	queryString := requests.NewQueryString(c)
	if appErr := queryString.Bind(); appErr != nil {
		return nil, appErr
	}

	return &structures.RequestValues{
		Variables: vars,
		Query:     *queryString,
	}, nil
}

func MustGetAll(c *gin.Context) (*structures.Variables, *service_errors.Error) {

	v := structures.Variables{}
	ok := false

	co := c.MustGet("conf")
	v.Conf, ok = co.(*conf.Conf)
	if !ok {
		return nil, service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::MustGetAll", "conf", "invalid conf pointer")
	}

	d := c.MustGet("db")
	v.Db, ok = d.(*database.Db)
	if !ok {
		return nil, service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::MustGetAll", "db", "invalid database pointer")
	}

	cf := c.MustGet("configuration")
	v.Configuration, ok = cf.(*configuration.Configuration)
	if !ok {
		return nil, service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::MustGetAll", "configuration", "invalid configuration pointer")
	}

	v.User = nil
	a, exists := c.Get("user")
	if exists {
		v.User, ok = a.(*token.User)
		if !ok {
			return nil, service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::MustGetAll", "auth", "invalid user pointer")
		}
	}

	sDb := c.MustGet("sessionsDb")
	v.SessionsDB, ok = sDb.(*redisdb.RedisDB)
	if !ok {
		return nil, service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::MustGetAll", "logsDb", "invalid sessions database pointer")
	}

	pDb := c.MustGet("permissionsDb")
	v.PermissionsDB, ok = pDb.(*redisdb.RedisDB)
	if !ok {
		return nil, service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::MustGetAll", "permissionsDb", "invalid permissions database pointer")
	}

	return &v, nil
}
