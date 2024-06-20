package app_configurations

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

type AppConfigurationsService struct {
	services      *structures.RequestValues
	db            *gorm.DB
	user          *token.User
	query         *requests.QueryString
	sessionDb     *redisdb.RedisDB
	permissionsDb *redisdb.RedisDB
	historyDb     *redisdb.RedisDB
}

func New(services *structures.RequestValues) *AppConfigurationsService {
	s := &AppConfigurationsService{}
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
func (s *AppConfigurationsService) List(filter string, args ...any) (*models.AppConfigurations, *logger.HttpError) {

	model := models.AppConfiguration{}
	models := models.AppConfigurations{}

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
func (s *AppConfigurationsService) Get(model *models.AppConfiguration, filter string, args ...any) *logger.HttpError {

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

	return nil
}

// Create creates a new user document or returns a logger.HttpError in case of error
func (s *AppConfigurationsService) Create(model *models.AppConfiguration) *logger.HttpError {

	if s.user.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to create this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
	}

	if err := s.Validate(model); err != nil {
		return err
	}

	//TODO: organization config any can create user or only the organization

	exists := models.AppConfiguration{}

	if err := s.db.Where(
		"app_environment_id = ? AND application_id = ? AND configuration_key = ?",
		model.AppEnvironmentID,
		model.ApplicationID,
		model.ConfigurationKey,
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
func (s *AppConfigurationsService) Update(model *models.AppConfiguration) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	// Security
	if s.user.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to update this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
	}

	exists := models.AppConfiguration{}
	if err := s.db.Where(
		"app_environment_id = ? AND application_id = ? AND configuration_key = ?",
		model.AppEnvironmentID,
		model.ApplicationID,
		model.ConfigurationKey,
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
func (s *AppConfigurationsService) Delete(id uint) *logger.HttpError {

	exists := &models.AppConfiguration{}

	// Security
	if s.user.OrganizationID != 1 {
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

func (s *AppConfigurationsService) Validate(model *models.AppConfiguration) *logger.HttpError {

	validate := validator.New()

	if err := validate.Var(model.AppEnvironmentID, "required,gt=0"); err != nil {
		fields := []string{"appEnvironmentId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid appEnvironmentId ", err, nil)
	}
	appEnvironmentModel := models.AppEnvironment{}
	if err := s.db.Model(&appEnvironmentModel).Where("id = ?", model.AppEnvironmentID).First(&appEnvironmentModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"appEnvironmentId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid appEnvironmentId ", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if err := validate.Var(model.ApplicationID, "required,gt=0"); err != nil {
		fields := []string{"applicationId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid applicationId ", err, nil)
	}
	applicationModel := models.Application{}
	if err := s.db.Model(&applicationModel).Where("id = ?", model.ApplicationID).First(&applicationModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"applicationId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid applicationId ", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if err := validate.Var(model.ConfigurationKey, "required;gt=1"); err != nil {
		fields := []string{"configurationKey"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid configurationKey ", err, nil)
	}

	return nil
}

func (s *AppConfigurationsService) assign(to *models.AppConfiguration, from *models.AppConfiguration, operation services.Operation) {

	now := time.Now().UTC()

	if operation == services.SVC_OPERATION_CREATE {

		to.CreatedBy = s.user.ID
		to.CreatedAt = now

	} else if operation == services.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.user.ID
		to.DeletedAt = &now

	} else {

		to.AppEnvironmentID = from.AppEnvironmentID
		to.ApplicationID = from.ApplicationID
		to.ConfigurationKey = from.ConfigurationKey
		to.Type = from.Type
		to.ValString = from.ValString
		to.ValInt = from.ValInt
		to.ValDouble = from.ValDouble
		to.ValDate = from.ValDate
		to.ValTime = from.ValTime
		to.ValDateTime = from.ValDateTime
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
