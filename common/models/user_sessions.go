package models

import (
	"time"

	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserSession = "user_sessions"
)

type UserSessionModel struct {
	BaseModel
	ID         primitive.ObjectID  `json:"_id" bson:"_id"`
	UserUserID primitive.ObjectID  `json:"userUserId" bson:"userUserId"`
	UserUser   UserUserModel       `json:"userUser,omitempty" bson:"-"`
	CreatedBy  primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt  time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy  primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt  time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy  *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt  *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *UserSessionModel) GetCollectionName() string {
	return collectionUserSession
}

func (m *UserSessionModel) AssignValues(to interface{}) error {

	dest, ok := to.(*UserSessionModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.UserUserID = m.UserUserID
	to = dest

	return nil
}

func (m *UserSessionModel) Validate() *logger.HttpError {

	// TODO: Validate related
	/* 	user := NewUserUserModel(m.ctx)
	   	if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
	   		return appErr
	   	} */

	return nil
}
