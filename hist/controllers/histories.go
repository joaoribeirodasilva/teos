package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/responses"
	"github.com/joaoribeirodasilva/teos/hist/models"
	"github.com/joaoribeirodasilva/teos/hist/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HistoriesList(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func HistoriesGet(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func HistoriesCreate(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	request := requests.HistoriesCreate{}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	request.AuthenticationKey = strings.TrimSpace(request.AuthenticationKey)
	if request.AuthenticationKey == "" {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "authenticationKey",
				Message: "authenticationKey filed is required",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	//TODO: check authentication key

	request.Table = strings.TrimSpace(request.Table)
	if request.Table == "" {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "table",
				Message: "table filed is required",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	appId, _ := primitive.ObjectIDFromHex("666758b475cf5396aea26a13")
	coll := vars.Db.Db.Collection("hists_history")

	record := models.HistHistory{
		AppAppID:   appId,
		Table:      request.Table,
		OriginalID: request.OriginalID,
		Data:       request.Data,
		CreatedBy:  request.CreatedBy,
		CreatedAt:  request.CreatedAt,
		UpdatedBy:  request.UpdatedBy,
		UpdatedAt:  request.UpdatedAt,
		DeletedBy:  request.DeletedBy,
		DeletedAt:  request.DeletedAt,
	}

	_, err = coll.InsertOne(context.TODO(), &record)
	if err != nil {
		slog.Error(fmt.Sprintf("[HISTORY] failed to insert history for collection: %srecord id:%s", request.Table, request.OriginalID.Hex()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
