package user_organizations

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

type UserOrganizationsService struct {
	services      *structures.RequestValues
	db            *gorm.DB
	user          *token.User
	query         *requests.QueryString
	sessionDb     *redisdb.RedisDB
	permissionsDb *redisdb.RedisDB
	historyDb     *redisdb.RedisDB
}

func New(services *structures.RequestValues) *UserOrganizationsService {
	s := &UserOrganizationsService{}
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
func (s *UserOrganizationsService) List(filter string, args ...any) (*models.UserOrganizations, *logger.HttpError) {

	model := models.UserOrganization{}
	models := models.UserOrganizations{}

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
func (s *UserOrganizationsService) Get(model *models.UserOrganization, filter string, args ...any) *logger.HttpError {

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
func (s *UserOrganizationsService) Create(model *models.UserOrganization) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	//TODO: organization config any can create user or only the organization

	exists := models.UserOrganization{}

	if err := s.db.Where("organization_id = ? AND user_id = ?", model.OrganizationID, model.UserID).First(&exists).Error; err != nil {

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
func (s *UserOrganizationsService) Update(model *models.UserOrganization) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	// Security
	// TODO: security

	exists := models.UserOrganization{}
	if err := s.db.Where("organization_id = ? AND user_id = ?", model.OrganizationID, model.UserID).First(&exists).Error; err != nil {

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
func (s *UserOrganizationsService) Delete(id uint) *logger.HttpError {

	exists := &models.UserOrganization{}

	// Security
	// TODO: security

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

func (s *UserOrganizationsService) Validate(model *models.UserOrganization) *logger.HttpError {

	validate := validator.New()

	if err := validate.Var(model.OrganizationID, "required,gte=1"); err != nil {
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid organizationId ", err, nil)
	}

	orgModel := models.Organization{}
	if err := s.db.Model(&orgModel).Where("id = ?", model.OrganizationID).First(&orgModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"organizationId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid organizationId ", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	userModel := models.User{}
	if err := s.db.Model(&userModel).Where("id = ?", model.UserID).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"userId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid userId ", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if err := validate.Var(model.Active, "required"); err != nil {
		fields := []string{"active"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid active ", err, nil)
	}

	return nil
}

func (s *UserOrganizationsService) assign(to *models.UserOrganization, from *models.UserOrganization, operation services.Operation) {

	now := time.Now().UTC()

	if operation == services.SVC_OPERATION_CREATE {

		to.CreatedBy = s.user.ID
		to.CreatedAt = now

	} else if operation == services.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.user.ID
		to.DeletedAt = &now

	} else {

		to.OrganizationID = from.OrganizationID
		to.UserID = from.UserID
		to.Active = from.Active
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
