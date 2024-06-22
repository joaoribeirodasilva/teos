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

type UserPermissionsService struct {
	payload *payload.Payload
	db      *gorm.DB
	request *payload.HttpRequest
	history *redisdb.RedisDB
}

func NewUserPermissionsService(payload *payload.Payload) *UserPermissionsService {
	return &UserPermissionsService{
		payload: payload,
		db:      payload.Services.Db.GetDatabase(),
		request: payload.Http.Request,
		history: payload.Services.HistoryDb,
	}

}

// List returns a list of users from the collection
func (s *UserPermissionsService) List(filter string, args ...any) (*models.UserPermissions, *logger.HttpError) {

	model := models.UserPermission{}
	models := models.UserPermissions{}

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
func (s *UserPermissionsService) Get(model *models.UserPermission, filter string, args ...any) *logger.HttpError {

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
func (s *UserPermissionsService) Create(model *models.UserPermission) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	//TODO: organization config any can create user or only the organization

	exists := models.UserPermission{}

	if err := s.db.Where(
		"organization_id = ? AND auth_role_id = ? AND app_route_id = ? user_id = ?",
		model.OrganizationID,
		model.AuthRoleID,
		model.AppRouteID,
		model.UserID,
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
func (s *UserPermissionsService) Update(model *models.UserPermission) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	// Security
	// TODO: security

	exists := models.UserPermission{}
	if err := s.db.Where(
		"organization_id = ? AND auth_role_id = ? AND app_route_id = ? user_id = ?",
		model.OrganizationID,
		model.AuthRoleID,
		model.AppRouteID,
		model.UserID,
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
func (s *UserPermissionsService) Delete(id uint) *logger.HttpError {

	exists := &models.UserPermission{}

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

	s.assign(exists, nil, payload.SVC_OPERATION_UPDATE)

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

func (s *UserPermissionsService) Validate(model *models.UserPermission) *logger.HttpError {

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

	if err := validate.Var(model.AuthRoleID, "required,gte=1"); err != nil {
		fields := []string{"authRoleId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid authRoleId ", err, nil)
	}

	authRoleModel := models.AuthRole{}
	if err := s.db.Model(&authRoleModel).Where("id = ?", model.AuthRoleID).First(&authRoleModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"authRoleId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid authRoleId ", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if err := validate.Var(model.AppRouteID, "required,gte=1"); err != nil {
		fields := []string{"authRoleId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid authRoleId ", err, nil)
	}

	appRouteModel := models.AppRoute{}
	if err := s.db.Model(&appRouteModel).Where("id = ?", model.AppRouteID).First(&appRouteModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"appRouteId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid appRouteId ", err, nil)
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

func (s *UserPermissionsService) assign(to *models.UserPermission, from *models.UserPermission, operation payload.Operation) {

	now := time.Now().UTC()

	if operation == payload.SVC_OPERATION_CREATE {

		to.CreatedBy = s.request.Session.Auth.UserSession.UserID
		to.CreatedAt = now

	} else if operation == payload.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.request.Session.Auth.UserSession.UserID
		to.DeletedAt = &now

	} else {

		to.OrganizationID = from.OrganizationID
		to.AuthRoleID = from.AuthRoleID
		to.AppRouteID = from.AppRouteID
		to.UserID = from.UserID
		to.Active = from.Active
	}

	to.UpdatedBy = s.request.Session.Auth.UserSession.UserID
	to.UpdatedAt = now
}
