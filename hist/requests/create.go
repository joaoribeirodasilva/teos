package requests

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HistoriesCreate struct {
	AuthenticationKey string              `json:"authenticationKey"`
	Table             string              `json:"table"`
	OriginalID        primitive.ObjectID  `json:"originalId"`
	Data              interface{}         `json:"data"`
	CreatedBy         primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt         time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy         primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt         time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy         *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt         *time.Time          `json:"deletedAt" bson:"deletedAt"`
}
