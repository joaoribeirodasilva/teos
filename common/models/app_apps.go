package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionAppAppsModel = "app_apps"
)

type AppAppModel struct {
	BaseModel   `json:"-" bson:"-"`
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	AppKey      string              `json:"appKey" bson:"appKey"`
	Active      bool                `json:"active" bson:"active"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *AppAppModel) GetCollectionName() string {
	return collectionAppAppsModel
}

func (m *AppAppModel) AssignValues(to interface{}) error {

	dest, ok := to.(*AppAppModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.Name = m.Name
	dest.Description = m.Description
	dest.AppKey = m.AppKey
	dest.Active = m.Active
	to = dest

	return nil
}

func (m *AppAppModel) Validate() *logger.HttpError {

	validate := validator.New()
	if err := validate.Var(m.Name, "required"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}

	if err := validate.Var(m.AppKey, "required,gte=1"); err != nil {
		fields := []string{"appKey"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid appKey ", err, nil)
	}

	return nil
}
