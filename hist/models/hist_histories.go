package models

import (
	"time"

	apps_models "github.com/joaoribeirodasilva/teos/apps/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HistHistory struct {
	ID         primitive.ObjectID  `json:"_id" bson:"_id"`
	AppAppID   primitive.ObjectID  `json:"appAppId" bson:"appAppId"`
	AppApp     apps_models.AppApp  `json:"appApp,omitempty" bson:"-"`
	Table      string              `json:"table" bson:"table"`
	OriginalID primitive.ObjectID  `json:"originalId" bson:"originalId"`
	Data       interface{}         `json:"data" bson:"data"`
	CreatedBy  primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt  time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy  primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt  time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy  *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt  *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type HistHistories struct {
	Count int64          `json:"count"`
	Rows  *[]HistHistory `json:"rows"`
}
