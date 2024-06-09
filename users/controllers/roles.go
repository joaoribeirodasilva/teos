package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
)

func UserRolesList(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func UserRolesGet(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func UserRolesCreate(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func UserRolesUpdate(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func UserRolesDelete(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
