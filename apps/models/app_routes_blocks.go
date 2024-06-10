package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppRoutesBlock struct {
	ID            primitive.ObjectID  `json:"_id" bson:"_id"`
	ApplicationID primitive.ObjectID  `json:"appAppId" bson:"appAppId"`
	Application   AppApp              `json:"appApp,omitempty" bson:"-"`
	Name          string              `json:"name" bson:"name"`
	Description   *string             `json:"description" bson:"description"`
	Route         string              `json:"appRoute" bson:"appRoute"`
	Active        bool                `json:"active" bson:"active"`
	CreatedBy     primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt     time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy     primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt     time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy     *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt     *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type AppRoutesBlocks struct {
	Count int64             `json:"count"`
	Rows  *[]AppRoutesBlock `json:"rows"`
}
