package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/apps/models"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/responses"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppAppsList", "", "failed to count records. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	cursor, err := coll.Find(context.TODO(), queryString.Filter, queryString.Options)
	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppAppsList", "", "failed to query database. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	records := []models.AppApp{}
	if err := cursor.All(context.TODO(), &records); err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppAppsList", "", "failed to fetch database cursor. ERR: %s", err.Error())
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
	if id == nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppAppsGet", "id", "invalid object id in path").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	record := models.AppApp{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppAppsGet", "", "failed to query database. ERR: %s", err.Error()).LogError()
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
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppAppsCreate", "", "failed to bind request. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "appKey", Value: record.AppKey}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppsCreate", "", "failed to query database. ERR: %s", err.Error()).LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
	} else {
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppAppsCreate", "appKey", "an application with the same appKey already exists. id:%s", record.ID.Hex()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	if appErr := AppAppValidate(&record, true); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	now := time.Now().UTC()
	record.ID = primitive.NewObjectID()
	record.CreatedBy = vars.User.ID
	record.CreatedAt = now
	record.UpdatedBy = vars.User.ID
	record.UpdatedAt = now

	result, err := coll.InsertOne(context.TODO(), record)
	if err != nil {
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppAppsCreate", "", "failed to insert into database. ERR: %s", err.Error()).LogError()
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
	if id == nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppAppsUpdate", "id", "invalid object id in path").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	request := models.AppApp{}

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppAppsUpdate", "", "failed to bind request. ERR: %s", err.Error()).LogError()
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
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppsUpdate", "", "failed to query database. ERR: %s", err.Error()).LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
	} else {
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppAppsUpdate", "appKey", "an application with the same appKey already exists. id:%s", record.ID.Hex()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	record.Name = request.Name
	record.Description = request.Description
	record.AppKey = request.AppKey
	record.Active = request.Active
	record.UpdatedBy = vars.User.ID
	record.UpdatedAt = time.Now().UTC()

	if appErr := AppAppValidate(&record, true); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	_, err := coll.UpdateOne(context.TODO(), bson.D{
		{Key: "_id", Value: id},
	}, bson.D{
		{Key: "$set", Value: record},
	})

	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppsUpdate", "", "failed to update into database. ERR: %s", err.Error()).LogError()
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
	if id == nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppAppsDelete", "id", "invalid object id in path").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	record := models.AppApp{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppsDelete", "", "failed to query database. ERR: %s", err.Error()).LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	now := time.Now().UTC()
	record.UpdatedBy = vars.User.ID
	record.UpdatedAt = now
	record.DeletedBy = &vars.User.ID
	record.DeletedAt = &now

	_, err := coll.UpdateOne(context.TODO(), bson.D{
		{Key: "_id", Value: id},
	}, bson.D{
		{Key: "$set", Value: record},
	})

	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppsDelete", "", "failed to delete from database. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	c.Status(http.StatusOK)
}

func AppAppValidate(record *models.AppApp, create bool) *service_errors.Error {

	validate := validator.New()
	if err := validate.Var(record.Name, "required"); err != nil {
		return service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppAppValidate", "name", "invalid application name").LogError()
	}

	if err := validate.Var(record.AppKey, "required,gte=1"); err != nil {
		return service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppAppValidate", "appKey", "invalid application key").LogError()
	}

	if err := validate.Var(record.CreatedBy, "required,mongodbId"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppValidate", "createdBy", "invalid created by").LogError()
	}

	if err := validate.Var(record.CreatedAt, "required"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppValidate", "createdAt", "invalid created date").LogError()
	}

	if err := validate.Var(record.CreatedBy, "required,mongodbId"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppValidate", "updatedBy", "invalid updated by").LogError()
	}

	if err := validate.Var(record.UpdatedAt, "required"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppAppValidate", "updatedAt", "invalid updated at").LogError()
	}

	return nil
}
