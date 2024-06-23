package models

import (
	"time"
)

type LogMessage struct {
	ID        uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	Time      time.Time  `json:"time"`
	Type      string     `json:"type"`
	Fields    *[]string  `json:"fields,omitempty"`
	Message   string     `json:"message"`
	Data      *string    `json:"data,omitempty"`
	CreatedBy uint       `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy uint       `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy *uint      `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type LogMessages struct {
	Count int64         `json:"count"`
	Docs  *[]LogMessage `json:"docs"`
}

func (m *LogMessage) GetID() uint {
	return m.ID
}

func (m *LogMessage) TableName() string {
	return "reset_types"
}
