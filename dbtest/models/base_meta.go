package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetaModel struct {
	BaseModel `json:"-" bson:"-"`
	ID        primitive.ObjectID  `json:"_id" bson:"_id"`
	CreatedBy primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt *time.Time          `json:"deletedAt" bson:"deletedAt"`
}
