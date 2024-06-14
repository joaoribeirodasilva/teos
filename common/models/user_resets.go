package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserReset = "user_reset"
	locationUserReset   = "COMMON::MODELS::UserReset"
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

type UserResetsModel struct {
	BaseModels
	Count int64             `json:"count"`
	Rows  *[]UserResetModel `json:"rows"`
}

func NewUserResetModel(c *gin.Context) *UserResetModel {
	m := &UserResetModel{}
	m.Init(c, locationUserReset, collectionUserReset)
	return m
}

func NewUserResetsModel(c *gin.Context) *UserResetsModel {
	m := &UserResetsModel{}
	m.Init(c, locationUserReset, collectionUserReset)
	return m
}

func (m *UserResetModel) FillMeta(create bool, delete bool) {

	now := time.Now().UTC()

	if create {
		m.ID = primitive.NewObjectID()
		m.CreatedBy = m.GetValues().Variables.User.ID
		m.CreatedAt = now
	} else if delete {
		m.DeletedBy = &m.GetValues().Variables.User.ID
		m.DeletedAt = &now
	}

	m.UpdatedBy = m.GetValues().Variables.User.ID
	m.UpdatedAt = now
}

func (m *UserResetModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *UserResetModel) Validate() *service_errors.Error {

	//validate := validator.New()

	user := NewUserUserModel(m.ctx)
	if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
		return appErr
	}

	userResetType := NewUserResetTypeModel(m.ctx)
	if appErr := m.FindByID(m.UserResetTypeID, userResetType); appErr != nil {
		return appErr
	}

	return nil
}
