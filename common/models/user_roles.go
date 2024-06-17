package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserRole = "user_roles"
)

type UserRoleModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	GroupPath   *string             `json:"userGroupPath" bson:"userGroupPath"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *UserRoleModel) GetCollectionName() string {
	return collectionUserRole
}

func (m *UserRoleModel) AssignValues(to interface{}) error {

	dest, ok := to.(*UserRoleModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.Name = m.Name
	dest.Description = m.Description
	dest.GroupPath = m.GroupPath
	to = dest

	return nil
}

func (m *UserRoleModel) Validate() *logger.HttpError {

	validate := validator.New()
	if err := validate.Var(m.Name, "required"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}

	return nil
}
