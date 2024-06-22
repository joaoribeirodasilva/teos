package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/apps/services"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/responses"
)

func AppRouteMethodsList(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAppRouteMethodsService(payload)

	docs, err := svc.List("")
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

func AppRouteMethodsGet(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAppRouteMethodsService(payload)
	doc := &models.AppRouteMethod{}

	if err := svc.Get(doc, "id = ?", payload.Http.Request.ID); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, &doc)
}

func AppRouteMethodsCreate(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAppRouteMethodsService(payload)
	doc := &models.AppRouteMethod{}

	if err := payload.Http.Request.Bind(doc); err != nil {

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

	payload, err := controllers.GetPayload(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAppRouteMethodsService(payload)
	doc := &models.AppRouteMethod{}
	doc.ID = payload.Http.Request.ID

	if err := payload.Http.Request.Bind(doc); err != nil {

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

	payload, err := controllers.GetPayload(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAppRouteMethodsService(payload)

	if err := svc.Delete(payload.Http.Request.ID); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.Status(http.StatusOK)

}
