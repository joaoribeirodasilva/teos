package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserPermission = "user_permissions"
)

type UserPermissionModel struct {
	BaseModel
	ID              primitive.ObjectID  `json:"_id" bson:"_id"`
	AppRouteID      primitive.ObjectID  `json:"appRouteId" bson:"routeId"`
	AppRoute        AppRouteModel       `json:"appRoute,omitempty" bson:"-"`
	UserRoleGroupID primitive.ObjectID  `json:"userRoleGroupId" bson:"roleGroupId"`
	UserRoleGroup   UserRolesGroupModel `json:"userRoleGroup,omitempty" bson:"-"`
	UserUserID      primitive.ObjectID  `json:"userUserId" bson:"userId"`
	UserUser        UserUserModel       `json:"userUser,omitempty" bson:"-"`
	Active          bool                `json:"active" bson:"active" `
	CreatedBy       primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt       time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy       primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt       time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy       *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt       *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func (m *UserPermissionModel) GetCollectionName() string {
	return collectionUserPermission
}

func (m *UserPermissionModel) AssignValues(to interface{}) error {

	dest, ok := to.(*UserPermissionModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.AppRouteID = m.AppRouteID
	dest.UserRoleGroupID = m.UserRoleGroupID
	dest.UserUserID = m.UserUserID
	dest.Active = m.Active
	to = dest

	return nil
}

func (m *UserPermissionModel) Validate() *logger.HttpError {

	validate := validator.New()
	// TODO: Validate related
	/*	userAppRoute := NewAppRouteModel(m.ctx)
		 	if appErr := m.FindByID(m.AppRouteID, userAppRoute); appErr != nil {
				return appErr
			}

			userRoleGroup := NewUserRolesGroupModel(m.ctx)
			if appErr := m.FindByID(m.UserRoleGroupID, userRoleGroup); appErr != nil {
				return appErr
			}

			user := NewUserUserModel(m.ctx)
			if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
				return appErr
			}
	*/

	if err := validate.Var(m.Active, "required"); err != nil {
		fields := []string{"active"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid active ", err, nil)
	}

	return nil
}
