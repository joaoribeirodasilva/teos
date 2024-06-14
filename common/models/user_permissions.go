package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserPermission = "user_permissions"
	locationUserPermission   = "COMMON::MODELS::UserPermission"
)

type UserPermissionModel struct {
	BaseModel
	ID              primitive.ObjectID  `json:"_id" bson:"_id"`
	AppRouteID      primitive.ObjectID  `json:"appRouteId" bson:"routeId"`
	AppRoute        AppRouteModel       `json:"appRoute,omitempty" bson:"-"`
	UserRoleGroupID primitive.ObjectID  `json:"userRoleGroupId" bson:"roleGroupId"`
	UserRoleGroup   UserRolesGroupModel `json:"userRoleGroup,omitempty" bson:"-"`
	UserUserID      primitive.ObjectID  `json:"userUserId" bson:"userId"`
	UserUser        UserUserModel       `json:"userUser,omitempty" bson:"-"`
	Active          bool                `json:"active" bson:"active" `
	CreatedBy       primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt       time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy       primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt       time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy       *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt       *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserPermissionsModel struct {
	BaseModels
	Count int64                  `json:"count"`
	Rows  *[]UserPermissionModel `json:"rows"`
}

func NewUserPermissionModel(c *gin.Context) *UserPermissionModel {
	m := &UserPermissionModel{}
	m.Init(c, locationUserPermission, collectionUserPermission)
	return m
}

func NewUserPermissionsModel(c *gin.Context) *UserPermissionsModel {
	m := &UserPermissionsModel{}
	m.Init(c, locationUserPermission, collectionUserPermission)
	return m
}

func (m *UserPermissionModel) FillMeta(create bool, delete bool) {

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

func (m *UserPermissionModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *UserPermissionModel) Validate() *service_errors.Error {

	//validate := validator.New()

	userAppRoute := NewAppRouteModel(m.ctx)
	if appErr := m.FindByID(m.AppRouteID, userAppRoute); appErr != nil {
		return appErr
	}

	userRoleGroup := NewUserRolesGroupModel(m.ctx)
	if appErr := m.FindByID(m.UserRoleGroupID, userRoleGroup); appErr != nil {
		return appErr
	}

	user := NewUserUserModel(m.ctx)
	if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
		return appErr
	}

	return nil
}
