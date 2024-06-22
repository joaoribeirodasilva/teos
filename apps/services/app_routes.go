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

type AppRoutesService struct {
	payload *payload.Payload
	db      *gorm.DB
	request *payload.HttpRequest
	history *redisdb.RedisDB
}

func NewAppRoutesService(payload *payload.Payload) *AppRoutesService {
	return &AppRoutesService{
		payload: payload,
		db:      payload.Services.Db.GetDatabase(),
		request: payload.Http.Request,
		history: payload.Services.HistoryDb,
	}
}

// List returns a list of users from the collection
func (s *AppRoutesService) List(filter string, args ...any) (*models.AppRoutes, *logger.HttpError) {

	model := models.AppRoute{}
	models := models.AppRoutes{}

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
func (s *AppRoutesService) Get(model *models.AppRoute, filter string, args ...any) *logger.HttpError {

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
func (s *AppRoutesService) Create(model *models.AppRoute) *logger.HttpError {

	if s.request.Session.Auth.UserSession.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to create this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
	}

	if err := s.Validate(model); err != nil {
		return err
	}

	//TODO: organization config any can create user or only the organization

	exists := models.AppRoute{}

	if err := s.db.Where(
		"application_id = ? AND app_environment_id = ? AND (name = ? OR uri = ?) ",
		model.ApplicationID,
		model.AppEnvironmentID,
		model.Name,
		model.Uri,
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
func (s *AppRoutesService) Update(model *models.AppRoute) *logger.HttpError {

	if s.request.Session.Auth.UserSession.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to update this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
	}

	if err := s.Validate(model); err != nil {
		return err
	}

	// Security
	// TODO: security

	exists := models.AppRoute{}
	if err := s.db.Where(
		"application_id = ? AND app_environment_id = ? AND (name = ? OR uri = ?) ",
		model.ApplicationID,
		model.AppEnvironmentID,
		model.Name,
		model.Uri,
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
func (s *AppRoutesService) Delete(id uint) *logger.HttpError {

	exists := &models.AppRoute{}

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

func (s *AppRoutesService) Validate(model *models.AppRoute) *logger.HttpError {

	validate := validator.New()

	if err := validate.Var(model.ApplicationID, "required;gt=1"); err != nil {
		fields := []string{"applicationId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid applicationId ", err, nil)
	}
	appModel := models.Application{}
	if err := s.db.Model(&appModel).Where("id = ?", model.ApplicationID).First(&appModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"applicationId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid applicationId ", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if err := validate.Var(model.AppEnvironmentID, "required;gt=1"); err != nil {
		fields := []string{"appEnvironmentId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid appEnvironmentId ", err, nil)
	}
	appEnvModel := models.AppEnvironment{}
	if err := s.db.Model(&appEnvModel).Where("id = ?", model.AppEnvironmentID).First(&appEnvModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"appEnvironmentId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid appEnvironmentId ", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if err := validate.Var(model.Name, "required;gt=1"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}

	if err := validate.Var(model.Uri, "required;gt=1"); err != nil {
		fields := []string{"uri"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid uri ", err, nil)
	}

	if err := validate.Var(model.Active, "required"); err != nil {
		fields := []string{"active"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid active ", err, nil)
	}

	return nil
}

func (s *AppRoutesService) assign(to *models.AppRoute, from *models.AppRoute, operation payload.Operation) {

	now := time.Now().UTC()

	if operation == payload.SVC_OPERATION_CREATE {

		to.CreatedBy = s.request.Session.Auth.UserSession.ID
		to.CreatedAt = now

	} else if operation == payload.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.request.Session.Auth.UserSession.ID
		to.DeletedAt = &now

	} else {

		to.ApplicationID = from.ApplicationID
		to.AppEnvironmentID = from.AppEnvironmentID
		to.Name = from.Name
		to.Description = from.Description
		to.Uri = from.Uri
		to.Active = from.Active
	}

	to.UpdatedBy = s.request.Session.Auth.UserSession.ID
	to.UpdatedAt = now
}
