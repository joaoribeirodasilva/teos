package models

import (
	"gorm.io/gorm"
)

type UserRole struct {
	gorm.Model
	Name        string  `json:"name" gorm:"column:name;type:string;size:255;not null"`
	Description *string `json:"description" gorm:"column:description;type:string;size:65536"`
	GroupPath   *string `json:"groupPath" gorm:"column:group_path;type:string;size:2048"`
	CreatedBy   uint    `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	UpdatedBy   uint    `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	DeletedBy   *uint   `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type UserRoles struct {
	Count int         `json:"count"`
	Rows  *[]UserRole `json:"rows"`
}
