package models

import (
	"time"
)

type UserGroup struct {
	ID             uint          `json:"id" gorm:"column:id;type:uint;primaryKey"`
	OrganizationID uint          `json:"organizationId" gorm:"column:organization_id;type:uint;not null;"`
	Organization   *Organization `json:"organization,omitempty"`
	AuthGroupID    uint          `json:"authGroupId" gorm:"column:auth_group_id;type:uint;not null;"`
	AuthGroup      *AuthGroup    `json:"authGroup,omitempty"`
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

type UserGroups struct {
	Count int64        `json:"count"`
	Docs  *[]UserGroup `json:"docs"`
}

func (m *UserGroup) GetID() uint {
	return m.ID
}

func (m *UserGroup) TableName() string {
	return "user_groups"
}
