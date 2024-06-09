package models

import "gorm.io/gorm"

type HistHistory struct {
	gorm.Model
	ApplicationID uint `json:"applicationId" gorm:"column:app_app_id;type:uint;not null;"`
	//Application   Application `json:"applicationId" gorm:"column:app_app_id;type:uint;not null;"`
	Table      string `json:"table" gorm:"column:table;type:string;size:255;not null;"`
	OriginalID uint   `json:"originalId" gorm:"column:original_id;type:uint;not null;"`
	Data       string `json:"data" gorm:"column:data;type:string;size:16777215;not null;"`
	CreatedBy  uint   `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	UpdatedBy  uint   `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	DeletedBy  *uint  `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type HistHistories struct {
	Count int            `json:"count"`
	Rows  *[]HistHistory `json:"rows"`
}
