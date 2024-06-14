package models

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserUser = "user_users"
	locationUserUser   = "COMMON::MODELS::UserUser"
)

type UserUserModel struct {
	BaseModel     `json:"-" bson:"-"`
	ID            primitive.ObjectID  `json:"_id" bson:"_id"`
	FirstName     string              `json:"firstName" bson:"firstName"`
	Surename      string              `json:"surename" bson:"surename"`
	Email         string              `json:"email" bson:"email"`
	Password      *string             `json:"password,omitempty" bson:"password"`
	Terms         *time.Time          `json:"terms" bson:"terms"`
	AvatarUrl     string              `json:"avatarUrl" bson:"avatarUrl"`
	EmailVerified *time.Time          `json:"emailVerified" bson:"emailVerified"`
	Active        uint                `json:"active" bson:"active"`
	CreatedBy     primitive.ObjectID  `json:"createdBy" bson:"createdBy"`
	CreatedAt     time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedBy     primitive.ObjectID  `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt     time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedBy     *primitive.ObjectID `json:"deletedBy" bson:"deletedBy"`
	DeletedAt     *time.Time          `json:"deletedAt" bson:"deletedAt"`
}

type UserUsersModel struct {
	BaseModels
	Count int64            `json:"count"`
	Rows  *[]UserUserModel `json:"rows"`
}

func NewUserUserModel(c *gin.Context) *UserUserModel {
	m := &UserUserModel{}
	m.Init(c, locationUserUser, collectionUserUser)
	return m
}

func NewUserUsersModel(c *gin.Context) *UserUsersModel {
	m := &UserUsersModel{}
	m.Init(c, locationUserUser, collectionUserUser)
	return m
}

func (m *UserUserModel) FillMeta(create bool, delete bool) {

	now := time.Now().UTC()

	if create {
		m.ID = primitive.NewObjectID()
		m.CreatedBy = m.GetValues().Variables.User.ID
		m.CreatedAt = now
	} else if delete {
		m.DeletedBy = &m.GetValues().Variables.User.ID
		m.DeletedAt = &now
	}

	m.UpdatedBy = m.GetValues().Variables.User.ID
	m.UpdatedAt = now
}

func (m *UserUserModel) Bind() *service_errors.Error {
	return m.BaseModel.Bind(m, m.ctx)
}

func (m *UserUserModel) Validate() *service_errors.Error {

	validate := validator.New()
	if err := validate.Var(m.FirstName, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::UserUserModel", "name", "invalid first name")
	}

	if err := validate.Var(m.Surename, "required,gte=1"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::UserUserModel", "surename", "invalid surename")
	}

	if err := validate.Var(m.Email, "required,email"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::UserUserModel", "email", "invalid email")
	}

	if err := validate.Var(m.Password, "required,gte=6"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::UserUserModel", "password", "invalid password")
	}

	if err := validate.Var(m.Terms, "required"); err != nil {
		return service_log.Error(0, http.StatusBadRequest, "MODELS::UserUserModel", "terms", "invalid password")
	}

	return nil
}
