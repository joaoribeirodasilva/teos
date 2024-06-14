package models

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"github.com/joaoribeirodasilva/teos/common/structures"
	models_history "github.com/joaoribeirodasilva/teos/hist/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionAppApps = "app_apps"
)

type AppApp struct {
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	AppKey      string              `json:"appKey" bson:"appKey"`
	Active      bool                `json:"active" bson:"active"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
	values      *structures.RequestValues
	NotFound    bool `json:"-" bson:"-"`
}

/*********************************************************
 * Single document
 *********************************************************/

// SetDb sstes the database connection pointer fo this model
func (m *AppApp) SetValues(values *structures.RequestValues) {
	m.values = values
}

func (m *AppApp) Bind(c *gin.Context) *service_errors.Error {

	// Binds the received JSON with the document
	if err := c.ShouldBindBodyWithJSON(m); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppApp::Bind", "", "failed to bind request. ERR: %s", err.Error())
	}

	// Validates the received document
	if appErr := m.Validate(false); appErr != nil {
		return appErr
	}
	return nil
}

// GetCollection gets the collection name for the model
func (m *AppApp) GetCollection() *mongo.Collection {
	return m.values.Variables.Db.Db.Collection(collectionAppApps)
}

// FindByID finds a single document based on it's unique id
func (m *AppApp) FindByID(id primitive.ObjectID, opts ...*options.FindOneOptions) *service_errors.Error {

	return m.First(bson.D{{Key: "_id", Value: id}}, opts...)
}

// First finds a single document based on the filter
func (m *AppApp) First(filter interface{}, opts ...*options.FindOneOptions) *service_errors.Error {

	if err := m.GetCollection().FindOne(context.TODO(), filter, opts...).Decode(m); err != nil {
		if err != mongo.ErrNoDocuments {
			m.NotFound = true
			return nil
		}
		return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApp::First", "", "failed to query database. ERR: %s", err.Error())
	}
	return nil
}

// Exists checks if a document exits based on the filter
func (m *AppApp) Exists(filter interface{}) *service_errors.Error {

	if appErr := m.First(filter); appErr != nil || m.NotFound {
		return appErr
	}
	return nil
}

// Create creates a new document on the database
func (m *AppApp) Create(opts ...*options.InsertOneOptions) *service_errors.Error {

	// Checks if uniques are violated
	existsDoc := AppApp{}
	if appErr := existsDoc.Exists(bson.D{{Key: "appKey", Value: m.AppKey}}); appErr != nil {
		return appErr
	}
	if !existsDoc.NotFound {
		return service_log.Error(0, http.StatusConflict, "MODELS::AppApps::Create", "appKey", "document already exists. id: %s", existsDoc.ID.Hex())
	}

	// Fill metadata fields
	if appErr := m.FillMeta(true, false); appErr != nil {
		return appErr
	}

	// Validate full document
	if appErr := m.Validate(true); appErr != nil {
		return appErr
	}

	// Insert document into the database
	if _, err := m.GetCollection().InsertOne(context.TODO(), m, opts...); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApp::Create", "", "failed to insert document into database. ERR: %s", err.Error())
	}

	return nil
}

func (m *AppApp) Update(id primitive.ObjectID, opts ...*options.FindOneOptions) *service_errors.Error {

	return nil
}

func (m *AppApp) Delete(id primitive.ObjectID, opts ...*options.FindOneOptions) *service_errors.Error {
	return nil
}

// FillMeta fills the metadata fileds according with the operation.
// If create and delete are both false then an update is assumed.
func (m *AppApp) FillMeta(create bool, delete bool) *service_errors.Error {

	if create && delete {
		return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApps::FillMeta", "", "a document can't be created and deleted at the same time")
	}

	now := time.Now().UTC()
	m.UpdatedBy = m.values.Variables.User.ID
	m.UpdatedAt = now
	if create {
		m.ID = primitive.NewObjectID()
		m.CreatedBy = m.values.Variables.User.ID
		m.CreatedAt = now
		return nil
	}
	m.DeletedBy = &m.values.Variables.User.ID
	m.DeletedAt = &now

	return nil

}

func (m *AppApp) SaveHistory(opts ...*options.InsertOneOptions) *service_errors.Error {

	// Creates the history document
	history := &models_history.HistHistory{
		ID:         primitive.NewObjectID(),
		AppAppID:   m.values.Variables.Configuration.ID,
		Collection: collectionAppApps,
		OriginalID: m.ID,
		Data:       m,
	}

	// Insert document into the database
	if _, err := history.GetCollection().InsertOne(context.TODO(), m, opts...); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApp::Create", "", "failed to insert document into database. ERR: %s", err.Error())
	}

	return nil
}

func (m *AppApp) Validate(validateMeta bool) *service_errors.Error {

	validate := validator.New()
	if err := validate.Var(m.Name, "required"); err != nil {
		return service_log.Error(0, http.StatusConflict, "MODELS::AppApp", "name", "invalid application name")
	}

	if err := validate.Var(m.AppKey, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusConflict, "MODELS::AppApp", "appKey", "invalid application key")
	}

	if validateMeta {
		if err := validate.Var(m.CreatedBy, "required,mongodbId"); err != nil {
			return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApp", "createdBy", "invalid created by")
		}

		if err := validate.Var(m.CreatedAt, "required"); err != nil {
			return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApp", "createdAt", "invalid created date")
		}

		if err := validate.Var(m.CreatedBy, "required,mongodbId"); err != nil {
			return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApp", "updatedBy", "invalid updated by")
		}

		if err := validate.Var(m.UpdatedAt, "required"); err != nil {
			return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApp", "updatedAt", "invalid updated at")
		}
	}

	return nil
}

/*********************************************************
 * Many documents
 *********************************************************/

type AppApps struct {
	Count  int64     `json:"count"`
	Rows   *[]AppApp `json:"rows"`
	values *structures.RequestValues
}

// SetDb sstes the database connection pointer fo this model
func (m *AppApps) SetValues(values *structures.RequestValues) {
	m.values = values
}

// GetCollection gets the collection name for the model
func (m *AppApps) GetCollection() *mongo.Collection {
	return m.values.Variables.Db.Db.Collection(collectionAppApps)
}

// FindAll return the documents from the database that match the filter
func (m *AppApps) FindAll(filter interface{}, opts ...*options.FindOptions) *service_errors.Error {

	if appErr := m.CountDocs(filter, nil); appErr != nil {
		return appErr
	}

	if m.Count == 0 {
		return nil
	}

	cursor, err := m.GetCollection().Find(context.TODO(), filter, opts...)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil
		}
		return service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApps", "", "failed to query database. ERR: %s", err.Error())
	}

	rows := make([]AppApp, 0)

	if err := cursor.All(context.TODO(), &rows); err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "CONTROLER::AppApps", "", "failed to fetch database cursor. ERR: %s", err.Error())
	}

	m.Rows = &rows

	return nil
}

// CountDocs counts the number of documents in the database that match the filter
func (m *AppApps) CountDocs(filter interface{}, opts ...*options.CountOptions) *service_errors.Error {

	var err error
	m.Count, err = m.GetCollection().CountDocuments(context.TODO(), filter)
	if err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "MODELS::AppApps", "", "failed to count records. ERR: %s", err.Error())
		return appErr
	}
	return nil
}
