package models

import (
	"time"
)

type UserPermission struct {
	ID             uint          `json:"id" gorm:"column:id;type:uint;primaryKey"`
	OrganizationID uint          `json:"organizationId" gorm:"column:organization_id;type:uint;not null;"`
	Organization   *Organization `json:"organization,omitempty"`
	AuthRoleID     uint          `json:"authRoleId" gorm:"column:auth_role_id;type:uint;not null;"`
	AuthRole       *AuthRole     `json:"authRole,omitempty"`
	AppRouteID     uint          `json:"appRouteId" gorm:"column:app_route_id;type:uint;not null;"`
	AppRoute       *AppRoute     `json:"appRoute,omitempty"`
	UserID         uint          `json:"userId" gorm:"column:user_id;type:uint;not null;"`
	User           *User         `json:"user,omitempty"`
	Active         int           `json:"active" gorm:"column:active;type:int;size(1);"`
	CreatedBy      uint          `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt      time.Time     `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy      uint          `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt      time.Time     `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy      *uint         `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt      *time.Time    `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type UserPermissions struct {
	Count int64             `json:"count"`
	Docs  *[]UserPermission `json:"docs"`
}

func (m *UserPermission) GetID() uint {
	return m.ID
}

func (m *UserPermission) TableName() string {
	return "user_permissions"
}
