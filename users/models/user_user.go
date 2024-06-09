package models

import (
	"time"

	"gorm.io/gorm"
)

type UserUser struct {
	gorm.Model
	FirstName     string     `json:"firstName" gorm:"column:first_name;type:string;size:255;not null;"`
	Surename      string     `json:"surename" gorm:"column:surename;type:string;size:255;not null;"`
	Email         string     `json:"email" gorm:"column:email;type:string;size:255;not null;"`
	Password      *string    `json:"password,omitempty" gorm:"column:password;type:string;size:255;"`
	Terms         *time.Time `json:"terms" gorm:"column:terms;type:time;"`
	AvatarUrl     string     `json:"avatarUrl" gorm:"column:avatar_url;type:string;size:255;"`
	EmailVerified *time.Time `json:"emailVerified" gorm:"column:email_verified;type:time;"`
	Active        uint       `json:"active" gorm:"column:active;type:uint;;not null;default:0"`
	CreatedBy     uint       `json:"createdBy" gorm:"column:created_by;type:uint;;not null;"`
	UpdatedBy     uint       `json:"updatedBy" gorm:"column:updated_by;type:uint;;not null;"`
	DeletedBy     *uint      `json:"deletedBy" gorm:"column:deleted_by;type:uint;"`
}

type UserUsers struct {
	Count int         `json:"count"`
	Rows  *[]UserUser `json:"rows"`
}
