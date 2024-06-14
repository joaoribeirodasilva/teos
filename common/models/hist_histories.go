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
	collectionHistHistory = "hist_history"
	locationHistHistory   = "COMMON::MODELS::HistHistory"
)

type HistHistoryModel struct {
	BaseModel
	ID         primitive.ObjectID  `json:"_id" bson:"_id"`
	AppAppID   primitive.ObjectID  `json:"appAppId" bson:"appAppId"`
	AppApp     interface{}         `json:"appApp,omitempty" bson:"-"`
	Collection string              `json:"collection" bson:"collection"`
	OriginalID primitive.ObjectID  `json:"originalId" bson:"originalId"`
	Data       interface{}         `json:"data" bson:"data"`
	CreatedBy  primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt  time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy  primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt  time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy  *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt  *time.Time          `json:"deletedAt" bson:"deletedAt"`
	NotFound   bool                `json:"-" bson:"-"`
}

type HistHistoriesModel struct {
	BaseModels
	Count int64               `json:"count"`
	Rows  *[]HistHistoryModel `json:"rows"`
}

func NewHistHistoryModel(c *gin.Context) *HistHistoryModel {
	m := &HistHistoryModel{}
	m.Init(c, locationAppAppsModel, collectionHistHistory)
	return m
}

func NewHistHistoriesModel(c *gin.Context) *HistHistoriesModel {
	m := &HistHistoriesModel{}
	m.Init(c, locationAppAppsModel, collectionAppAppsModel)
	return m
}

func (m *HistHistoryModel) FillMeta(create bool, delete bool) {

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

func (m *HistHistoryModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *HistHistoryModel) Validate() *service_errors.Error {

	validate := validator.New()

	appApp := NewAppAppModel(m.ctx)
	if appErr := m.FindByID(m.AppAppID, appApp); appErr != nil {
		return appErr
	}
	if err := validate.Var(m.Collection, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::HistHistoryModel", "collection", "invalid collection")
	}

	return nil
}
