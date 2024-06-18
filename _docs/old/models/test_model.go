package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	defaultTestCollectionName = "dbtest"
)

type TestCollection struct {
	BaseCollection
}

type TestModel struct {
	BaseModel `json:"-" bson:"-"`
	ID        primitive.ObjectID  `json:"_id" bson:"_id" structs:"_id"`
	Name      string              `json:"name" bson:"name" structs:"name"`
	Age       int                 `json:"age" bson:"age" structs:"age"`
	CreatedBy primitive.ObjectID  `json:"createdBy" bson:"createdBy" structs:"createdBy"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt" structs:"createdAt"`
	UpdatedBy primitive.ObjectID  `json:"updatedBy" bson:"updatedBy" structs:"updatedBy"`
	UpdatedAt time.Time           `json:"updatedAt" bson:"updatedAt" structs:"updatedAt"`
	DeletedBy *primitive.ObjectID `json:"deletedBy" bson:"deletedBy" structs:"deletedBy"`
	DeletedAt *time.Time          `json:"deletedAt" bson:"deletedAt" structs:"deletedAt"`
}

func (m *TestModel) GetCollectionName() string {
	return defaultTestCollectionName
}

func (m *TestModel) Validate() error {

	if m.Name == "" {
		return errors.New("name is required")
	}

	if m.Age <= 0 {
		return errors.New("age must be greater than 0")
	}

	return nil
}

func (m *TestModel) AssignValues(to interface{}) error {

	dest, ok := to.(*TestModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.Name = m.Name
	dest.Age = m.Age
	to = dest

	return nil
}
