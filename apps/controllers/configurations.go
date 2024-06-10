package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
)

func AppConfigurationsList(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func AppConfigurationsGet(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func AppConfigurationsCreate(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func AppConfigurationsUpdate(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func AppConfigurationsDelete(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
