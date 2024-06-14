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
	collectionUserResetType = "user_reset_types"
	locationUserResetType   = "COMMON::MODELS::UserResetType"
)

type UserResetTypeModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserResetTypesModel struct {
	BaseModels
	Count int64             `json:"count"`
	Rows  *[]UserResetModel `json:"rows"`
}

func NewUserResetTypeModel(c *gin.Context) *UserResetTypeModel {
	m := &UserResetTypeModel{}
	m.Init(c, locationUserResetType, collectionUserResetType)
	return m
}

func NewUserResetTypesModel(c *gin.Context) *UserResetTypesModel {
	m := &UserResetTypesModel{}
	m.Init(c, locationUserResetType, collectionUserResetType)
	return m
}

func (m *UserResetTypeModel) FillMeta(create bool, delete bool) {

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

func (m *UserResetTypeModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *UserResetTypeModel) Validate() *service_errors.Error {

	validate := validator.New()
	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::UserResetTypeModel", "name", "invalid name")
	}
	return nil
}
