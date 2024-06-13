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

func AppRoutesBlocksList(c *gin.Context) {

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

	coll := vars.Db.Db.Collection("app_routes_blocks")

	count, err := coll.CountDocuments(context.TODO(), queryString.Filter)
	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppRoutesBlocksList", "", "failed to count records. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	cursor, err := coll.Find(context.TODO(), queryString.Filter, queryString.Options)
	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppRoutesBlocksList", "", "failed to query database. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	records := []models.AppRoutesBlock{}
	if err := cursor.All(context.TODO(), &records); err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "AppRoutesBlocksList", "", "failed to fetch database cursor. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	response := models.AppRoutesBlocks{
		Count: count,
		Rows:  &records,
	}

	c.JSON(http.StatusOK, &response)
}

func AppRoutesBlocksGet(c *gin.Context) {

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
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppRoutesBlocksGet", "id", "invalid object id in path").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	record := models.AppRoutesBlock{}
	coll := vars.Db.Db.Collection("app_routes_blocks")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksGet", "", "failed to query database. ERR: %s", err.Error()).LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, &record)
}

func AppRoutesBlocksCreate(c *gin.Context) {

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	record := models.AppRoutesBlock{}

	if err := c.ShouldBindBodyWithJSON(&record); err != nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppRoutesBlocksCreate", "", "failed to bind request. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "appAppId", Value: record.ApplicationID}},
			bson.D{{Key: "route", Value: record.Route}},
		}},
	}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksCreate", "", "failed to query database. ERR: %s", err.Error()).LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
	} else {
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppRoutesBlocksCreate", "appAppId, route", "an application with the same appAppId and route already exists. id:%s", record.ID.Hex()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	if appErr := AppRoutesBlocksValidate(&record, true); appErr != nil {
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
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppRoutesBlocksCreate", "", "failed to insert into database. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	response := responses.ResponseCreated{
		ID: result.InsertedID,
	}

	c.JSON(http.StatusCreated, response)
}

func AppRoutesBlocksUpdate(c *gin.Context) {

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
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppRoutesBlocksUpdate", "id", "invalid object id in path").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	request := models.AppRoutesBlock{}

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppRoutesBlocksUpdate", "", "failed to bind request. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	record := models.AppRoutesBlock{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "appAppId", Value: request.ApplicationID}},
			bson.D{{Key: "route", Value: request.Route}},
			bson.D{{Key: "_id", Value: bson.D{
				{Key: "$ne", Value: id},
			},
			}},
		}},
	}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksUpdate", "", "failed to query database. ERR: %s", err.Error()).LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
	} else {
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppRoutesBlocksUpdate", "appAppId, route", "an application with the same appAppId and route already exists. id:%s", record.ID.Hex()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	now := time.Now().UTC()
	record.ApplicationID = request.ApplicationID
	record.Name = request.Name
	record.Description = request.Description
	record.Route = request.Route
	record.Active = request.Active
	record.UpdatedBy = vars.User.ID
	record.UpdatedAt = now

	if appErr := AppRoutesBlocksValidate(&record, true); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	_, err := coll.UpdateOne(context.TODO(), bson.D{
		{Key: "_id", Value: id},
	}, bson.D{
		{Key: "$set", Value: record},
	})
	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksUpdate", "", "failed to update into database. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	c.Status(http.StatusOK)
}

func AppRoutesBlocksDelete(c *gin.Context) {

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
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "AppRoutesBlocksDelete", "id", "invalid object id in path").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	record := models.AppApp{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksDelete", "", "failed to query database. ERR: %s", err.Error()).LogError()
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
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksDelete", "", "failed to delete from database. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	c.Status(http.StatusOK)
}

func AppRoutesBlocksValidate(record *models.AppRoutesBlock, create bool) *service_errors.Error {

	validate := validator.New()

	if err := validate.Var(record.ApplicationID, "required,mongodbId"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksValidate", "createdBy", "invalid created by").LogError()
	}

	if err := validate.Var(record.Name, "required,gte=1"); err != nil {
		return service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppRoutesBlocksValidate", "name", "invalid application name").LogError()
	}

	if err := validate.Var(record.Route, "required,gte=1"); err != nil {
		return service_errors.New(0, http.StatusConflict, "CONTROLLER", "AppRoutesBlocksValidate", "route", "invalid application key").LogError()
	}

	if err := validate.Var(record.CreatedBy, "required,mongodbId"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksValidate", "createdBy", "invalid created by").LogError()
	}

	if err := validate.Var(record.CreatedAt, "required"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksValidate", "createdAt", "invalid created date").LogError()
	}

	if err := validate.Var(record.CreatedBy, "required,mongodbId"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksValidate", "updatedBy", "invalid updated by").LogError()
	}

	if err := validate.Var(record.UpdatedAt, "required"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "AppRoutesBlocksValidate", "updatedAt", "invalid updated at").LogError()
	}

	return nil
}
