package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/payload"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"gorm.io/gorm"
)

type ApplicationsService struct {
	payload *payload.Payload
	db      *gorm.DB
	request *payload.HttpRequest
	history *redisdb.RedisDB
}

func NewApplicationsService(payload *payload.Payload) *ApplicationsService {
	return &ApplicationsService{
		payload: payload,
		db:      payload.Services.Db.GetDatabase(),
		request: payload.Http.Request,
		history: payload.Services.HistoryDb,
	}

}

// List returns a list of users from the collection
func (s *ApplicationsService) List(filter string, args ...any) (*models.Applications, *logger.HttpError) {

	model := models.Application{}
	models := models.Applications{}

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

	return &models, nil
}

// Get returns a single user from the collection
func (s *ApplicationsService) Get(model *models.Application, filter string, args ...any) *logger.HttpError {

	query := s.db.Model(model)
	if filter == "" {
		query.Where("id = ?", s.request.ID)
	} else {
		query.Where(s.request.Query.Filter)
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

	return nil
}

// Create creates a new user document or returns a logger.HttpError in case of error
func (s *ApplicationsService) Create(model *models.Application) *logger.HttpError {

	if s.request.Session.Auth.UserSession.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to create this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
	}

	if err := s.Validate(model); err != nil {
		return err
	}

	exists := models.Application{}

	if err := s.db.Where(
		"name = ? OR code = ?",
		model.Name,
		model.Code,
	).First(&exists).Error; err != nil {

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

		s.assign(&exists, model, payload.SVC_OPERATION_CREATE)
	} else {
		s.assign(model, nil, payload.SVC_OPERATION_CREATE)
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
func (s *ApplicationsService) Update(model *models.Application) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	// Security
	if s.request.Session.Auth.UserSession.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to update this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
	}

	exists := models.Application{}
	if err := s.db.Where(
		"name = ? OR code = ?",
		model.Name,
		model.Code,
	).First(&exists).Error; err != nil {

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

	s.assign(&exists, model, payload.SVC_OPERATION_UPDATE)

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
func (s *ApplicationsService) Delete(id uint) *logger.HttpError {

	exists := &models.Application{}

	// Security
	if s.request.Session.Auth.UserSession.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to delete this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
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

	s.assign(exists, nil, payload.SVC_OPERATION_DELETE)

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

func (s *ApplicationsService) Validate(model *models.Application) *logger.HttpError {

	validate := validator.New()

	if err := validate.Var(model.Name, "required;gt=1"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}

	if err := validate.Var(model.Code, "required;gt=1"); err != nil {
		fields := []string{"code"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid code ", err, nil)
	}

	if err := validate.Var(model.Internal, "required"); err != nil {
		fields := []string{"internal"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid internal ", err, nil)
	}

	return nil
}

func (s *ApplicationsService) assign(to *models.Application, from *models.Application, operation payload.Operation) {

	now := time.Now().UTC()

	if operation == payload.SVC_OPERATION_CREATE {

		to.CreatedBy = s.request.Session.Auth.UserSession.ID
		to.CreatedAt = now

	} else if operation == payload.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.request.Session.Auth.UserSession.ID
		to.DeletedAt = &now

	} else {

		to.Name = from.Name
		to.Description = from.Description
		to.Code = from.Code
		to.Internal = from.Internal
	}

	to.UpdatedBy = s.request.Session.Auth.UserSession.ID
	to.UpdatedAt = now
}
