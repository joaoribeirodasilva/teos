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
	collectionAppRoutesBlock = "app_routes_blocks"
	locationAppRoutesBlock   = "COMMON::MODELS::AppRoutesBlock"
)

type AppRoutesBlockModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	AppAppID    primitive.ObjectID  `json:"appAppId" bson:"appAppId"`
	AppApp      AppAppModel         `json:"appApp,omitempty" bson:"-"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	Route       string              `json:"appRoute" bson:"appRoute"`
	Active      bool                `json:"active" bson:"active"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type AppRoutesBlocksModels struct {
	BaseModel
	Count int64                  `json:"count"`
	Rows  *[]AppRoutesBlockModel `json:"rows"`
}

func NewAppRoutesBlockModels(c *gin.Context) *AppRoutesBlockModel {
	m := &AppRoutesBlockModel{}
	m.Init(c, locationAppRoutesBlock, collectionAppRoutesBlock)
	return m
}

func NewAppRoutesBlocksModels(c *gin.Context) *AppRoutesBlocksModels {
	m := &AppRoutesBlocksModels{}
	m.Init(c, locationAppRoutesBlock, collectionAppRoutesBlock)
	return m
}

func (m *AppRoutesBlockModel) FillMeta(create bool, delete bool) {

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

func (m *AppRoutesBlockModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *AppRoutesBlockModel) Validate() *service_errors.Error {

	validate := validator.New()

	appApp := NewAppAppModel(m.ctx)
	if appErr := m.FindByID(m.AppAppID, appApp); appErr != nil {
		return appErr
	}

	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppRouteModel", "name", "invalid name")
	}

	if err := validate.Var(m.Route, "required"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppRouteModel", "route", "invalid route")
	}

	if err := validate.Var(m.Active, "required"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::AppRouteModel", "active", "invalid active")
	}

	return nil
}
