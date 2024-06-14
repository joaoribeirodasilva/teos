package models

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionAppAppsModel = "app_apps"
	locationAppAppsModel   = "COMMON::MODELS::AppAppModel"
)

type AppAppModel struct {
	BaseModel   `json:"-" bson:"-"`
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
}

type AppAppsModel struct {
	BaseModels
	Count int64
	Rows  []AppAppModel
}

func NewAppAppModel(c *gin.Context) *AppAppModel {
	m := &AppAppModel{}
	m.Init(c, locationAppAppsModel, collectionAppAppsModel)
	return m
}

func NewAppAppsModel(c *gin.Context) *AppAppsModel {
	m := &AppAppsModel{}
	m.Init(c, locationAppAppsModel, collectionAppAppsModel)
	return m
}

func (m *AppAppModel) FillMeta(create bool, delete bool) {

	now := time.Now().UTC()

	if create {
		m.ID = primitive.NewObjectID()
		m.CreatedBy = m.GetValues().Variables.User.ID
		m.CreatedAt = now
	} else if delete {
		m.DeletedBy = &m.GetValues().Variables.User.ID
		m.DeletedAt = &now
	}

	m.UpdatedBy = m.GetValues().Variables.User.ID
	m.UpdatedAt = now
}

func (m *AppAppModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *AppAppModel) Validate() *service_errors.Error {

	validate := validator.New()
	if err := validate.Var(m.Name, "required"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppApp", "name", "invalid application name")
	}

	if err := validate.Var(m.AppKey, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppApp", "appKey", "invalid application key")
	}

	return nil
}
