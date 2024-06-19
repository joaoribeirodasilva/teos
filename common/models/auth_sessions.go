package models

import (
	"time"
)

type AuthSession struct {
	ID        uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	UserID    uint       `json:"userId" gorm:"column:user_id;type:uint;not null;"`
	User      *User      `json:"user,omitempty"`
	CreatedBy uint       `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy uint       `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy *uint      `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type AuthSessions struct {
	Count int64          `json:"count"`
	Docs  *[]AuthSession `json:"docs"`
}

func (m *AuthSession) GetID() uint {
	return m.ID
}

func (m *AuthSession) TableName() string {
	return "auth_sessions"
}
