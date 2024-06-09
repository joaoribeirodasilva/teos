package models

import "gorm.io/gorm"

type UserSession struct {
	gorm.Model
	UserUserID uint     `json:"userId" gorm:"column:user_user_id;type:uint;not null;"`
	UserUser   UserUser `json:"user,omitempty"`
	CreatedBy  uint     `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	UpdatedBy  uint     `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	DeletedBy  *uint    `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type UserSessions struct {
	Count int            `json:"count"`
	Rows  *[]UserSession `json:"rows"`
}
