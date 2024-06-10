package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppRoute struct {
	ID               primitive.ObjectID  `json:"_id" bson:"_id"`
	AppRoutesBlockID primitive.ObjectID  `json:"appRoutesBlockId" bson:"appRoutesBlockId"`
	AppRoutesBlock   AppRoutesBlock      `json:"appRoutesBlock,omitempty" bson:"-"`
	Name             string              `json:"name" bson:"name"`
	Description      *string             `json:"description" bson:"description"`
	Method           string              `json:"method" bson:"method"`
	Route            string              `json:"route" bson:"route"`
	Open             uint                `json:"open" bson:"open"`
	Active           bool                `json:"active" bson:"active"`
	CreatedBy        primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt        time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy        primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt        time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy        *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt        *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type AppRoutes struct {
	Count int64       `json:"count"`
	Rows  *[]AppRoute `json:"rows"`
}
