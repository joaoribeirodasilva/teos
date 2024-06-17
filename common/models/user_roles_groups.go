package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserRolesGroup = "user_roles_groups"
)

type UserRolesGroupModel struct {
	BaseModel
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	UserRoleID  primitive.ObjectID  `json:"userRoleId" bson:"roleId"`
	UserRole    UserRoleModel       `json:"userRole,omitempty" bson:"-"`
	UserGroupID primitive.ObjectID  `json:"userGroupId" bson:"userGroupId"`
	UserGroup   UserGroupModel      `json:"userGroup,omitempty" bson:"-"`
	Active      bool                `json:"active" bson:"active"`
	CreatedBy   primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy   primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy   *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt   *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *UserRolesGroupModel) GetCollectionName() string {
	return collectionUserRolesGroup
}

func (m *UserRolesGroupModel) AssignValues(to interface{}) error {

	dest, ok := to.(*UserRolesGroupModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.UserRoleID = m.UserRoleID
	dest.UserGroupID = m.UserGroupID
	dest.Active = m.Active
	to = dest

	return nil
}

func (m *UserRolesGroupModel) Validate() *logger.HttpError {

	validate := validator.New()

	// TODO: Validate related
	if err := validate.Var(m.Active, "required"); err != nil {
		fields := []string{"active"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid active ", err, nil)
	}

	/*
		userRole := NewUserRoleModel(m.ctx)
		if appErr := m.FindByID(m.UserRoleID, userRole); appErr != nil {
			return appErr
		}

		userGroup := NewUserGroupModel(m.ctx)
		if appErr := m.FindByID(m.UserGroupID, userGroup); appErr != nil {
			return appErr
		}
	*/
	return nil
}
