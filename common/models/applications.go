package models

import (
	"time"
)

type Application struct {
	ID          uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	Name        string     `json:"name" gorm:"column:name;type:string;size:255;not null;"`
	Description *string    `json:"description" gorm:"column:description;type:string;size:65536;"`
	Code        string     `json:"code" gorm:"column:code;type:string;size:255;"`
	Internal    int        `json:"internal" gorm:"column:internal;type:int;size:1;"`
	CreatedBy   uint       `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy   uint       `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy   *uint      `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt   *time.Time `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type Applications struct {
	Count int64          `json:"count"`
	Docs  *[]Application `json:"docs"`
}

func (m *Application) GetID() uint {
	return m.ID
}

func (m *Application) TableName() string {
	return "applications"
}
