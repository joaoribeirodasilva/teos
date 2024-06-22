package models

import (
	"time"
)

type User struct {
	ID            uint       `json:"id" gorm:"column:id;type:uint;primaryKey"`
	FirstName     string     `json:"firstName" gorm:"column:first_name;type:string;size:255;not null;"`
	Surname       string     `json:"surname" gorm:"column:surname;type:string;size:255;not null;"`
	Email         string     `json:"email" gorm:"column:email;type:string;size:255;not null;unique;"`
	Password      string     `json:"password,omitempty" gorm:"column:password;type:string;size:255;not null;"`
	Terms         *time.Time `json:"terms" gorm:"column:terms;type:time;"`
	AvatarUrl     string     `json:"avatarUrl" gorm:"column:avatar_url;type:string;size:255;not null;default:0"`
	EmailVerified *time.Time `json:"emailVerified" gorm:"column:email_verified;type:time;not null;"`
	Active        uint       `json:"active" gorm:"column:active;type:int;default:0"`
	CreatedBy     uint       `json:"createdBy" gorm:"column:created_by;type:uint;not null;"`
	CreatedAt     time.Time  `json:"createdAt" gorm:"column:created_at;type:time;not null;"`
	UpdatedBy     uint       `json:"updatedBy" gorm:"column:updated_by;type:uint;not null;"`
	UpdatedAt     time.Time  `json:"updatedAt" gorm:"column:updated_at;type:time;not null;"`
	DeletedBy     *uint      `json:"deletedBy" gorm:"column:deleted_by;type:uint"`
	DeletedAt     *time.Time `json:"deletedAt" gorm:"column:deleted_at;type:time"`
}

type Users struct {
	Count int64   `json:"count"`
	Docs  *[]User `json:"docs"`
}

func (m *User) GetID() uint {
	return m.ID
}

func (m *User) TableName() string {
	return "users"
}
