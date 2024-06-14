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
	collectionUserRole = "user_role"
	locationUserRole   = "COMMON::MODELS::UserRole"
)

type UserRoleModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	GroupPath   *string             `json:"userGroupPath" bson:"userGroupPath"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserRolesModel struct {
	BaseModels
	Count int64            `json:"count"`
	Rows  *[]UserRoleModel `json:"rows"`
}

func NewUserRoleModel(c *gin.Context) *UserRoleModel {
	m := &UserRoleModel{}
	m.Init(c, locationUserRole, collectionUserRole)
	return m
}

func NewUserRolesModel(c *gin.Context) *UserRolesModel {
	m := &UserRolesModel{}
	m.Init(c, locationUserRole, collectionUserRole)
	return m
}

func (m *UserRoleModel) FillMeta(create bool, delete bool) {

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

func (m *UserRoleModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *UserRoleModel) Validate() *service_errors.Error {

	validate := validator.New()
	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::UserRoleModel", "name", "invalid name")
	}

	return nil
}
