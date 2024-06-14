package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HistoriesList(c *gin.Context) {

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	queryString := requests.NewQueryString(c)
	if err := queryString.Bind(); err != nil {
		c.AbortWithStatusJSON(err.HttpCode, err)
		return
	}

	coll := vars.Db.Db.Collection("hist_history")

	count, err := coll.CountDocuments(context.TODO(), queryString.Filter)
	if err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLER::HistoriesList", "", "failed to count records. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	cursor, err := coll.Find(context.TODO(), queryString.Filter, queryString.Options)
	if err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLER::HistoriesList", "", "failed to query database. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	records := []models.HistHistory{}
	if err := cursor.All(context.TODO(), &records); err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLER::HistoriesList", "", "failed to fetch database cursor. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	response := models.HistHistories{
		Count: count,
		Rows:  &records,
	}

	c.JSON(http.StatusOK, &response)
}

func HistoriesGet(c *gin.Context) {

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	queryString := requests.NewQueryString(c)
	if appErr := queryString.Bind(); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	id := queryString.ID
	if id == nil {
		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::HistoriesGet", "id", "invalid object id in path")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	record := models.HistHistory{}
	coll := vars.Db.Db.Collection("hist_history")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLER::HistoriesGet", "", "failed to query database. ERR: %s", err.Error())
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

/* func HistoriesCreate(c *gin.Context) {

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	request := models.HistHistory{}

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
} */
