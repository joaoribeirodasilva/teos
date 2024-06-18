package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResetModel struct {
	ID              primitive.ObjectID  `json:"_id" bson:"_id"`
	UserResetTypeID primitive.ObjectID  `json:"userResetTypeId" bson:"resetTypeId"`
	UserResetType   UserResetTypeModel  `json:"userResetType,omitempty" bson:"-"`
	UserUserID      primitive.ObjectID  `json:"userUserId" bson:"userId"`
	UserUser        UserUserModel       `json:"userUser,omitempty" bson:"-"`
	ResetKey        string              `json:"-" bson:"resetKey"`
	Used            *time.Time          `json:"used" bson:"used"`
	Expire          time.Time           `json:"expire" bson:"expire"`
	CreatedBy       primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt       time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy       primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt       time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy       *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt       *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserResetsModel struct {
	Count int64             `json:"count"`
	Docs  *[]UserResetModel `json:"docs"`
}

func (m *UserResetModel) GetID() primitive.ObjectID {
	return m.ID
}
