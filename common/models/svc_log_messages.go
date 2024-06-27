package models

import (
	"time"
)

type LogMessage struct {
	ID        uint      `json:"id" gorm:"column:id;type:uint;primaryKey"`
	App       string    `json:"app" gorm:"column:app;type:string;size:255;not null"`
	Time      time.Time `json:"time"`
	Type      string    `json:"type"`
	Data      *string   `json:"data,omitempty"`
	CreatedBy uint      `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
}

type LogMessages struct {
	Count int64         `json:"count"`
	Docs  *[]LogMessage `json:"docs"`
}

func (m *LogMessage) GetID() uint {
	return m.ID
}

func (m *LogMessage) TableName() string {
	return "svc_log_messages"
}
