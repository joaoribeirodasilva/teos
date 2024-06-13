package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/responses"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserUsersList(c *gin.Context) {

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

	coll := vars.Db.Db.Collection("user_users")

	count, err := coll.CountDocuments(context.TODO(), queryString.Filter)
	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "UserUsersList", "", "failed to count records. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	cursor, err := coll.Find(context.TODO(), queryString.Filter, queryString.Options)
	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "UserUsersList", "", "failed to query database. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	records := []models.UserUser{}
	if err := cursor.All(context.TODO(), &records); err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "UserUsersList", "", "failed to fetch database cursor. ERR: %s", err.Error())
		appErr.LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	response := models.UserUsers{
		Count: count,
		Rows:  &records,
	}

	c.JSON(http.StatusOK, &response)
}

func UserUsersGet(c *gin.Context) {

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
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "UserUsersGet", "id", "invalid object id in path").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	record := models.UserUser{}
	coll := vars.Db.Db.Collection("user_users")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLER", "UserUsersGet", "", "failed to query database. ERR: %s", err.Error()).LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, &record)
}

func UserUsersCreate(c *gin.Context) {

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	record := models.UserUser{}

	if err := c.ShouldBindBodyWithJSON(&record); err != nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "UserUsersCreate", "", "failed to bind request. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	coll := vars.Db.Db.Collection("user_users")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: record.Email}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUsersCreate", "", "failed to query database. ERR: %s", err.Error()).LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
	} else {
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "UserUsersCreate", "email", "an application with the same email already exists. id:%s", record.ID.Hex()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	if appErr := UserUserValidate(&record, true); appErr != nil {
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
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "UserUsersCreate", "", "failed to insert into database. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	response := responses.ResponseCreated{
		ID: result.InsertedID,
	}

	c.JSON(http.StatusCreated, response)
}

func UserUsersUpdate(c *gin.Context) {

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
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "UserUsersUpdate", "id", "invalid object id in path").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	request := models.UserUser{}

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "UserUsersUpdate", "", "failed to bind request. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	record := models.UserUser{}
	coll := vars.Db.Db.Collection("user_users")
	if err := coll.FindOne(context.TODO(), bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "email", Value: request.Email}},
			bson.D{{Key: "_id", Value: bson.D{
				{Key: "$ne", Value: id},
			},
			}},
		}},
	}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUsersUpdate", "", "failed to query database. ERR: %s", err.Error()).LogError()
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
	} else {
		appErr := service_errors.New(0, http.StatusConflict, "CONTROLLER", "UserUsersUpdate", "appKey", "an application with the same appKey already exists. id:%s", record.ID.Hex()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	record.FirstName = request.FirstName
	record.Surename = request.Surename
	record.Email = request.Email
	record.AvatarUrl = request.AvatarUrl
	// record.Password is only set via auth reset
	// record.Terms not set as it is set internally only
	// record.EmailVerified not set as it is set internally only
	record.Active = request.Active
	record.UpdatedBy = vars.User.ID
	record.UpdatedAt = time.Now().UTC()

	if appErr := UserUserValidate(&record, true); appErr != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	_, err := coll.UpdateOne(context.TODO(), bson.D{
		{Key: "_id", Value: id},
	}, bson.D{
		{Key: "$set", Value: record},
	})

	if err != nil {
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUsersUpdate", "", "failed to update into database. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	c.Status(http.StatusOK)
}

func UserUsersDelete(c *gin.Context) {

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
		appErr := service_errors.New(0, http.StatusBadRequest, "CONTROLLER", "UserUsersDelete", "id", "invalid object id in path").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	record := models.UserUser{}
	coll := vars.Db.Db.Collection("user_users")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUsersDelete", "", "failed to query database. ERR: %s", err.Error()).LogError()
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
		appErr := service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUsersDelete", "", "failed to delete from database. ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	c.Status(http.StatusOK)
}

func UserUserValidate(record *models.UserUser, create bool) *service_errors.Error {

	validate := validator.New()
	if err := validate.Var(record.FirstName, "required,gte=2"); err != nil {
		return service_errors.New(0, http.StatusConflict, "CONTROLLER", "UserUserValidate", "firstName", "invalid user first name").LogError()
	}

	if err := validate.Var(record.Surename, "required,gte=2"); err != nil {
		return service_errors.New(0, http.StatusConflict, "CONTROLLER", "UserUserValidate", "appKey", "invalid user surename").LogError()
	}

	if err := validate.Var(record.Email, "required,email"); err != nil {
		return service_errors.New(0, http.StatusConflict, "CONTROLLER", "UserUserValidate", "appKey", "invalid user email address").LogError()
	}

	if err := validate.Var(record.CreatedBy, "required,mongodbId"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUserValidate", "createdBy", "invalid created by").LogError()
	}

	if err := validate.Var(record.Password, "required"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUserValidate", "createdBy", "invalid created by").LogError()
	}

	if err := validate.Var(record.CreatedAt, "required"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUserValidate", "createdAt", "invalid created date").LogError()
	}

	if err := validate.Var(record.CreatedBy, "required,mongodbId"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUserValidate", "updatedBy", "invalid updated by").LogError()
	}

	if err := validate.Var(record.UpdatedAt, "required"); err != nil {
		return service_errors.New(0, http.StatusInternalServerError, "CONTROLLER", "UserUserValidate", "updatedAt", "invalid updated at").LogError()
	}

	return nil
}
