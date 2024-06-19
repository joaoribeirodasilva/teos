package models

import (
	"time"
)

type Methods string

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodPatch  = "PATCH"
	MethodDelete = "DELETE"
)

type AppRouteMethod struct {
	ID          uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	AppRouteID  uint       `json:"appRouteId" gorm:"column:app_route_id;type:uint;not null;"`
	AppRoute    *AppRoute  `json:"appRoute,omitempty"`
	Name        string     `json:"name" gorm:"column:name;type:string;size:255;not null;"`
	Description *string    `json:"description" gorm:"column:description;type:string;size:65536;"`
	Method      Methods    `json:"method" gorm:"column:method;type:string"`
	Uri         string     `json:"uri" gorm:"column:uri;type:string;size:255;not null;"`
	Open        int        `json:"open" gorm:"column:open;type:int;size(1);"`
	Active      int        `json:"active" gorm:"column:active;type:int;size(1);"`
	CreatedBy   uint       `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy   uint       `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy   *uint      `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt   *time.Time `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type AppRouteMethods struct {
	Count int64             `json:"count"`
	Docs  *[]AppRouteMethod `json:"docs"`
}

func (m *AppRouteMethod) GetID() uint {
	return m.ID
}

func (m *AppRouteMethod) TableName() string {
	return "app_route_methods"
}
