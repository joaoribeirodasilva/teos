package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/database"
)

type Variables struct {
	Conf          *conf.Conf
	Db            *database.Db
	Configuration *configuration.Configuration
	// User *token.User
}

func MustGetAll(c *gin.Context) (*Variables, error) {

	v := Variables{}
	ok := false

	co := c.MustGet("conf")
	v.Conf, ok = co.(*conf.Conf)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, errors.New("invalid conf pointer")
	}

	d := c.MustGet("db")
	v.Db, ok = d.(*database.Db)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, errors.New("invalid database pointer")
	}

	cf := c.MustGet("configuration")
	v.Configuration, ok = cf.(*configuration.Configuration)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, errors.New("invalid configuration pointer")
	}

	/* v.User = nil
	a, exists := c.Get("auth")
	if exists {
		v.User = a.(*token.User)
	} */

	return &v, nil

}
