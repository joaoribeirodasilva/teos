package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/apps/models"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/responses"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//TODO: Validation of data in fields

type RequestAppList struct {
	Text     string `json:"text"`
	Active   bool   `json:"active"`
	Page     int64  `json:"p"`
	PageSize int64  `json:"ps"`

	All bool `json:"a"`
}

func AppAppsList(c *gin.Context) {

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

	coll := vars.Db.Db.Collection("app_apps")

	count, err := coll.CountDocuments(context.TODO(), queryString.Filter)
	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppAppsList", "failed to count records. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	cursor, err := coll.Find(context.TODO(), queryString.Filter, queryString.Options)
	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppAppsList", "failed to query database. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	records := []models.AppApp{}
	if err := cursor.All(context.TODO(), &records); err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppAppsList", "failed to fetch database cursor. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	response := models.AppApps{
		Count: count,
		Rows:  &records,
	}

	c.JSON(http.StatusOK, &response)
}

func AppAppsGet(c *gin.Context) {

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

	record := models.AppApp{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppAppsGet", "failed to query database. ERR: %s", err.Error())
			appErr.LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, &record)
}

func AppAppsCreate(c *gin.Context) {

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	record := models.AppApp{}

	if err := c.ShouldBindBodyWithJSON(&record); err != nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppAppsCreate", "failed to bind request. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "appKey", Value: record.AppKey}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {

			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppsCreate", "", "failed to query database. ERR: %s", err.Error())
			appErr.LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
	} else {
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppAppsCreate", "appKey", "failed to query database. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	result, err := coll.InsertOne(context.TODO(), record)
	if err != nil {
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppAppsCreate", "", "failed to insert into database. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	response := responses.ResponseCreated{
		ID: result.InsertedID,
	}

	c.JSON(http.StatusCreated, response)
}

func AppAppsUpdate(c *gin.Context) {

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

	request := models.AppApp{}

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppAppsUpdate", "", "failed to bind request. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	record := models.AppApp{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "appKey", Value: request.AppKey}},
			bson.D{{Key: "_id", Value: bson.D{
				{Key: "$ne", Value: id},
			},
			}},
		}},
	}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppsUpdate", "", "failed to query database. ERR: %s", err.Error())
			appErr.LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	record.ID = *id
	record.Name = request.Name
	record.Description = request.Description
	record.AppKey = request.AppKey
	record.Active = request.Active
	record.UpdatedAt = time.Now().UTC()

	_, err := coll.UpdateOne(context.TODO(), bson.D{
		{Key: "_id", Value: id},
	}, bson.D{
		{Key: "$set", Value: record},
	})

	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppsUpdate", "", "failed to update into database. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	c.Status(http.StatusOK)
}

func AppAppsDelete(c *gin.Context) {

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

	record := models.AppApp{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppsDelete", "", "failed to query database. ERR: %s", err.Error())
			appErr.LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
