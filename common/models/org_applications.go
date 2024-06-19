package models

import (
	"time"
)

type OrgApplication struct {
	ID               uint            `json:"id" gorm:"column:id;type:uint;primaryKey"`
	OrganizationId   uint            `json:"organizationId" gorm:"column:organization_id;type:uint;not null;"`
	Organization     *Organization   `json:"organization,omitempty"`
	ApplicationID    uint            `json:"applicationId" gorm:"column:application_id;type:uint;not null;"`
	Application      *Application    `json:"application,omitempty"`
	AppEnvironmentId uint            `json:"appEnvironmentId" gorm:"column:app_environment_id;type:uint;not null;"`
	AppEnvironment   *AppEnvironment `json:"appEnvironment,omitempty"`
	Active           int             `json:"active" gorm:"column:active;type:int;size(1);"`
	CreatedBy        uint            `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt        time.Time       `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy        uint            `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt        time.Time       `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy        *uint           `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt        *time.Time      `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type OrgApplications struct {
	Count int64              `json:"count"`
	Docs  *[]OrgApplications `json:"docs"`
}

func (m *OrgApplication) GetID() uint {
	return m.ID
}

func (m *OrgApplication) TableName() string {
	return "org_applications"
}