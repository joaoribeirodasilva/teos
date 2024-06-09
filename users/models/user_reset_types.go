package models

import "gorm.io/gorm"

type UserResetType struct {
	gorm.Model
	Name        string  `json:"name" gorm:"column:name;type:string;size:255;not null;"`
	Description *string `json:"description" gorm:"column:description;type:string;size:65536;"`
	CreatedBy   uint    `json:"createdBy" gorm:"column:created_by;type:uint;;not null;"`
	UpdatedBy   uint    `json:"updatedBy" gorm:"column:updated_by;type:uint;;not null;"`
	DeletedBy   *uint   `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}
