package models

import (
	"time"
)

type History struct {
	ID             uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	OrganizationID uint       `json:"organizationId" gorm:"column:organization_id;type:uint;not null;"`
	Organization   *User      `json:"organization,omitempty"`
	Tablename      string     `json:"tableName" gorm:"column:table_name;type:string;size:255;not null;"`
	OriginalID     uint       `json:"originalId" gorm:"column:original_id;type:uint;not null;"`
	JsonData       string     `json:"jsonData" gorm:"column:json_data;type:string;size:65536;not null;"`
	CreatedBy      uint       `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy      uint       `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy      *uint      `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type Histories struct {
	Count int64      `json:"count"`
	Docs  *[]History `json:"docs"`
}

func (m *History) GetID() uint {
	return m.ID
}

func (m *History) TableName() string {
	return "history"
}
