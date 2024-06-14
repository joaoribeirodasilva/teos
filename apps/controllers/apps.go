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
	"github.com/joaoribeirodasilva/teos/common/service_log"
	history_models "github.com/joaoribeirodasilva/teos/hist/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AppAppsList(c *gin.Context) {

	values, appErr := controllers.GetValues(c)
	if appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	docs := &models.AppApps{}
	docs.SetValues(values)

	if appErr := docs.FindAll(values.Query.Filter); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}
	if docs.Rows == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, &docs)
}

func AppAppsGet(c *gin.Context) {

	values, appErr := controllers.GetValues(c)
	if appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	id := values.Query.ID
	if id == nil {
		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::AppAppsGet", "id", "invalid object id in path")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	doc := &models.AppApp{}
	doc.SetValues(values)

	if appErr := doc.FindByID(*id); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	if doc.NotFound {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, &doc)
}

func AppAppsCreate(c *gin.Context) {

	values, appErr := controllers.GetValues(c)
	if appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	doc := models.AppApp{}
	doc.SetValues(values)

	// Binds the received JSON with the model and validates it
	if appErr := doc.Bind(c); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	// TODO: how to do if user is nill to bring the system user

	if appErr := doc.Create(); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	response := responses.ResponseCreated{
		ID: doc.ID,
	}

	c.JSON(http.StatusCreated, response)
}

func AppAppsUpdate(c *gin.Context) {

	values, appErr := controllers.GetValues(c)
	if appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	id := values.Query.ID
	if id == nil {
		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::AppAppsUpdate", "id", "invalid object id in path")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	doc := models.AppApp{}
	doc.SetValues(values)

	// Binds the received JSON with the model and validates it
	if appErr := doc.Bind(c); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	if appErr := doc.Update(*id); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
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
		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::AppAppsDelete", "id", "invalid object id in path")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	record := models.AppApp{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AppAppsDelete", "", "failed to query database. ERR: %s", err.Error())
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	oldRecord := record

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
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AppAppsDelete", "", "failed to delete from database. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	history := history_models.HistHistory{
		ID:         primitive.NewObjectID(),
		AppAppID:   vars.Configuration.ID,
		Collection: "app_apps",
		OriginalID: record.ID,
		Data:       oldRecord,
		CreatedBy:  vars.User.ID,
		CreatedAt:  now,
		UpdatedBy:  vars.User.ID,
		UpdatedAt:  now,
	}

	collHist := vars.Db.Db.Collection("hist_history")
	if _, err := collHist.InsertOne(context.TODO(), &history); err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AppAppsDelete", "", "failed to insert record history into database. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	c.Status(http.StatusOK)
}

func AppAppValidate(record *models.AppApp, create bool) *service_errors.Error {

	validate := validator.New()
	if err := validate.Var(record.Name, "required"); err != nil {
		return service_log.Error(0, http.StatusConflict, "CONTROLLER::AppAppValidate", "name", "invalid application name")
	}

	if err := validate.Var(record.AppKey, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusConflict, "CONTROLLER::AppAppValidate", "appKey", "invalid application key")
	}

	if err := validate.Var(record.CreatedBy, "required,mongodbId"); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AppAppValidate", "createdBy", "invalid created by")
	}

	if err := validate.Var(record.CreatedAt, "required"); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AppAppValidate", "createdAt", "invalid created date")
	}

	if err := validate.Var(record.CreatedBy, "required,mongodbId"); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AppAppValidate", "updatedBy", "invalid updated by")
	}

	if err := validate.Var(record.UpdatedAt, "required"); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AppAppValidate", "updatedAt", "invalid updated at")
	}

	return nil
}
