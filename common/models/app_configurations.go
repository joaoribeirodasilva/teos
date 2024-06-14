package models

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	validTypes = []string{"string", "int", "float", "bool"}
)

const (
	collectionAppConfiguration = "app_configurations"
	locationAppConfiguration   = "COMMON::MODELS::AppConfiguration"
)

type AppConfigurationModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	AppAppID    *primitive.ObjectID `json:"appAppId" bson:"appAppId"`
	AppApp      AppAppModel         `json:"appApp,omitempty" bson:"-"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	Key         string              `json:"key" bson:"key"`
	Type        string              `json:"type" bson:"type"`
	ValueInt    *int                `json:"valueInt" bson:"valueInt"`
	ValueString *string             `json:"valueString" bson:"valueString"`
	ValueFloat  *float64            `json:"valueFloat" bson:"valueFloat"`
	ValueBool   *bool               `json:"valueBool" bson:"valueBool"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type AppConfigurationsModel struct {
	BaseModels
	Count int64                    `json:"count"`
	Rows  *[]AppConfigurationModel `json:"rows"`
}

func NewAppConfigurationModel(c *gin.Context) *AppConfigurationModel {
	m := &AppConfigurationModel{}
	m.Init(c, locationAppConfiguration, collectionAppConfiguration)
	return m
}

func NewBaseModels(c *gin.Context) *AppConfigurationsModel {
	m := &AppConfigurationsModel{}
	m.Init(c, locationAppConfiguration, collectionAppConfiguration)
	return m
}

func (m *AppConfigurationModel) FillMeta(create bool, delete bool) {

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

func (m *AppConfigurationModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *AppConfigurationModel) Validate() *service_errors.Error {

	validate := validator.New()

	if m.AppAppID != nil {
		appApp := NewAppAppModel(m.ctx)
		if appErr := m.FindByID(*m.AppAppID, appApp); appErr != nil {
			return appErr
		}
	}

	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppConfigurationModel", "name", "invalid name")
	}
	if err := validate.Var(m.Key, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppConfigurationModel", "key", "invalid key")
	}
	if err := validate.Var(m.Type, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppConfigurationModel", "type", "invalid type")
	}
	if !slices.Contains(validTypes, m.Type) {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppConfigurationModel", "type", "invalid type")
	}

	return nil
}
