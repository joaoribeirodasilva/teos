package configuration

import "time"

type ConfigurationApp struct {
	ID   uint   `json:"_id" gorm:"column:id"`
	Code string `json:"code" gorm:"column:code"`
}

func (ca *ConfigurationApp) TableName() string {
	return "applications"
}

type EnvironmentApp struct {
	ID  uint   `json:"_id" gorm:"column:id"`
	Key string `json:"key" gorm:"column:key"`
}

func (ca *EnvironmentApp) TableName() string {
	return "app_environments"
}

type ConfigurationValue struct {
	ID               uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	AppEnvironmentID uint       `json:"appEnvironmentId" gorm:"column:app_environment_id;type:uint;not null;"`
	ApplicationID    uint       `json:"applicationId" gorm:"column:application_ịId;type:uint;not null;"`
	ConfigurationKey string     `json:"configurationKey" gorm:"column:configuration_key;type:string;size:255;not null;"`
	ValString        *string    `json:"valString" gorm:"column:val_string;type:string;size:65536;"`
	ValInt           *int64     `json:"valInt" gorm:"column:val_int;type:int64;"`
	ValDouble        *float64   `json:"valDouble" gorm:"column:val_double;type:float64;"`
	ValBool          *int       `json:"valBoolean" gorm:"column:val_boolean;type:int;size(1);"`
	ValDate          *time.Time `json:"valDate" gorm:"column:val_date;type:DATE;"`
	ValTime          *time.Time `json:"valTime" gorm:"column:val_time;type:TIME;"`
	ValDateTime      *time.Time `json:"valDateTime" gorm:"column:val_datetime;type:TIMESTAMP;"`
	Type             string     `json:"type" gorm:"column:type;type:string;"`
}

func (ca *ConfigurationValue) TableName() string {
	return "app_configurations"
}
