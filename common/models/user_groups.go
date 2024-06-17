package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserGroup = "user_groups"
)

type UserGroupModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	ParentID    *primitive.ObjectID `json:"parentId" bson:"parentId"`
	Parent      *UserGroupModel     `json:"parent,omitempty"`
	Name        string              `json:"name,omitempty" `
	Description *string             `json:"description" bson:"description"`
	GroupPath   *string             `json:"userGroupPath" bson:"userGroupPath"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *UserGroupModel) GetCollectionName() string {
	return collectionUserGroup
}

func (m *UserGroupModel) AssignValues(to interface{}) error {

	dest, ok := to.(*UserGroupModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.ParentID = m.ParentID
	dest.Name = m.Name
	dest.Description = m.Description
	dest.GroupPath = m.GroupPath
	to = dest

	return nil
}

func (m *UserGroupModel) Validate() *logger.HttpError {

	validate := validator.New()
	// TODO: Validate related
	/* 	if m.ParentID != nil {
		userGroup := NewUserGroupModel(m.ctx)
		if appErr := m.FindByID(*m.ParentID, userGroup); appErr != nil {
			return appErr
		}
	} */

	if err := validate.Var(m.Name, "required,gte=1"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}

	return nil
}
