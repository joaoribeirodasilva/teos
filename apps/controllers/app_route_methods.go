package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	app_route_methods "github.com/joaoribeirodasilva/teos/apps/services/app_route_methods"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/responses"
)

func AppRouteMethodsList(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := app_route_methods.New(services)

	docs, err := svc.List("")
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

func AppRouteMethodsGet(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := app_route_methods.New(services)
	doc := &models.AppRouteMethod{}

	if err := svc.Get(doc, "id = ?", services.Query.ID); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, &doc)
}

func AppRouteMethodsCreate(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := app_route_methods.New(services)
	doc := &models.AppRouteMethod{}

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

func AppRouteMethodsUpdate(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := app_route_methods.New(services)
	doc := &models.AppRouteMethod{}
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

func AppRouteMethodsDelete(c *gin.Context) {

	services, err := controllers.GetValues(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := app_route_methods.New(services)
	doc := &models.AppRouteMethod{}

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
