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
	collectionUserRolesGroup = "user_roles_grups"
	locationUserRolesGroup   = "COMMON::MODELS::UserRolesGroup"
)

type UserRolesGroupModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	UserRoleID  primitive.ObjectID  `json:"userRoleId" bson:"roleId"`
	UserRole    UserRoleModel       `json:"userRole,omitempty" bson:"-"`
	UserGroupID primitive.ObjectID  `json:"userGroupId" bson:"userGroupId"`
	UserGroup   UserGroupModel      `json:"userGroup,omitempty" bson:"-"`
	Active      bool                `json:"active" bson:"active"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserRolesGroupsModel struct {
	BaseModels
	Count int64                  `json:"count"`
	Rows  *[]UserRolesGroupModel `json:"rows"`
}

func NewUserRolesGroupModel(c *gin.Context) *UserRolesGroupModel {
	m := &UserRolesGroupModel{}
	m.Init(c, locationUserRolesGroup, collectionUserRolesGroup)
	return m
}

func NewUserRolesGroupsModel(c *gin.Context) *UserRolesGroupsModel {
	m := &UserRolesGroupsModel{}
	m.Init(c, locationUserRolesGroup, collectionUserRolesGroup)
	return m
}

func (m *UserRolesGroupModel) FillMeta(create bool, delete bool) {

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

func (m *UserRolesGroupModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *UserRolesGroupModel) Validate() *service_errors.Error {

	validate := validator.New()

	if err := validate.Var(m.Active, "required"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::UserRolesGroupModel", "avtive", "invalid active")
	}

	userRole := NewUserRoleModel(m.ctx)
	if appErr := m.FindByID(m.UserRoleID, userRole); appErr != nil {
		return appErr
	}

	userGroup := NewUserGroupModel(m.ctx)
	if appErr := m.FindByID(m.UserGroupID, userGroup); appErr != nil {
		return appErr
	}

	return nil
}
