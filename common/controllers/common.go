package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
)

func GetValues(c *gin.Context) (*structures.RequestValues, *logger.HttpError) {

	vars, httpErr := MustGetAll(c)
	if httpErr != nil {
		return nil, httpErr
	}

	queryString := requests.NewQueryString(c)
	if httpErr := queryString.Bind(); httpErr != nil {
		return nil, httpErr
	}

	return &structures.RequestValues{
		Variables: vars,
		Query:     *queryString,
	}, nil
}

func MustGetAll(c *gin.Context) (*structures.Variables, *logger.HttpError) {

	v := structures.Variables{}
	ok := false

	co := c.MustGet("conf")
	v.Conf, ok = co.(*conf.Conf)
	if !ok {
		err := errors.New("invalid conf pointer")
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "invalid conf pointer", err, nil)
	}

	d := c.MustGet("db")
	v.Db, ok = d.(*database.Db)
	if !ok {
		err := errors.New("invalid database pointer")
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "invalid database pointer", err, nil)
	}

	cf := c.MustGet("configuration")
	v.Configuration, ok = cf.(*configuration.Configuration)
	if !ok {
		err := errors.New("invalid configuration pointer")
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "invalid configuration pointer", err, nil)
	}

	sDb := c.MustGet("sessionsDb")
	v.SessionsDB, ok = sDb.(*redisdb.RedisDB)
	if !ok {
		err := errors.New("invalid sessions database pointer")
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "invalid sessions database pointer", err, nil)
	}

	pDb := c.MustGet("permissionsDb")
	v.PermissionsDB, ok = pDb.(*redisdb.RedisDB)
	if !ok {
		err := errors.New("invalid permissions database pointer")
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "invalid permissions database pointer", err, nil)
	}

	v.User = nil
	a, exists := c.Get("user")
	if exists {
		v.User, ok = a.(*token.User)
		if !ok {
			err := errors.New("invalid user pointer")
			return nil, logger.Error(logger.LogStatusInternalServerError, nil, "invalid user pointer", err, nil)
		}
	}

	return &v, nil
}
