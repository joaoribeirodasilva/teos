package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserResetType = "user_reset_types"
)

type UserResetTypeModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *UserResetTypeModel) GetCollectionName() string {
	return collectionUserResetType
}

func (m *UserResetTypeModel) AssignValues(to interface{}) error {

	dest, ok := to.(*UserResetTypeModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.Name = m.Name
	dest.Description = m.Description
	to = dest

	return nil
}

func (m *UserResetTypeModel) Validate() *logger.HttpError {

	validate := validator.New()

	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}

	return nil
}
