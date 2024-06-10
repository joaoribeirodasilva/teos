package models

import (
	apps_models "github.com/joaoribeirodasilva/teos/apps/models"
	"gorm.io/gorm"
)

type HistHistory struct {
	gorm.Model
	AppAppID   uint               `json:"applicationId" gorm:"column:app_app_id;type:uint;not null;"`
	AppApp     apps_models.AppApp `json:"application,omitempty"`
	Table      string             `json:"table" gorm:"column:table;type:string;size:255;not null;"`
	OriginalID uint               `json:"originalId" gorm:"column:original_id;type:uint;not null;"`
	Data       string             `json:"data" gorm:"column:data;type:string;size:16777215;not null;"`
	CreatedBy  uint               `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	UpdatedBy  uint               `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	DeletedBy  *uint              `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type HistHistories struct {
	Count int            `json:"count"`
	Rows  *[]HistHistory `json:"rows"`
}
