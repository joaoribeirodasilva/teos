package models

import (
	"time"
)

type AppRoute struct {
	ID               uint            `json:"id" gorm:"column:id;type:uint;primaryKey"`
	ApplicationID    uint            `json:"applicationId" gorm:"column:application_id;type:uint;not null;"`
	Application      *Application    `json:"application,omitempty"`
	AppEnvironmentID uint            `json:"appEnvironmentId" gorm:"column:app_environment_id;type:uint;not null;"`
	AppEnvironment   *AppEnvironment `json:"appEnvironment,omitempty"`
	Name             string          `json:"name" gorm:"column:name;type:string;size:255;not null;"`
	Description      *string         `json:"description" gorm:"column:description;type:string;size:65536;"`
	Uri              string          `json:"uri" gorm:"column:uri;type:string;size:255;not null;"`
	Active           string          `json:"active" gorm:"column:active;type:int;not null;default:0"`
	CreatedBy        uint            `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt        time.Time       `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy        uint            `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt        time.Time       `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy        *uint           `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt        *time.Time      `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type AppRoutes struct {
	Count int64       `json:"count"`
	Docs  *[]AppRoute `json:"docs"`
}

func (m *AppRoute) GetID() uint {
	return m.ID
}

func (m *AppRoute) TableName() string {
	return "app_routes"
}
