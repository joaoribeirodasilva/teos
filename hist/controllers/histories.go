package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/hist/services"
)

func HistoryList(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewHistoryService(payload)

	docs, err := svc.List("")
	if err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

func HistoryGet(c *gin.Context) {

	payload, err := controllers.GetPayload(c)
	if err != nil {
		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	svc := services.NewHistoryService(payload)
	doc := &models.History{}

	if err := svc.Get(doc, "id = ?", payload.Http.Request.ID); err != nil {

		c.AbortWithStatusJSON(int(err.Status), err)
		return
	}

	c.JSON(http.StatusOK, &doc)
}

// func HistoryCreate(c *gin.Context) {

// 	services, err := controllers.GetValues(c)
// 	if err != nil {

// 		c.AbortWithStatusJSON(int(err.Status), err)
// 		return
// 	}

// 	svc := history.New(services)
// 	doc := &models.History{}

// 	if err := c.ShouldBindBodyWithJSON(doc); err != nil {

// 		httpError := logger.Error(logger.LogStatusBadRequest, nil, "invalid JSON body", err, nil)
// 		c.AbortWithStatusJSON(int(httpError.Status), httpError)
// 		return
// 	}

// 	if err := svc.Create(doc); err != nil {

// 		c.AbortWithStatusJSON(int(err.Status), err)
// 		return
// 	}

// 	response := responses.ResponseCreated{
// 		ID: doc.ID,
// 	}

// 	c.JSON(http.StatusCreated, response)

// }
