package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/payload"
)

func GetPayload(c *gin.Context) (*payload.Payload, *logger.HttpError) {

	serviceInterface := c.MustGet("service")
	services, ok := serviceInterface.(*payload.Payload)
	if !ok {
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "invalid service pointer", nil, nil)
	}

	return services, nil
}
