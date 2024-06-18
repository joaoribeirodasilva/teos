package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/responses"
	"github.com/joaoribeirodasilva/teos/users/services/groups"
)

func UserGroupsList(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := groups.New(services)

	docs, err := svc.List(nil)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

func UserGroupsGet(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := groups.New(services)
	doc := &models.UserGroupModel{}

	if err := svc.Get(nil, doc); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, &doc)
}

func UserGroupsCreate(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := groups.New(services)
	doc := &models.UserGroupModel{}

	if err := c.ShouldBindBodyWithJSON(doc); err != nil {

		httpError := logger.Error(logger.LogStatusBadRequest, nil, "invalid JSON body", err, nil)
		c.AbortWithStatusJSON(int(httpError.Status), httpError)
		return
	}

	if err := svc.Create(doc); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	response := responses.ResponseCreated{
		ID: doc.ID,
	}

	c.JSON(http.StatusCreated, response)
}

func UserGroupsUpdate(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := groups.New(services)
	doc := &models.UserGroupModel{}

	if err := c.ShouldBindBodyWithJSON(doc); err != nil {

		httpError := logger.Error(logger.LogStatusBadRequest, nil, "invalid JSON body", err, nil)
		c.AbortWithStatusJSON(int(httpError.Status), httpError)
		return
	}

	if err := svc.Update(doc); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.Status(http.StatusOK)
}

func UserGroupsDelete(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := groups.New(services)
	doc := &models.UserGroupModel{}

	if err := c.ShouldBindBodyWithJSON(doc); err != nil {

		httpError := logger.Error(logger.LogStatusBadRequest, nil, "invalid JSON body", err, nil)
		c.AbortWithStatusJSON(int(httpError.Status), httpError)
		return
	}

	if err := svc.Delete(doc); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.Status(http.StatusOK)
}
