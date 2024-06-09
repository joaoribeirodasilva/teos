package models

import (
	"gorm.io/gorm"
)

type AppConfiguration struct {
	gorm.Model
	ApplicationID uint `json:"appID" gorm:"column:app_app_id;type:uint;not null"`
	//Application   Application   `json:"application,omitempty"`
	Name        string   `json:"name" gorm:"column:name;type:string;size:255;not null"`
	Description *string  `json:"description" gorm:"column:description;type:string;size:65536"`
	Key         string   `json:"key" gorm:"column:key;type:string;size:128;not null"`
	ValueInt    *int     `json:"valueInt" gorm:"column:value_int;type:int"`
	ValueChar   *string  `json:"valueChar" gorm:"column:value_char;type:string;size:65536"`
	ValueFloat  *float64 `json:"valueFloat" gorm:"column:value_float;type:float64;"`
	CreatedBy   uint     `json:"createdBy" gorm:"column:created_by;type:uint;;not null;"`
	UpdatedBy   uint     `json:"updatedBy" gorm:"column:updated_by;type:uint;;not null;"`
	DeletedBy   *uint    `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type AppConfigurations struct {
	Count int                 `json:"count"`
	Rows  *[]AppConfiguration `json:"rows"`
}
