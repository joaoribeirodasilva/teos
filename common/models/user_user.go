package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/dbtest/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	collectionUserUser = "user_users"
)

type UserUserModel struct {
	BaseModel     `json:"-" bson:"-"`
	ID            primitive.ObjectID  `json:"_id" bson:"_id"`
	FirstName     string              `json:"firstName" bson:"firstName"`
	Surname       string              `json:"surname" bson:"surname"`
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

func (m *UserUserModel) GetCollectionName() string {
	return collectionUserUser
}

func (m *UserUserModel) AssignValues(to interface{}) error {

	dest, ok := to.(*UserUserModel)
	if !ok {
		return ErrWrongModelType
	}
	dest.ID = m.ID
	dest.FirstName = m.FirstName
	dest.Email = m.Email
	dest.Password = m.Password
	dest.Terms = m.Terms
	dest.AvatarUrl = m.AvatarUrl
	dest.EmailVerified = m.EmailVerified
	dest.Active = m.Active
	to = dest

	return nil
}

func (m *UserUserModel) Validate() *logger.HttpError {

	validate := validator.New()
	if err := validate.Var(m.FirstName, "required"); err != nil {
		fields := []string{"firstName"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid firstName ", err, nil)
	}

	if err := validate.Var(m.Surname, "required,gte=1"); err != nil {
		fields := []string{"surname"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid surname ", err, nil)
	}

	if err := validate.Var(m.Email, "required,email"); err != nil {
		fields := []string{"email"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid email ", err, nil)
	}

	if err := validate.Var(m.Password, "required,gte=6"); err != nil {
		fields := []string{"password"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid password ", err, nil)
	}

	if err := validate.Var(m.Terms, "required"); err != nil {
		fields := []string{"terms"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid terms ", err, nil)
	}

	return nil
}
