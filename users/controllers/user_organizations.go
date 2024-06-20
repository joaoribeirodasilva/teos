package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/responses"
	user_organizations "github.com/joaoribeirodasilva/teos/users/services/user_organizations"
)

func UserOrganizationsList(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := user_organizations.New(services)

	docs, err := svc.List("")
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

func UserOrganizationsGet(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := user_organizations.New(services)
	doc := &models.UserOrganization{}

	if err := svc.Get(doc, "id = ?", services.Query.ID); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, &doc)
}

func UserOrganizationsCreate(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := user_organizations.New(services)
	doc := &models.UserOrganization{}

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

func UserOrganizationsUpdate(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := user_organizations.New(services)
	doc := &models.UserOrganization{}
	doc.ID = *services.Query.ID

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

func UserOrganizationsDelete(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := user_organizations.New(services)
	doc := &models.UserOrganization{}

	if err := c.ShouldBindBodyWithJSON(doc); err != nil {

		httpError := logger.Error(logger.LogStatusBadRequest, nil, "invalid JSON body", err, nil)
		c.AbortWithStatusJSON(int(httpError.Status), httpError)
		return
	}

	if err := svc.Delete(doc.ID); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.Status(http.StatusOK)

}
