package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserSession = "user_sessions"
	locationUserSession   = "COMMON::MODELS::UserSession"
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

type UserSessionsModel struct {
	BaseModels
	Count int64               `json:"count"`
	Rows  *[]UserSessionModel `json:"rows"`
}

func NewUserSessionModel(c *gin.Context) *UserSessionModel {
	m := &UserSessionModel{}
	m.Init(c, locationUserSession, collectionUserSession)
	return m
}

func NewUserSessionsModel(c *gin.Context) *UserSessionsModel {
	m := &UserSessionsModel{}
	m.Init(c, locationUserSession, collectionUserSession)
	return m
}

func (m *UserSessionModel) FillMeta(create bool, delete bool) {

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

func (m *UserSessionModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *UserSessionModel) Validate() *service_errors.Error {

	user := NewUserUserModel(m.ctx)
	if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
		return appErr
	}

	return nil
}
