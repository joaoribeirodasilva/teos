package models

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseModels struct {
	values         *structures.RequestValues
	location       string
	collectionName string
	collection     *mongo.Collection
}

func (m *BaseModels) Init(c *gin.Context, location string, collectionName string) *service_errors.Error {
	values, appErr := controllers.GetValues(c)
	if appErr != nil {
		return appErr
	}
	m.values = values
	m.collectionName = collectionName
	m.location = location
	m.collection = m.values.Variables.Db.Db.Collection(m.collectionName)

	return nil
}

func (m *BaseModels) GetValues() *structures.RequestValues {
	return m.values
}

func (m *BaseModels) FindAllBuQueryFilter(model interface{}, opts ...*options.FindOptions) (int64, *service_errors.Error) {
	return m.FindAll(m.values.Query.Filter, model, opts...)
}

func (m *BaseModels) FindAll(filter interface{}, model interface{}, opts ...*options.FindOptions) (int64, *service_errors.Error) {

	count, appErr := m.CountDocs(filter, nil)
	if appErr != nil {
		return 0, appErr
	}

	if count == 0 {
		return 0, nil
	}

	cursor, err := m.collection.Find(context.TODO(), filter, opts...)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return 0, service_log.Error(0, http.StatusNotFound, m.location, "", "documents not found")
		}
		return 0, service_log.Error(0, http.StatusInternalServerError, m.location, "", "failed to query database. ERR: %s", err.Error())
	}

	if err := cursor.All(context.TODO(), model); err != nil {
		return 0, service_log.Error(0, http.StatusInternalServerError, m.location, "", "failed to fetch database cursor. ERR: %s", err.Error())
	}

	return count, nil
}

func (m *BaseModels) CountDocs(filter interface{}, opts ...*options.CountOptions) (int64, *service_errors.Error) {

	count, err := m.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApps", "", "failed to count records. ERR: %s", err.Error())
		return 0, appErr
	}
	return count, nil
}
