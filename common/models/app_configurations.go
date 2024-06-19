package models

import (
	"time"
)

type ConfigurationType string

const (
	TypeString   = "string"
	TypeInt      = "int"
	TypeDouble   = "double"
	TypeBoolean  = "boolean"
	TypeDate     = "date"
	TypeTime     = "time"
	TypeDateTime = "datetime"
)

type AppConfiguration struct {
	ID               uint              `json:"id" gorm:"column:id;type:uint;primaryKey"`
	AppEnvironmentID uint              `json:"appEnvironmentId" gorm:"column:app_environment_id;type:uint;not null;"`
	AppEnvironment   *AppEnvironment   `json:"appEnvironment,omitempty"`
	ApplicationID    uint              `json:"applicationId" gorm:"column:application_ịId;type:uint;not null;"`
	Application      *Application      `json:"application,omitempty"`
	ConfigurationKey string            `json:"configurationKey" gorm:"column:configuration_key;type:string;size:255;not null;"`
	ValString        *string           `json:"valString" gorm:"column:val_string;type:string;size:65536;"`
	ValInt           *int64            `json:"valInt" gorm:"column:val_int;type:int64;"`
	ValDouble        *float64          `json:"valDouble" gorm:"column:val_double;type:float64;"`
	ValBool          *int              `json:"valBoolean" gorm:"column:val_boolean;type:int;size(1);"`
	ValDate          *time.Time        `json:"valDate" gorm:"column:val_date;type:DATE;"`
	ValTime          *time.Time        `json:"valTime" gorm:"column:val_time;type:TIME;"`
	ValDateTime      *time.Time        `json:"valDateTime" gorm:"column:val_datetime;type:TIMESTAMP;"`
	Type             ConfigurationType `json:"type" gorm:"column:type;type:string;"`
	CreatedBy        uint              `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt        time.Time         `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy        uint              `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt        time.Time         `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy        *uint             `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt        *time.Time        `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type AppConfigurations struct {
	Count int64               `json:"count"`
	Docs  *[]AppConfiguration `json:"docs"`
}

func (m *AppConfiguration) GetID() uint {
	return m.ID
}

func (m *AppConfiguration) TableName() string {
	return "app_configurations"
}
