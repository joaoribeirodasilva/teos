package models

import (
	"time"
)

type SvcRoute struct {
	ID        uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	App       string     `json:"name" gorm:"column:name;type:string;size:255;not null;"`
	Method    string     `json:"method" gorm:"column:method;type:string;size:255;not null;"`
	Uri       string     `json:"uri" gorm:"column:uri;type:string;size:255;not null;"`
	Open      uint       `json:"open" gorm:"column:open;type:int;default:0"`
	CreatedBy uint       `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy uint       `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy *uint      `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type SvcRoutes struct {
	Count int64       `json:"count"`
	Docs  *[]SvcRoute `json:"docs"`
}

func (m *SvcRoute) GetID() uint {
	return m.ID
}

func (m *SvcRoute) TableName() string {
	return "reset_types"
}
