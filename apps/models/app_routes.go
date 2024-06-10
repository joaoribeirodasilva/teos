package models

import (
	"gorm.io/gorm"
)

type AppRoute struct {
	gorm.Model
	AppRoutesBlockID uint           `json:"routesBlockID" gorm:"column:app_routes_block_id;type:uint;not null"`
	AppRoutesBlock   AppRoutesBlock `json:"routesBlock,omitempty"`
	Name             string         `json:"name" gorm:"column:name;type:string;size:255;not null"`
	Description      *string        `json:"description" gorm:"column:description;type:string;size:65536"`
	Method           string         `json:"method" gorm:"column:method;type:string;size:255;not null"`
	Route            string         `json:"route" gorm:"column:route;type:string;size:255;not null"`
	Open             uint           `json:"open" gorm:"column:open;type:uint;not null;default:0"`
	Active           uint           `json:"active" gorm:"column:active;type:uint;not null;default:0"`
	CreatedBy        uint           `json:"createdBy" gorm:"column:created_by;type:uint;;not null;"`
	UpdatedBy        uint           `json:"updatedBy" gorm:"column:updated_by;type:uint;;not null;"`
	DeletedBy        *uint          `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type AppRoutes struct {
	Count int         `json:"count"`
	Rows  *[]AppRoute `json:"rows"`
}
