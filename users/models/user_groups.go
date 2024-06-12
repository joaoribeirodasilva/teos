package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserGroup struct {
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	ParentID    primitive.ObjectID  `json:"parentId" bson:"parentId"`
	Parent      *UserGroup          `json:"parent,omitempty"`
	Name        string              `json:"name,omitempty" `
	Description *string             `json:"description" bson:"description"`
	GroupPath   *string             `json:"userGroupPath" bson:"userGroupPath"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserGroups struct {
	Count int64        `json:"count"`
	Rows  *[]UserGroup `json:"rows"`
}
