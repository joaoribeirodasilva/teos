package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/consumers/perms/services"
)

func RoutesList(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func RoutesGet(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func RoutesReload(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewRoutesService(payload)

	if err := svc.Refresh(); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.Status(http.StatusOK)
}
