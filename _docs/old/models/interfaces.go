package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type iBaseCollection interface {
	GetCollection() *mongo.Collection
	SetCollection(collection *mongo.Collection) *mongo.Collection
	GetCollectionName() string
	SetCollectionName(name string)
	GetDatabase() *mongo.Database
	SetDatabase(db *mongo.Database)
	SetUserID(user *primitive.ObjectID)
	GetUserID() *primitive.ObjectID
	AssignMeta(model iBaseModel, create bool, delete bool) error
	FindAll(filter interface{}, model iBaseModel, opts ...*options.FindOptions) (int64, error)
	Find(filter interface{}, model iBaseModel, opts ...*options.FindOptions) (int64, error)
	First(filter interface{}, model iBaseModel, opts ...*options.FindOneOptions) error
	FirstByID(id *primitive.ObjectID, model iBaseModel, opts ...*options.FindOneOptions) error
	Create(uniqueFilter interface{}, model iBaseModel, opts ...*options.InsertOneOptions) error
	Update(id *primitive.ObjectID, uniqueFilter interface{}, model iBaseModel, opts ...*options.UpdateOptions) error
	Delete(id *primitive.ObjectID, opts ...*options.DeleteOptions) error
	CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error)
}

type iBaseModel interface {
	GetID() *primitive.ObjectID
	SetID(id primitive.ObjectID)
	GetCreatedBy() primitive.ObjectID
	SetCreatedBy(id primitive.ObjectID)
	GetCreatedAt() time.Time
	SetCreatedAt(timestamp time.Time)
	GetUpdateddBy() primitive.ObjectID
	SetUpdatedBy(id primitive.ObjectID)
	GetUpdatedAt() time.Time
	SetUpdatedAt(timestamp time.Time)
	GetDeletedBy() *primitive.ObjectID
	SetDeletedBy(id *primitive.ObjectID)
	GetDeleteddAt() *time.Time
	SetDeletedAt(timestamp *time.Time)
	GetCollectionName() string
	Validate() error
	AssignValues(to interface{}) error
	Normalize(to interface{}) (*bson.D, error)
}
