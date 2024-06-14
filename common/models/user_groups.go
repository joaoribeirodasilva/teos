package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserGroup = "user_groups"
	locationUserGroup   = "COMMON::MODELS::UserGroup"
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

type UserGroupsModel struct {
	BaseModels
	Count int64             `json:"count"`
	Rows  *[]UserGroupModel `json:"rows"`
}

func NewUserGroupModel(c *gin.Context) *UserGroupModel {
	m := &UserGroupModel{}
	m.Init(c, locationUserGroup, collectionUserGroup)
	return m
}

func NewUserGroupsModel(c *gin.Context) *UserGroupsModel {
	m := &UserGroupsModel{}
	m.Init(c, locationUserGroup, collectionUserGroup)
	return m
}

func (m *UserGroupModel) FillMeta(create bool, delete bool) {

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

func (m *UserGroupModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *UserGroupModel) Validate() *service_errors.Error {

	//validate := validator.New()

	if m.ParentID != nil {
		userGroup := NewUserGroupModel(m.ctx)
		if appErr := m.FindByID(*m.ParentID, userGroup); appErr != nil {
			return appErr
		}
	}

	return nil
}
