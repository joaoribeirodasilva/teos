package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole struct {
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

type UserRoles struct {
	Count int64       `json:"count"`
	Rows  *[]UserRole `json:"rows"`
}
