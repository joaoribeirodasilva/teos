package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppConfiguration struct {
	ID            primitive.ObjectID  `json:"_id" bson:"_id"`
	ApplicationID *primitive.ObjectID `json:"appAppId" bson:"appAppId"`
	Application   AppAppModel         `json:"appApp,omitempty" bson:"-"`
	Name          string              `json:"name" bson:"name"`
	Description   *string             `json:"description" bson:"description"`
	Key           string              `json:"key" bson:"key"`
	Type          string              `json:"type" bson:"type"`
	ValueInt      *int                `json:"valueInt" bson:"valueInt"`
	ValueString   *string             `json:"valueString" bson:"valueString"`
	ValueFloat    *float64            `json:"valueFloat" bson:"valueFloat"`
	ValueBool     *bool               `json:"valueBool" bson:"valueBool"`
	CreatedBy     primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt     time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy     primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt     time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy     *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt     *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type AppConfigurations struct {
	Count int64               `json:"count"`
	Rows  *[]AppConfiguration `json:"rows"`
}
