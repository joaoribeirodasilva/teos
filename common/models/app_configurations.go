package models

import (
	"errors"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	validTypes = []string{"string", "int", "float", "bool"}
)

const (
	collectionAppConfiguration = "app_configurations"
)

type AppConfigurationModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	AppAppID    *primitive.ObjectID `json:"appAppId" bson:"appAppId"`
	AppApp      AppAppModel         `json:"appApp,omitempty" bson:"-"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	Key         string              `json:"key" bson:"key"`
	Type        string              `json:"type" bson:"type"`
	ValueInt    *int                `json:"valueInt" bson:"valueInt"`
	ValueString *string             `json:"valueString" bson:"valueString"`
	ValueFloat  *float64            `json:"valueFloat" bson:"valueFloat"`
	ValueBool   *bool               `json:"valueBool" bson:"valueBool"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *AppConfigurationModel) GetCollectionName() string {
	return collectionAppConfiguration
}

func (m *AppConfigurationModel) AssignValues(to interface{}) error {

	dest, ok := to.(*AppConfigurationModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.AppAppID = m.AppAppID
	dest.Name = m.Name
	dest.Description = m.Description
	dest.Key = m.Key
	dest.Type = m.Type
	dest.ValueInt = m.ValueInt
	dest.ValueString = m.ValueString
	dest.ValueFloat = m.ValueFloat
	dest.ValueBool = m.ValueBool
	to = dest

	return nil
}

func (m *AppConfigurationModel) Validate() *logger.HttpError {

	validate := validator.New()

	// TODO: Validate related
	/* 	if m.AppAppID != nil {
		appApp := NewAppAppModel(m.ctx)
		if appErr := m.FindByID(*m.AppAppID, appApp); appErr != nil {
			return appErr
		}
	} */

	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}
	if err := validate.Var(m.Key, "required,gte=1"); err != nil {
		fields := []string{"key"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid key ", err, nil)
	}
	if err := validate.Var(m.Type, "required,gte=1"); err != nil {
		fields := []string{"type"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid type ", err, nil)
	}
	if !slices.Contains(validTypes, m.Type) {
		fields := []string{"type"}
		err := errors.New("invalid field type")
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid type ", err, nil)
	}

	return nil
}
