package controllers

import (
	"context"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/responses"
	"github.com/joaoribeirodasilva/teos/common/utils/password"
	"github.com/joaoribeirodasilva/teos/users/models"
	"github.com/joaoribeirodasilva/teos/users/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UserUsersList(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func UserUsersGet(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//id := c.Param("id")

	c.Status(http.StatusOK)
}

func UserUsersCreate(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var request requests.UserUsersCreate

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	request.FirstName = strings.TrimSpace(request.FirstName)
	if request.FirstName == "" {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "firstName",
				Message: "firstName is required",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	request.Surename = strings.TrimSpace(request.Surename)
	if request.FirstName == "" {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "surename",
				Message: "surename is required",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	_, err = mail.ParseAddress(request.Email)
	if err != nil {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "email",
				Message: "invalid email address",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	request.Password = strings.TrimSpace(request.Password)
	if request.Password == "" || len(request.Password) < 6 {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "password",
				Message: "the password must have at least 6 characters",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	request.CheckPassword = strings.TrimSpace(request.CheckPassword)
	if request.CheckPassword != request.Password {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "password",
				Message: "the passwords don't match",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !request.Terms {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "email",
				Message: "you must accept the terms of service",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	record := models.UserUser{}

	coll := vars.Db.Db.Collection("user_users")
	opts := options.FindOneOptions{Collation: &options.Collation{CaseLevel: false, Strength: 1}}
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: request.Email}}, &opts).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		response := responses.ResponseErrorMessage{
			Error: responses.ErrMessage{
				Message: "user already exists",
			},
		}
		c.AbortWithStatusJSON(http.StatusConflict, response)
		return
	}

	var tempPassword string
	record.FirstName = request.FirstName
	record.Surename = request.Surename
	record.Email = request.Email
	tempPassword, err = password.Hash(request.Password)
	record.Active = 0
	//TODO: replace with the right user

	now := time.Now().UTC()
	record.CreatedBy = 1
	record.CreatedAt = now
	record.UpdatedBy = 1
	record.UpdatedAt = now
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	record.Password = &tempPassword
	record.ID = primitive.NewObjectID()

	result, err := coll.InsertOne(context.TODO(), record)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := responses.ResponseCreated{
		ID: result.InsertedID,
	}

	c.JSON(http.StatusCreated, response)
}

func UserUsersUpdate(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func UserUsersDelete(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
