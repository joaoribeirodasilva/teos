package models

import (
	"time"

	"gorm.io/gorm"
)

type UserReset struct {
	gorm.Model
	UserResetTypeID uint          `json:"resetTypeID" gorm:"column:user_reset_type_id;type:uint;not null"`
	UserResetType   UserResetType `json:"resetType,omitempty"`
	UserUserID      uint          `json:"userID" gorm:"column:user_user_id;type:uint;not null"`
	UserUser        UserUser      `json:"user,omitempty"`
	ResetKey        string        `json:"-" gorm:"column:reset_key;type:string;size:255;not null;"`
	Used            *time.Time    `json:"used" gorm:"column:used;type:time;"`
	Expire          time.Time     `json:"expire" gorm:"column:expire;type:time;not null;"`
	CreatedBy       uint          `json:"createdBy" gorm:"column:created_by;type:uint;;not null;"`
	UpdatedBy       uint          `json:"updatedBy" gorm:"column:updated_by;type:uint;;not null;"`
	DeletedBy       *uint         `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}
