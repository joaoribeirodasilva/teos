package models

import (
	"time"

	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserReset = "user_reset"
)

type UserResetModel struct {
	BaseModel
	ID              primitive.ObjectID  `json:"_id" bson:"_id"`
	UserResetTypeID primitive.ObjectID  `json:"userResetTypeId" bson:"resetTypeId"`
	UserResetType   UserResetTypeModel  `json:"userResetType,omitempty" bson:"-"`
	UserUserID      primitive.ObjectID  `json:"userUserId" bson:"userId"`
	UserUser        UserUserModel       `json:"userUser,omitempty" bson:"-"`
	ResetKey        string              `json:"-" bson:"resetKey"`
	Used            *time.Time          `json:"used" bson:"used"`
	Expire          time.Time           `json:"expire" bson:"expire"`
	CreatedBy       primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt       time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy       primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt       time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy       *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt       *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *UserResetModel) GetCollectionName() string {
	return collectionUserReset
}

func (m *UserResetModel) AssignValues(to interface{}) error {

	dest, ok := to.(*UserResetModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.UserResetTypeID = m.UserResetTypeID
	dest.UserUserID = m.UserUserID
	dest.UserUser = m.UserUser
	dest.ResetKey = m.ResetKey
	dest.Used = m.Used
	dest.Expire = m.Expire
	to = dest

	return nil
}

func (m *UserResetModel) Validate() *logger.HttpError {

	//validate := validator.New()
	// TODO: Validate related
	/* 	user := NewUserUserModel(m.ctx)
	   	if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
	   		return appErr
	   	}

	   	userResetType := NewUserResetTypeModel(m.ctx)
	   	if appErr := m.FindByID(m.UserResetTypeID, userResetType); appErr != nil {
	   		return appErr
	   	}
	*/

	return nil
}
