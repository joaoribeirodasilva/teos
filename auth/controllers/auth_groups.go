package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/auth/services"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/responses"
)

func AuthGroupsList(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAuthGroupsService(payload)

	docs, err := svc.List("")
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

func AuthGroupsGet(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAuthGroupsService(payload)
	doc := &models.AuthGroup{}

	if err := svc.Get(doc, "id = ?", payload.Http.Request.ID); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, &doc)
}

func AuthGroupsCreate(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAuthGroupsService(payload)
	doc := &models.AuthGroup{}

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

func AuthGroupsUpdate(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAuthGroupsService(payload)
	doc := &models.AuthGroup{}
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

func AuthGroupsDelete(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewAuthGroupsService(payload)

	if err := svc.Delete(payload.Http.Request.ID); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.Status(http.StatusOK)

}
