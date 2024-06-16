package models

import (
	"sort"
	"time"

	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	ID             primitive.ObjectID  `json:"_id" bson:"_id" structs:"_id"`
	CreatedBy      primitive.ObjectID  `json:"createdBy" bson:"createdBy" structs:"createdBy"`
	CreatedAt      time.Time           `json:"createdAt" bson:"createdAt" structs:"createdAt"`
	UpdatedBy      primitive.ObjectID  `json:"updatedBy" bson:"updatedBy" structs:"updatedBy"`
	UpdatedAt      time.Time           `json:"updatedAt" bson:"updatedAt" structs:"updatedAt"`
	DeletedBy      *primitive.ObjectID `json:"deletedBy" bson:"deletedBy" structs:"deletedBy"`
	DeletedAt      *time.Time          `json:"deletedAt" bson:"deletedAt" structs:"deletedAt"`
	collectionName string
}

func (bm *BaseModel) GetCollectionName() string {
	return bm.collectionName
}

func (bm *BaseModel) GetID() *primitive.ObjectID {
	return &bm.ID
}

func (bm *BaseModel) SetID(id primitive.ObjectID) {
	bm.ID = id
}

func (bm *BaseModel) GetCreatedBy() primitive.ObjectID {
	return bm.CreatedBy
}

func (bm *BaseModel) SetCreatedBy(id primitive.ObjectID) {
	bm.CreatedBy = id
}

func (bm *BaseModel) GetCreatedAt() time.Time {
	return bm.CreatedAt
}

func (bm *BaseModel) SetCreatedAt(timestamp time.Time) {
	bm.CreatedAt = timestamp
}

func (bm *BaseModel) GetUpdateddBy() primitive.ObjectID {
	return bm.UpdatedBy
}

func (bm *BaseModel) SetUpdatedBy(id primitive.ObjectID) {
	bm.UpdatedBy = id
}

func (bm *BaseModel) GetUpdatedAt() time.Time {
	return bm.UpdatedAt
}

func (bm *BaseModel) SetUpdatedAt(timestamp time.Time) {
	bm.UpdatedAt = timestamp
}

func (bm *BaseModel) GetDeletedBy() *primitive.ObjectID {
	return bm.DeletedBy
}

func (bm *BaseModel) SetDeletedBy(id *primitive.ObjectID) {
	bm.DeletedBy = id
}

func (bm *BaseModel) GetDeleteddAt() *time.Time {
	return bm.DeletedAt
}

func (bm *BaseModel) SetDeletedAt(timestamp *time.Time) {
	bm.DeletedAt = timestamp
}

func (bm *BaseModel) Validate() error { return nil }

func (bm *BaseModel) AssignValues(to interface{}) error { return nil }

// func (bm *BaseModel) Normalize(to interface{}) error { return nil }
func (bm *BaseModel) Normalize(to interface{}) *bson.D {

	fromMap := structs.Map(bm)
	toMap := structs.Map(to)

	toMap["_id"] = fromMap["_id"]
	toMap["createdBy"] = fromMap["createdBy"]
	toMap["createdAt"] = fromMap["createdAt"]
	toMap["updatedBy"] = fromMap["updatedBy"]
	toMap["updatedAt"] = fromMap["updatedAt"]
	toMap["deletedBy"] = fromMap["deletedBy"]
	toMap["deletedAt"] = fromMap["deletedAt"]
	delete(toMap, "BaseModel")

	keys := make([]string, 0, len(toMap))
	for k := range toMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	d := bson.D{}
	for _, k := range keys {
		val := primitive.E{Key: k, Value: toMap[k]}
		d = append(d, val)
	}

	//dump.PrintJson(toMap)
	return &d
}
