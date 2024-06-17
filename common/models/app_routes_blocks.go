package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionAppRoutesBlock = "app_routes_blocks"
)

type AppRoutesBlockModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	AppAppID    primitive.ObjectID  `json:"appAppId" bson:"appAppId"`
	AppApp      AppAppModel         `json:"appApp,omitempty" bson:"-"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	Route       string              `json:"appRoute" bson:"appRoute"`
	Active      bool                `json:"active" bson:"active"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *AppRoutesBlockModel) GetCollectionName() string {
	return collectionAppRoutesBlock
}

func (m *AppRoutesBlockModel) AssignValues(to interface{}) error {

	dest, ok := to.(*AppRoutesBlockModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.AppAppID = m.AppAppID
	dest.Name = m.Name
	dest.Description = m.Description
	dest.Route = m.Route
	dest.Active = m.Active
	to = dest

	return nil
}

func (m *AppRoutesBlockModel) Validate() *logger.HttpError {

	validate := validator.New()

	// TODO: Validate related
	/* 	appApp := NewAppAppModel(m.ctx)
	   	if appErr := m.FindByID(m.AppAppID, appApp); appErr != nil {
	   		return appErr
	   	} */

	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}

	if err := validate.Var(m.Route, "required"); err != nil {
		fields := []string{"route"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid route ", err, nil)
	}

	if err := validate.Var(m.Active, "required"); err != nil {
		fields := []string{"active"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid active ", err, nil)
	}

	return nil
}
