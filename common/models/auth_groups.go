package models

import (
	"time"
)

type AuthGroup struct {
	ID             uint          `json:"id" gorm:"column:id;type:uint;primaryKey"`
	OrganizationID uint          `json:"organizationId" gorm:"column:organization_id;type:uint;not null;"`
	Organization   *Organization `json:"organization,omitempty"`
	Name           string        `json:"name" gorm:"column:name;type:string;size:255;not null;"`
	Description    *string       `json:"description" gorm:"column:description;type:string;size:65536;"`
	CreatedBy      uint          `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt      time.Time     `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy      uint          `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt      time.Time     `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy      *uint         `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt      *time.Time    `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type AuthGroups struct {
	Count int64        `json:"count"`
	Docs  *[]AuthGroup `json:"docs"`
}

func (m *AuthGroup) GetID() uint {
	return m.ID
}

func (m *AuthGroup) TableName() string {
	return "auth_groups"
}
