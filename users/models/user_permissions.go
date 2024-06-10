package models

import (
	apps_models "github.com/joaoribeirodasilva/teos/apps/models"
	"gorm.io/gorm"
)

type UserPermission struct {
	gorm.Model
	UserRouteID     uint                 `json:"routeId" gorm:"column:app_route_id;type:uint;not null;"`
	UserRoute       apps_models.AppRoute `json:"route,omitempty"`
	UserRoleGroupID uint                 `json:"roleGroupId" gorm:"column:user_role_group_id;type:uint;not null;"`
	UserRoleGroup   UserRolesGroup       `json:"roleGroup,omitempty"`
	UserUserID      uint                 `json:"userId" gorm:"column:user_user_id;type:uint;not null;"`
	UserUser        UserUser             `json:"user,omitempty"`
	Active          uint                 `json:"active" gorm:"column:active;type:uint;not null;default:0"`
	CreatedBy       uint                 `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	UpdatedBy       uint                 `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	DeletedBy       *uint                `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type UserPermissions struct {
	Count int               `json:"count"`
	Rows  *[]UserPermission `json:"rows"`
}
