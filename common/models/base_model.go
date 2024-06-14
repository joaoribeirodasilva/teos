package models

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseModel struct {
	ID             primitive.ObjectID  `json:"_id" bson:"_id"`
	CreatedBy      primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt      time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy      primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt      time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy      *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt      *time.Time          `json:"deletedAt" bson:"deletedAt"`
	values         *structures.RequestValues
	location       string
	collectionName string
	collection     *mongo.Collection
	ctx            *gin.Context
}

func (m *BaseModel) Init(c *gin.Context, location string, collectionName string) *service_errors.Error {
	values, appErr := controllers.GetValues(c)
	if appErr != nil {
		return appErr
	}
	m.ctx = c
	m.values = values
	m.collectionName = collectionName
	m.location = location
	m.collection = m.values.Variables.Db.Db.Collection(m.collectionName)

	return nil
}

func (m *BaseModel) GetValues() *structures.RequestValues {
	return m.values
}

func (m *BaseModel) Bind(model interface{}, c *gin.Context) *service_errors.Error {

	// Binds the received JSON with the document
	if err := m.ctx.ShouldBindBodyWithJSON(model); err != nil {
		return service_log.Error(0, http.StatusBadRequest, m.location+"::Bind", "", "failed to bind request. ERR: %s", err.Error())
	}

	return nil
}

func (m *BaseModel) FindByQueryID(model interface{}, opts ...*options.FindOneOptions) *service_errors.Error {

	if m.values.Query.ID == nil {
		return service_log.Error(0, http.StatusNotFound, m.location+"::FinfQueryByID", "", "document not found")
	}
	return m.FindByID(*m.values.Query.ID, model, opts...)
}

func (m *BaseModel) FindByID(id primitive.ObjectID, model interface{}, opts ...*options.FindOneOptions) *service_errors.Error {

	return m.First(bson.D{{Key: "_id", Value: id}}, model, opts...)
}

func (m *BaseModel) First(filter interface{}, model interface{}, opts ...*options.FindOneOptions) *service_errors.Error {

	if err := m.collection.FindOne(context.TODO(), filter, opts...).Decode(model); err != nil {
		if err != mongo.ErrNoDocuments {
			return service_log.Error(0, http.StatusNotFound, m.location+"::First", "", "document not found")
		}
		return service_log.Error(0, http.StatusNotFound, m.location+"::First", "", "failed to query database. ERR: %s", err.Error())
	}
	return nil
}

func (m *BaseModel) Create(uniqueFilter bson.D, model interface{}, unique string, opts ...*options.InsertOneOptions) *service_errors.Error {

	exists := BaseModel{}
	exists.Init(m.ctx, m.location, m.collectionName)

	if err := m.collection.FindOne(context.TODO(), uniqueFilter).Decode(&exists); err != nil {
		if err != mongo.ErrNoDocuments {
			return service_log.Error(0, http.StatusInternalServerError, m.location+"::Create", "", "failed to query database. ERR: %s", err.Error())
		}
	} else {
		return service_log.Error(0, http.StatusConflict, m.location+"::Create", "", "document already exists with id: %s", exists.ID.Hex())
	}

	if _, err := m.collection.InsertOne(context.TODO(), model, opts...); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, m.location+"::Create", "", "failed to insert document into database. ERR: %s", err.Error())
	}

	return nil
}

func (m *BaseModel) Update(filter bson.D, model interface{}, unique string, opts ...*options.InsertOneOptions) *service_errors.Error {

	if _, err := m.collection.UpdateOne(context.TODO(), filter, bson.D{
		{Key: "$set", Value: model},
	}); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, m.location+"::Update", "", "failed to update into database. ERR: %s", err.Error())
	}
	return nil
}

func (m *BaseModel) Delete(id primitive.ObjectID, opts ...*options.DeleteOptions) *service_errors.Error {

	exists := BaseModel{}
	exists.Init(m.ctx, m.location, m.collectionName)
	if err := m.collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&exists); err != nil {
		if err == mongo.ErrNoDocuments {
			return service_log.Error(0, http.StatusInternalServerError, m.location+"::Delete", "", "document not found")
		} else {
			return service_log.Error(0, http.StatusInternalServerError, m.location+"::Delete", "", "failed to query database. ERR: %s", err.Error())
		}
	}

	if _, err := m.collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}}); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, m.location+"::Delete", "", "failed to delete document from database. ERR: %s", err.Error())
	}

	return nil
}

func (m *BaseModel) GetLocation() string {
	return m.location
}

func (m *BaseModel) SetLocation(location string) {
	m.location = location
}

func (m *BaseModel) SetCollectionName(collectionName string) {
	m.collectionName = collectionName
	//m.collection = m.values.Variables.Db.Db.Collection(m.collectionName)
}

func (m *BaseModel) Validate(validateMeta bool) *service_errors.Error {
	return nil
}
