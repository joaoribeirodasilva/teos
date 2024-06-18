package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
)

func GetValues(c *gin.Context) (*structures.RequestValues, *logger.HttpError) {

	values, httpErr := MustGetAll(c)
	if httpErr != nil {
		return nil, httpErr
	}

	tq := requests.NewQueryString(c)
	values.Query = *tq
	if httpErr := values.Query.Bind(); httpErr != nil {
		return nil, httpErr
	}

	return values, nil

}

func MustGetAll(c *gin.Context) (*structures.RequestValues, *logger.HttpError) {

	v := structures.RequestValues{}
	var ok bool

	services := c.MustGet("services")
	v.Services, ok = services.(*structures.Services)
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
