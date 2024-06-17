package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionHistHistory = "hist_history"
)

type HistHistoryModel struct {
	BaseModel
	ID         primitive.ObjectID  `json:"_id" bson:"_id"`
	AppAppID   primitive.ObjectID  `json:"appAppId" bson:"appAppId"`
	AppApp     interface{}         `json:"appApp,omitempty" bson:"-"`
	Collection string              `json:"collection" bson:"collection"`
	OriginalID primitive.ObjectID  `json:"originalId" bson:"originalId"`
	Data       interface{}         `json:"data" bson:"data"`
	CreatedBy  primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt  time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy  primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt  time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy  *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt  *time.Time          `json:"deletedAt" bson:"deletedAt"`
	NotFound   bool                `json:"-" bson:"-"`
}

func (m *HistHistoryModel) GetCollectionName() string {
	return collectionHistHistory
}

func (m *HistHistoryModel) AssignValues(to interface{}) error {

	dest, ok := to.(*HistHistoryModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.AppAppID = m.AppAppID
	dest.Collection = m.Collection
	dest.OriginalID = m.OriginalID
	dest.Data = m.Data
	to = dest

	return nil
}

func (m *HistHistoryModel) Validate() *logger.HttpError {

	validate := validator.New()

	// TODO: Validate related
	/* 	appApp := NewAppAppModel(m.ctx)
	   	if appErr := m.FindByID(m.AppAppID, appApp); appErr != nil {
	   		return appErr
	   	} */

	if err := validate.Var(m.Collection, "required,gte=1"); err != nil {
		fields := []string{"collection"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid collection ", err, nil)
	}

	return nil
}
