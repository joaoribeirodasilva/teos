package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/apps/models"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AppAppsList(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	pageSize := int64(100)
	page := int64(0)

	skip := page * pageSize

	coll := vars.Db.Db.Collection("app_apps")

	count, err := coll.EstimatedDocumentCount(context.TODO())
	if err != nil {
		slog.Error(fmt.Sprintf("[APP_APPS::LIST] failed to count records. ERR: %s", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	options := options.FindOptions{Limit: &pageSize, Skip: &skip}

	cursor, err := coll.Find(context.TODO(), bson.D{{}}, &options)
	if err != nil {
		slog.Error(fmt.Sprintf("[APP_APPS::LIST] failed to query database. ERR: %s", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	records := []models.AppApp{}
	if err := cursor.All(context.TODO(), &records); err != nil {
		slog.Error(fmt.Sprintf("[APP_APPS::LIST] failed to fetch database cursor.ERR: %s", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := models.AppApps{
		Count: count,
		Rows:  &records,
	}

	c.JSON(http.StatusOK, &response)
}

func AppAppsGet(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	strId := c.Param("id")
	id, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		slog.Error(fmt.Sprintf("[APP_APPS::Get] invalid id. ERR: %s", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	record := models.AppApp{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			slog.Error(fmt.Sprintf("[APP_APPS::Get] failed to query database. ERR: %s", err.Error()))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, &record)
}

func AppAppsCreate(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	record := models.AppApp{}

	if err := c.ShouldBindBodyWithJSON(&record); err != nil {
		slog.Error(fmt.Sprintf("[APP_APPS::Create] failed to bind request. ERR: %s", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "appKey", Value: record.AppKey}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			slog.Error(fmt.Sprintf("[APP_APPS::Create] failed to query database. ERR: %s", err.Error()))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	result, err := coll.InsertOne(context.TODO(), record)
	if err != nil {
		slog.Error(fmt.Sprintf("[APP_APPS::Create] failed to insert into database. ERR: %s", err.Error()))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := responses.ResponseCreated{
		ID: result.InsertedID,
	}

	c.JSON(http.StatusCreated, response)
}

func AppAppsUpdate(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	strId := c.Param("id")
	id, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		slog.Error(fmt.Sprintf("[APP_APPS::Update] invalid id. ERR: %s", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	request := models.AppApp{}

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		slog.Error(fmt.Sprintf("[APP_APPS::Update] failed to bind request. ERR: %s", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
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
			slog.Error(fmt.Sprintf("[APP_APPS::Update] failed to query database. ERR: %s", err.Error()))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	record.ID = id
	record.Name = request.Name
	record.Description = request.Description
	record.AppKey = request.AppKey
	record.Active = request.Active
	record.UpdatedAt = time.Now().UTC()

	_, err = coll.UpdateOne(context.TODO(), bson.D{
		{Key: "_id", Value: id},
	}, bson.D{
		{Key: "$set", Value: record},
	})

	if err != nil {
		slog.Error("[APP_APPS::Update] failed to update into database")
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.Status(http.StatusOK)
}

func AppAppsDelete(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	strId := c.Param("id")
	id, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		slog.Error(fmt.Sprintf("[APP_APPS::Delete] invalid id. ERR: %s", err.Error()))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	record := models.AppApp{}
	coll := vars.Db.Db.Collection("app_apps")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			slog.Error(fmt.Sprintf("[APP_APPS::Delete] failed to query database. ERR: %s", err.Error()))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
