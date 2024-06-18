package models

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/joaoribeirodasilva/teos/common/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseCollectionOptions struct {
	Ctx        context.Context
	UseMetas   bool
	UseUserID  bool
	UseDates   bool
	SoftDelete bool
	Debug      bool
}

type BaseCollection struct {
	CollectionName string
	Collection     *mongo.Collection
	Database       *database.Db
	UserID         *primitive.ObjectID
	Cursor         *mongo.Cursor
	options        *BaseCollectionOptions
}

// mongo base errors and original mongo errors mapped as mongo base errors
var (
	ErrInvalidObjectId       = errors.New("mongobase: invalid object id")
	ErrDocumentExists        = errors.New("mongobase: document exists")
	ErrModelIsNotPointer     = errors.New("mongobase: model is not a pointer")
	ErrWrongModelType        = errors.New("mongobase: model is not of the same type in variable assignment")
	ErrNoUserId              = errors.New("mongobase: no user id defined")
	ErrClientDisconnected    = mongo.ErrClientDisconnected
	ErrEmptySlice            = mongo.ErrEmptySlice
	ErrInvalidIndexValue     = mongo.ErrInvalidIndexValue
	ErrMapForOrderedArgument = mongo.ErrMapForOrderedArgument{}
	ErrMissingResumeToken    = mongo.ErrMissingResumeToken
	ErrMultipleIndexDrop     = mongo.ErrMultipleIndexDrop
	ErrNilCursor             = mongo.ErrNilCursor
	ErrNilDocument           = mongo.ErrNilDocument
	ErrNoDocuments           = mongo.ErrNoDocuments
	ErrNonStringIndexName    = mongo.ErrNonStringIndexName
	ErrUnacknowledgedWrite   = mongo.ErrUnacknowledgedWrite
	ErrWrongClient           = mongo.ErrWrongClient
)

// NewBaseCollection initializes a new collection to handle queries in it.
//
// The 'name' is the name of the collection to be queried and not be empty.
// The 'db' must be a valid pointer to the mongo Collection structure and not be null.
// The 'ctx' parameter must be a valid context to be applied to all queries that will be ran in this collection. Thos parameter if null defaults to context.TODO()
func NewBaseCollection(name string, db *database.Db, opts *BaseCollectionOptions) *BaseCollection {

	bc := new(BaseCollection)
	bc.CollectionName = name
	bc.Database = db
	bc.Collection = bc.Database.GetDatabase().Collection(bc.CollectionName)

	bc.options = opts
	if opts == nil {
		bc.options = &BaseCollectionOptions{
			Ctx:        bc.Database.GetContext(),
			UseMetas:   true,
			UseUserID:  true,
			UseDates:   true,
			SoftDelete: true,
			Debug:      true,
		}
	}
	if !bc.options.UseMetas {
		bc.options.UseDates = false
		bc.options.UseUserID = false
	}
	return bc
}

// FindAll finds documents in this collection according to the parameters defined in the 'filter' parameter.
// Internally it calls the function 'Find'. After that it returns the number of documents found in the collection that match the filter
// and fills the model parameter with the documents found.
//
// The 'filter' parameter must be a valid mongo filter containing the query to perform.
// The 'model' parameter must be a pointer to an array of 'iBaseModel' structures of and never be null.
// The 'opts' parameter is optional and receives a pointer to a mongo 'options.FindOptions' type.
func (bc *BaseCollection) FindAll(filter interface{}, model interface{}, opts ...*options.FindOptions) (int64, error) {

	count, err := bc.Find(filter, model, opts...)
	if err != nil {
		return 0, err
	}

	if err := bc.Cursor.All(bc.options.Ctx, model); err != nil {
		return 0, err
	}

	return count, nil
}

// Find finds documents in this collection according to the parameters defined in the 'filter' parameter.
// After that it returns the number of documents found in the collection that match the filter
// and sets the 'BaseCollection.Cursor' variable with the database document cursor.
//
// The 'filter' parameter must be a valid mongo filter containing the query to perform.
// The 'opts' parameter is optional and receives a pointer to a mongo 'options.FindOptions' type.
func (bc *BaseCollection) Find(filter interface{}, model interface{}, opts ...*options.FindOptions) (int64, error) {

	bc.Cursor = nil

	count, err := bc.CountDocuments(filter, nil)
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, nil
	}

	cursor, err := bc.Collection.Find(bc.options.Ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	bc.Cursor = cursor

	return count, nil
}

// First fills the model pointer passed in the 'model' parameter with the first occurrence of the document found
// that matches the filter in the parameter 'filter' returning an error in any error situation.
//
// The 'filter' parameter must be a valid mongo filter containing the query to perform.
// The 'model' parameter must be a pointer to an array of 'iBaseModel' structures of and never be null.
// The 'opts' parameter is optional and receives a pointer to a mongo 'options.FindOneOptions' type.
func (bc *BaseCollection) First(filter interface{}, model iBaseModel, opts ...*options.FindOneOptions) error {

	if !isPtr(model) {
		return ErrModelIsNotPointer
	}

	if err := bc.Collection.FindOne(bc.options.Ctx, filter, opts...).Decode(model); err != nil {
		return err
	}

	return nil
}

// FirstByID fills the model pointer passed in the 'model' parameter with the document that has
// the id specified in the 'id' parameter, returning an error in any error situation.
// Internally it calls the 'BaseCollection.First' with the proper id filter.
//
// The 'id' parameter must be a valid 'mongo.objectID' pointer and not be null. If this pointer is null an error is returned
// The 'model' parameter must be a pointer to an array of 'iBaseModel' structures of and never be null.
// The 'opts' parameter is optional and receives a pointer to a mongo 'options.FindOneOptions' type.
func (bc *BaseCollection) FirstByID(id *primitive.ObjectID, model iBaseModel, opts ...*options.FindOneOptions) error {

	if id == nil {
		return ErrInvalidObjectId
	}

	if err := bc.First(bson.D{{Key: "_id", Value: bson.D{{Key: "$eq", Value: id}}}}, model, opts...); err != nil {
		return err
	}

	return nil
}

// Create creates a new mongo document in the database and checks for uniqueness based on the 'uniqueFilter' parameter. If
// a unique filter is not defined then uniqueness is not verified and any document can be freely inserted in the collection.
// This function also assign the 'model' metadata values incliding the new document _id.
//
// The 'uniqueFilter' parameter must be a valid mongo filter containing the query to find possible document conflicts at the
// unique field values or null if no unique filter is to be used.
// The 'model' parameter must be a pointer to an array of 'iBaseModel' structures of and never be null.
// The 'opts' parameter is optional and receives a pointer to a mongo 'options.InsertOneOptions' type.
func (bc *BaseCollection) Create(uniqueFilter interface{}, model iBaseModel, opts ...*options.InsertOneOptions) error {

	if !isPtr(model) {
		return ErrModelIsNotPointer
	}

	if uniqueFilter != nil {
		exists := BaseModel{}
		if err := bc.First(uniqueFilter, &exists); err != nil {
			if err != mongo.ErrNoDocuments {
				return err
			}
		} else {
			return ErrDocumentExists
		}
	}

	if err := bc.AssignMeta(model, true, false); err != nil {
		return err
	}

	if err := model.AssignValues(model); err != nil {
		return err
	}

	doc, err := model.Normalize(model)
	if err != nil {
		return err
	}

	if _, err := bc.Collection.InsertOne(bc.options.Ctx, doc, opts...); err != nil {
		return err
	}

	return nil
}

// Update updates an existing document in the database and checks for uniqueness based on the 'uniqueFilter' parameter. If
// a unique filter is not defined then uniqueness is not verified and any document can be freely updated in the collection.
// If the document being updated doesn't exist it return an ErrNoDocuments error.
// This function also assign the 'model' metadata values.
//
// The 'uniqueFilter' parameter must be a valid mongo filter containing the query to find possible document conflicts at the
// unique field values or null if no unique filter is to be used.
// The 'model' parameter must be a pointer to an array of 'iBaseModel' structures of and never be null.
// The 'opts' parameter is optional and receives a pointer to a mongo 'options.UpdateOptions' type.
func (bc *BaseCollection) Update(id *primitive.ObjectID, uniqueFilter interface{}, model iBaseModel, opts ...*options.UpdateOptions) error {

	if id == nil {
		return ErrInvalidObjectId
	}

	if !isPtr(model) {
		return ErrModelIsNotPointer
	}

	if uniqueFilter != nil {
		duplicated := BaseModel{}
		if err := bc.First(uniqueFilter, &duplicated); err != nil {
			if err != mongo.ErrNoDocuments {
				return err
			}
		} else if *id != duplicated.ID {
			return ErrDocumentExists
		}
	}

	original := cloneInterfacePtr(model)

	if err := bc.First(uniqueFilter, original); err != nil {
		return err
	}

	if err := model.AssignValues(original); err != nil {
		return err
	}

	if err := bc.AssignMeta(original, false, false); err != nil {
		return err
	}

	doc, err := model.Normalize(model)
	if err != nil {
		return err
	}

	if _, err := bc.Collection.UpdateOne(bc.options.Ctx, bson.D{{Key: "_id", Value: id}}, doc); err != nil {
		return err
	}

	return nil
}

// Delete deleted an existing document in the database based on the 'id' parameter. It tries to
// find the document first. If the document being deleted doesn't exist it return an ErrNoDocuments error.
func (bc *BaseCollection) Delete(id *primitive.ObjectID, opts ...*options.DeleteOptions) error {

	if id == nil {
		return ErrInvalidObjectId
	}

	exists := &BaseModel{}
	if err := bc.FirstByID(id, exists); err != nil {
		return err
	}

	if bc.options.SoftDelete {
		if _, err := bc.Collection.UpdateOne(bc.options.Ctx, bson.D{{Key: "_id", Value: exists.ID}}, bson.D{{Key: "$set", Value: exists}}); err != nil {
			return err
		}
	} else {
		if _, err := bc.Collection.DeleteOne(bc.options.Ctx, bson.D{{Key: "_id", Value: id}}, opts...); err != nil {
			return err
		}
	}

	return nil
}

// CountDocuments counts the number of documents in a collection that match the filter in the 'filter' parameter.
//
// The 'filter' parameter must be a valid mongo filter containing the query to perform.
// The 'opts' parameter is optional and receives a pointer to a mongo 'options.CountOptions' type.
func (bc *BaseCollection) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error) {

	count, err := bc.Collection.CountDocuments(bc.options.Ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Assigns metadata filed models according to the operation is course if metas are enabled on options.
// It will only assign those meta properties that are enabled. Although it always assign the 'ID' property
// of the document being created.
func (bc *BaseCollection) AssignMeta(model iBaseModel, create bool, delete bool) error {

	now := time.Now().UTC()
	if bc.options.UseUserID {
		if bc.UserID == nil {
			return ErrNoUserId
		}
		model.SetUpdatedBy(*bc.UserID)
	}
	if bc.options.UseDates {
		model.SetUpdatedAt(now)
	}
	if create {
		id := primitive.NewObjectID()
		model.SetID(id)
		if bc.options.UseUserID {
			model.SetCreatedBy(*bc.UserID)
		}
		if bc.options.UseDates {
			model.SetCreatedAt(now)
		}
	} else if delete {
		if bc.options.UseUserID {
			model.SetDeletedBy(bc.UserID)
		}
		if bc.options.UseDates {
			model.SetDeletedAt(&now)
		}
	}

	return nil
}

func (bc *BaseCollection) GetCollection() *mongo.Collection {
	return bc.Collection
}

// SetCollection sets a new collection pointer into 'BaseCollection.Collection'.
//
// If the pointer is null all references to the collection
// will be cleared including the 'BaseCollection.CollectionName' string.
func (bc *BaseCollection) SetCollection(collection *mongo.Collection) {
	if collection != nil {
		bc.CollectionName = bc.Collection.Name()
		return
	}
	bc.CollectionName = ""
	bc.Collection = nil
}

// GetCollectionName returns currently set collection
// name for this collection.
func (bc *BaseCollection) GetCollectionName() string {
	return bc.CollectionName
}

// SetCollectionName sets a new collection for this structure.
//
// By setting a new collection mae a new 'mongo.Collection' pointer
// will also be set automatically, no need to manually set
// a new 'BaseCollection.Collection' through 'BaseCollection.SetCollection'.
// If the collection name is an empty string all references to the collection
// will be cleared including the 'BaseCollection.Collection' pointer.
func (bc *BaseCollection) SetCollectionName(name string) {

	bc.CollectionName = name
	if name == "" {
		bc.Collection = nil
		return
	}
	bc.Collection = bc.Database.GetDatabase().Collection(bc.CollectionName)
}

// GetDatabase return the current 'mongo.Database' set for this collection.
func (bc *BaseCollection) GetDatabase() *database.Db {
	return bc.Database
}

// SetDatabase sets the 'BaseCollection.Database' pointer
//
// The 'db' parameter can't be null, otherwise errors will be thrown everywhere.
// When set the current collection name will be also set for the new database in the
// 'db' parameter.
func (bc *BaseCollection) SetDatabase(db *database.Db) {
	bc.Database = db
	bc.Collection = bc.Database.GetDatabase().Collection(bc.CollectionName)
}

// GetCursor returns the last Find call 'BaseCollection.Cursor' instance.
func (bc *BaseCollection) GetCursor() *mongo.Cursor {
	return bc.Cursor
}

func (bc *BaseCollection) SetUserID(user *primitive.ObjectID) {
	bc.UserID = user
}

func (bc *BaseCollection) GetUserID() *primitive.ObjectID {
	return bc.UserID
}

func isPtr(ptr interface{}) bool {
	return reflect.ValueOf(ptr).Type().Kind() == reflect.Ptr
}

func cloneInterfacePtr(st interface{}) iBaseModel {
	return reflect.New(reflect.ValueOf(st).Elem().Type()).Interface().(iBaseModel)
}
