package models

import (
	"gorm.io/gorm"
)

type AppRoutesBlock struct {
	gorm.Model
	ApplicationID uint    `json:"appId" gorm:"column:app_app_id;type:uint;not null"`
	Application   AppApp  `json:"application,omitempty"`
	Name          string  `json:"name" gorm:"column:name;type:string;size:255;not null"`
	Description   *string `json:"description" gorm:"column:description;type:string;size:65536"`
	Route         string  `json:"route" gorm:"column:route;type:string;size:255;not null"`
	Active        uint    `json:"active" gorm:"column:active;type:uint;not null;default:0"`
	CreatedBy     uint    `json:"createdBy" gorm:"column:created_by;type:uint;;not null;"`
	UpdatedBy     uint    `json:"updatedBy" gorm:"column:updated_by;type:uint;;not null;"`
	DeletedBy     *uint   `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type AppRoutesBlocks struct {
	Count int               `json:"count"`
	Rows  *[]AppRoutesBlock `json:"rows"`
}
