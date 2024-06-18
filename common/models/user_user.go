package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUserModel struct {
	ID            primitive.ObjectID  `json:"_id" bson:"_id"`
	FirstName     string              `json:"firstName" bson:"firstName"`
	Surname       string              `json:"surname" bson:"surname"`
	Email         string              `json:"email" bson:"email"`
	Password      *string             `json:"password,omitempty" bson:"password"`
	Terms         *time.Time          `json:"terms" bson:"terms"`
	AvatarUrl     string              `json:"avatarUrl" bson:"avatarUrl"`
	EmailVerified *time.Time          `json:"emailVerified" bson:"emailVerified"`
	Active        uint                `json:"active" bson:"active"`
	CreatedBy     primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt     time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy     primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt     time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy     *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt     *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserUsersModel struct {
	Count int64            `json:"count"`
	Docs  *[]UserUserModel `json:"docs"`
}

func (m *UserUserModel) GetID() primitive.ObjectID {
	return m.ID
}
