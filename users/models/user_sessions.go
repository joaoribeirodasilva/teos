package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSession struct {
	ID         primitive.ObjectID  `json:"_id" bson:"_id"`
	UserUserID primitive.ObjectID  `json:"userUserId" bson:"userUserId"`
	UserUser   UserUser            `json:"userUser,omitempty" bson:"-"`
	CreatedBy  primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt  time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy  primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt  time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy  *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt  *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserSessions struct {
	Count int64          `json:"count"`
	Rows  *[]UserSession `json:"rows"`
}
