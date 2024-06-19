package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/services"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
	"gorm.io/gorm"
)

type UsersService struct {
	services      *structures.RequestValues
	db            *gorm.DB
	user          *token.User
	query         *requests.QueryString
	sessionDb     *redisdb.RedisDB
	permissionsDb *redisdb.RedisDB
	historyDb     *redisdb.RedisDB
}

func New(services *structures.RequestValues) *UsersService {
	s := &UsersService{}
	s.services = services
	s.db = services.Services.Db.GetDatabase()
	s.user = services.User
	s.query = &services.Query
	s.sessionDb = services.Services.SessionsDB
	s.permissionsDb = services.Services.PermissionsDB
	s.historyDb = services.Services.HistoryDB
	return s
}

// List returns a list of users from the collection
func (s *UsersService) List(filter string, args ...any) (*models.Users, *logger.HttpError) {

	model := models.Users{}
	models := models.Users{}

	if err := s.db.Model(&model).Where(filter, args).Count(&models.Count).Error; err != nil {

		return nil, logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to query database",
			err,
			nil,
		)
	}

	if err := s.db.Model(&model).Where(filter, args).Find(&models.Docs).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, logger.Error(
				logger.LogStatusNotFound,
				nil,
				"no documents found",
				nil,
				nil,
			)
		}

		return nil, logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to query database",
			err,
			nil,
		)
	}

	for idx := range *models.Docs {
		(*models.Docs)[idx].Password = ""
	}

	return &models, nil
}

// Get returns a single user from the collection
func (s *UsersService) Get(model *models.User, filter string, args ...any) *logger.HttpError {

	query := s.db.Model(model)
	if filter == "" {
		query.Where("id = ?", s.query.ID)
	} else {
		query.Where(filter, args)
	}

	if err := query.First(model).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return logger.Error(
				logger.LogStatusNotFound,
				nil,
				"no documents found",
				nil,
				nil,
			)
		}

		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to query database",
			err,
			nil,
		)

	}

	model.Password = ""

	return nil
}

// Create creates a new user document or returns a logger.HttpError in case of error
func (s *UsersService) Create(model *models.User) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	//TODO: organization config any can create user or only the organization

	exists := models.User{}

	if err := s.db.Where("email = ?", model.Email).First(&exists).Error; err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {

			return logger.Error(
				logger.LogStatusInternalServerError,
				nil,
				"failed to query database",
				err,
				nil,
			)
		}
	}

	// Send exists to history

	if exists.DeletedBy != nil || exists.DeletedAt != nil {

		exists.DeletedAt = nil
		exists.DeletedBy = nil

		s.assign(&exists, model, services.SVC_OPERATION_CREATE)
	} else {
		s.assign(model, nil, services.SVC_OPERATION_CREATE)
	}

	if err := s.db.Create(model).Error; err != nil {

		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to save document into database",
			err,
			nil,
		)
	}

	return nil
}

// Create updates a user document or returns a logger.HttpError in case of error
func (s *UsersService) Update(model *models.User) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	// Security
	if s.user.ID != model.ID && s.user.OrganizationID != 1 {

		err := errors.New("user documents can only be changed by the owner")
		return logger.Error(
			logger.LogStatusUnauthorized,
			nil,
			"you don't have enough privileges to change an user document",
			err,
			nil,
		)
	}

	exists := models.User{}
	if err := s.db.Where("email = ?", model.Email).First(&exists).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return logger.Error(
				logger.LogStatusNotFound,
				nil,
				"document not found",
				err,
				nil,
			)
		}
	}

	if exists.ID != model.ID {

		err := errors.New("document exists")
		return logger.Error(
			logger.LogStatusConflict,
			nil,
			fmt.Sprintf("document already exists with the same uniqueness and id: %d", exists.ID),
			err,
			nil,
		)
	}

	s.assign(&exists, model, services.SVC_OPERATION_UPDATE)

	if err := s.db.Save(exists).Error; err != nil {

		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to save document into database",
			err,
			nil,
		)
	}

	return nil
}

// Delete deletes a user document or returns a logger.HttpError in case of error
func (s *UsersService) Delete(id uint) *logger.HttpError {

	exists := &models.User{}

	// Security
	if s.user.ID != id && s.user.OrganizationID != 1 {

		err := errors.New("user documents can only be deleted by the owner")
		return logger.Error(
			logger.LogStatusUnauthorized,
			nil,
			"you don't have enough privileges to delete an user document",
			err,
			nil,
		)
	}

	if err := s.db.Where("id = ?", exists).First(&exists).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return logger.Error(
				logger.LogStatusNotFound,
				nil,
				"document not found",
				err,
				nil,
			)
		}
	}

	if err := s.db.Delete("id = ?", id).Error; err != nil {

		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to save document into database",
			err,
			nil,
		)
	}

	return nil
}

func (m *UsersService) Validate(model *models.User) *logger.HttpError {

	validate := validator.New()
	if err := validate.Var(model.FirstName, "required"); err != nil {
		fields := []string{"firstName"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid firstName ", err, nil)
	}

	if err := validate.Var(model.Surname, "required,gte=1"); err != nil {
		fields := []string{"surname"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid surname ", err, nil)
	}

	if err := validate.Var(model.Email, "required,email"); err != nil {
		fields := []string{"email"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid email ", err, nil)
	}

	if err := validate.Var(model.Password, "required,gte=6"); err != nil {
		fields := []string{"password"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid password ", err, nil)
	}

	if err := validate.Var(model.Terms, "required"); err != nil {
		fields := []string{"terms"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid terms ", err, nil)
	}

	return nil
}

func (s *UsersService) assign(to *models.User, from *models.User, operation services.Operation) {

	now := time.Now().UTC()

	if operation == services.SVC_OPERATION_CREATE {

		to.CreatedBy = s.user.ID
		to.CreatedAt = now

	} else if operation == services.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.user.ID
		to.DeletedAt = &now

	} else {

		to.FirstName = from.FirstName
		to.Surname = from.Surname
		to.Email = from.Email
		to.Password = from.Password
		to.Terms = from.Terms
		to.AvatarUrl = from.AvatarUrl
		to.EmailVerified = from.EmailVerified
		to.Active = from.Active
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
