package app_route_methods

import (
	"errors"
	"fmt"
	"slices"
	"strings"
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

var (
	methods = []models.Methods{"GET", "POST", "PUT", "PATCH", "DELETE"}
)

type AppRouteMethodsService struct {
	services      *structures.RequestValues
	db            *gorm.DB
	user          *token.User
	query         *requests.QueryString
	sessionDb     *redisdb.RedisDB
	permissionsDb *redisdb.RedisDB
	historyDb     *redisdb.RedisDB
}

func New(services *structures.RequestValues) *AppRouteMethodsService {
	s := &AppRouteMethodsService{}
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
func (s *AppRouteMethodsService) List(filter string, args ...any) (*models.AppRouteMethods, *logger.HttpError) {

	model := models.AppRouteMethod{}
	models := models.AppRouteMethods{}

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
func (s *AppRouteMethodsService) Get(model *models.AppRouteMethod, filter string, args ...any) *logger.HttpError {

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
func (s *AppRouteMethodsService) Create(model *models.AppRouteMethod) *logger.HttpError {

	if s.user.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to create this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
	}

	if err := s.Validate(model); err != nil {
		return err
	}

	//TODO: organization config any can create user or only the organization

	exists := models.AppRouteMethod{}

	if err := s.db.Where(
		"app_route_id = ? AND ((name = ?) OR (method_id = ? AND uri = ?) ",
		model.AppRouteID,
		model.Name,
		model.Method,
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
func (s *AppRouteMethodsService) Update(model *models.AppRouteMethod) *logger.HttpError {

	if s.user.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to update this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
	}

	if err := s.Validate(model); err != nil {
		return err
	}

	// Security
	// TODO: security

	exists := models.AppRouteMethod{}
	if err := s.db.Where(
		"app_route_id = ? AND ((name = ?) OR (method_id = ? AND uri = ?) ",
		model.AppRouteID,
		model.Name,
		model.Method,
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
func (s *AppRouteMethodsService) Delete(id uint) *logger.HttpError {

	if s.user.OrganizationID != 1 {
		err := errors.New("the current user does not have permission to delete this record")
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusUnauthorized, &fields, "user not authorized", err, nil)
	}

	exists := &models.AppRouteMethod{}

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

func (s *AppRouteMethodsService) Validate(model *models.AppRouteMethod) *logger.HttpError {

	validate := validator.New()

	if err := validate.Var(model.AppRouteID, "required,gt=0"); err != nil {
		fields := []string{"appRouteId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid appRouteId ", err, nil)
	}

	appRouteModel := models.AppRoute{}
	if err := s.db.Model(&appRouteModel).Where("id = ?", model.AppRouteID).First(&appRouteModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"appRouteId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid appRouteId ", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if err := validate.Var(model.Name, "required;gt=1"); err != nil {
		fields := []string{"name"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid name ", err, nil)
	}

	model.Method = models.Methods(strings.ToUpper(string(model.Method)))
	if err := validate.Var(model.Method, "required;gt=1"); err != nil {
		fields := []string{"method"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid method ", err, nil)
	}

	if !slices.Contains(methods, model.Method) {
		err := errors.New("method not found")
		fields := []string{"method"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid method ", err, nil)
	}

	if err := validate.Var(model.Active, "required"); err != nil {
		fields := []string{"active"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid active ", err, nil)
	}

	return nil
}

func (s *AppRouteMethodsService) assign(to *models.AppRouteMethod, from *models.AppRouteMethod, operation services.Operation) {

	now := time.Now().UTC()

	if operation == services.SVC_OPERATION_CREATE {

		to.CreatedBy = s.user.ID
		to.CreatedAt = now

	} else if operation == services.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.user.ID
		to.DeletedAt = &now

	} else {

		to.AppRouteID = from.AppRouteID
		to.Name = from.Name
		to.Description = from.Description
		to.Method = from.Method
		to.Uri = from.Uri
		to.Open = from.Open
		to.Active = from.Active
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
