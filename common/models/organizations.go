package models

import (
	"time"
)

type Organization struct {
	ID        uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	Name      string     `json:"name" gorm:"column:name;type:string;size:255;not null;"`
	CreatedBy uint       `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy uint       `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy *uint      `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type Organizations struct {
	Count int64           `json:"count"`
	Docs  *[]Organization `json:"docs"`
}

func (m *Organization) GetID() uint {
	return m.ID
}

func (m *Organization) TableName() string {
	return "organizations"
}
