package models

import (
	"time"

	apps_models "github.com/joaoribeirodasilva/teos/apps/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPermission struct {
	ID              primitive.ObjectID   `json:"_id" bson:"_id"`
	AppRouteID      primitive.ObjectID   `json:"appRouteId" bson:"routeId"`
	AppRoute        apps_models.AppRoute `json:"appRoute,omitempty" bson:"-"`
	UserRoleGroupID primitive.ObjectID   `json:"userRoleGroupId" bson:"roleGroupId"`
	UserRoleGroup   UserRolesGroup       `json:"userRoleGroup,omitempty" bson:"-"`
	UserUserID      primitive.ObjectID   `json:"userUserId" bson:"userId"`
	UserUser        UserUser             `json:"userUser,omitempty" bson:"-"`
	Active          bool                 `json:"active" bson:"active" `
	CreatedBy       primitive.ObjectID   `json:"createdBy" bson:"createdBy"`
	CreatedAt       time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedBy       primitive.ObjectID   `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt       time.Time            `json:"updatedAt" bson:"updatedAt"`
	DeletedBy       *primitive.ObjectID  `json:"deletedBy" bson:"deletedBy"`
	DeletedAt       *time.Time           `json:"deletedAt" bson:"deletedAt"`
}

type UserPermissions struct {
	Count int64             `json:"count"`
	Rows  *[]UserPermission `json:"rows"`
}
