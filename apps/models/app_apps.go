package models

import (
	"gorm.io/gorm"
)

type AppApp struct {
	gorm.Model
	Name        string  `json:"name" gorm:"column:name;type:string;size:255;not null"`
	Description *string `json:"description" gorm:"column:description;type:string;size:65536"`
	AppKey      string  `json:"appKey" gorm:"column:app_key;type:string;size:128;not null"`
	Active      uint    `json:"active" gorm:"column:active;type:uint;not null;default:0"`
	CreatedBy   uint    `json:"createdBy" gorm:"column:created_by;type:uint;;not null;"`
	UpdatedBy   uint    `json:"updatedBy" gorm:"column:updated_by;type:uint;;not null;"`
	DeletedBy   *uint   `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type AppApps struct {
	Count int       `json:"count"`
	Rows  *[]AppApp `json:"rows"`
}
