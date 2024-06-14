package models

import (
	"time"

	"github.com/joaoribeirodasilva/teos/common/structures"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionHistHistory = "hist_history"
)

type HistHistory struct {
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
	values     *structures.RequestValues
}

type HistHistories struct {
	Count int64          `json:"count"`
	Rows  *[]HistHistory `json:"rows"`
}

// SetDb sstes the database connection pointer fo this model
func (m *HistHistory) SetValues(values *structures.RequestValues) {
	m.values = values
}

// GetCollection gets the collection name for the model
func (m *HistHistory) GetCollection() *mongo.Collection {
	return m.values.Variables.Db.Db.Collection(collectionHistHistory)
}
