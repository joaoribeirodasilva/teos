package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRolesGroupModel struct {
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
	Count int64                  `json:"count"`
	Docs  *[]UserRolesGroupModel `json:"docs"`
}

func (m *UserRolesGroupModel) GetID() primitive.ObjectID {
	return m.ID
}
