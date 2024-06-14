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
	collectionUserRolesGroupsUser = "user_roles_grups_user"
	locationUserRolesGroupsUser   = "COMMON::MODELS::UserRolesGroupsUser"
)

type UserRolesGroupsUserModel struct {
	BaseModel
	ID              primitive.ObjectID  `json:"_id" bson:"_id"`
	UserRoleGroupID primitive.ObjectID  `json:"userRoleGroupId" bson:"roleGroupId"`
	UserRoleGroup   UserRolesGroupModel `json:"userRoleGroup,omitempty" bson:"-"`
	UserUserID      primitive.ObjectID  `json:"userUserId" bson:"userId"`
	UserUser        UserUserModel       `json:"userUser,omitempty" bson:"-"`
	Active          bool                `json:"active" bson:"active"`
	CreatedBy       primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt       time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy       primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt       time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy       *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt       *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserRolesGroupsUsersModel struct {
	BaseModels
	Count int64                       `json:"count"`
	Rows  *[]UserRolesGroupsUserModel `json:"rows"`
}

func NewUserRolesGroupsUserModel(c *gin.Context) *UserRolesGroupsUserModel {
	m := &UserRolesGroupsUserModel{}
	m.Init(c, locationUserRolesGroupsUser, collectionUserRolesGroupsUser)
	return m
}

func NewUserRolesGroupsUsersModel(c *gin.Context) *UserRolesGroupsUsersModel {
	m := &UserRolesGroupsUsersModel{}
	m.Init(c, locationUserRolesGroupsUser, collectionUserRolesGroupsUser)
	return m
}

func (m *UserRolesGroupsUserModel) FillMeta(create bool, delete bool) {

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

func (m *UserRolesGroupsUserModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *UserRolesGroupsUserModel) Validate() *service_errors.Error {

	validate := validator.New()

	if err := validate.Var(m.Active, "required"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::UserRolesGroupModel", "avtive", "invalid active")
	}

	user := NewUserUserModel(m.ctx)
	if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
		return appErr
	}

	userRoleGroup := NewUserRolesGroupModel(m.ctx)
	if appErr := m.FindByID(m.UserRoleGroupID, userRoleGroup); appErr != nil {
		return appErr
	}

	return nil
}
