package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPermissionModel struct {
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
	Count int64                  `json:"count"`
	Docs  *[]UserPermissionModel `json:"docs"`
}

func (m *UserPermissionModel) GetID() primitive.ObjectID {
	return m.ID
}
