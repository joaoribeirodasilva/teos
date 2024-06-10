package models

import (
	"gorm.io/gorm"
)

type UserRolesGroup struct {
	gorm.Model
	UserRoleID  uint      `json:"roleId" gorm:"column:user_role_id;type:uint;not null;"`
	UserRole    UserRole  `json:"role,omitempty"`
	UserGroupID uint      `json:"groupId" gorm:"column:user_group_id;type:uint;not null;"`
	UserGroup   UserGroup `json:"group,omitempty"`
	Active      uint      `json:"active" gorm:"column:active;type:uint;not null;default:0"`
	CreatedBy   uint      `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	UpdatedBy   uint      `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	DeletedBy   *uint     `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type UserRolesGroups struct {
	Count int               `json:"count"`
	Rows  *[]UserRolesGroup `json:"rows"`
}
