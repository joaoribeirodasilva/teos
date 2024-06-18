package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HistHistoryModel struct {
	ID         primitive.ObjectID  `json:"_id" bson:"_id"`
	AppAppID   primitive.ObjectID  `json:"appAppId" bson:"appAppId"`
	AppApp     interface{}         `json:"appApp,omitempty" bson:"-"`
	Collection string              `json:"collection" bson:"collection"`
	OriginalID primitive.ObjectID  `json:"originalId" bson:"originalId"`
	Data       interface{}         `json:"data" bson:"data"`
	CreatedBy  primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt  time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy  primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt  time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy  *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt  *time.Time          `json:"deletedAt" bson:"deletedAt"`
	NotFound   bool                `json:"-" bson:"-"`
}

type HistHistoriesModel struct {
	Count int64               `json:"count"`
	Docs  *[]HistHistoryModel `json:"docs"`
}

func (m *HistHistoryModel) GetID() primitive.ObjectID {
	return m.ID
}
