package services

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/payload"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"gorm.io/gorm"
)

type HistoryService struct {
	payload *payload.Payload
	db      *gorm.DB
	request *payload.HttpRequest
	history *redisdb.RedisDB
}

func NewHistoryService(payload *payload.Payload) *HistoryService {
	return &HistoryService{
		payload: payload,
		db:      payload.Services.Db.GetDatabase(),
		request: payload.Http.Request,
		history: payload.Services.HistoryDb,
	}
}

// List returns a list of users from the collection
func (s *HistoryService) List(filter string, args ...any) (*models.Histories, *logger.HttpError) {

	model := models.History{}
	models := models.Histories{}

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
func (s *HistoryService) Get(model *models.History, filter string, args ...any) *logger.HttpError {

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
func (s *HistoryService) Create(model *models.History) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	s.assign(model, nil, payload.SVC_OPERATION_CREATE)

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

func (s *HistoryService) Validate(model *models.History) *logger.HttpError {

	validate := validator.New()

	if err := validate.Var(model.OrganizationID, "required;gt=1"); err != nil {
		fields := []string{"organizationId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid organizationId ", err, nil)
	}

	if err := validate.Var(model.Tablename, "required;gt=1"); err != nil {
		fields := []string{"tableName"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid tableName ", err, nil)
	}

	if err := validate.Var(model.OriginalID, "required;gt=0"); err != nil {
		fields := []string{"originalId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid originalId ", err, nil)
	}

	if err := validate.Var(model.JsonData, "required;gt=1"); err != nil {
		fields := []string{"jsonData"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid jsonData ", err, nil)
	}

	return nil
}

func (s *HistoryService) assign(to *models.History, from *models.History, operation payload.Operation) {

	now := time.Now().UTC()

	if operation == payload.SVC_OPERATION_CREATE {

		to.CreatedBy = s.request.Session.Auth.UserSession.UserID
		to.CreatedAt = now

	} else if operation == payload.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.request.Session.Auth.UserSession.UserID
		to.DeletedAt = &now

	} else {

		to.OrganizationID = from.OrganizationID
		to.Tablename = from.Tablename
		to.OriginalID = from.OriginalID
		to.JsonData = from.JsonData
	}

	to.UpdatedBy = s.request.Session.Auth.UserSession.UserID
	to.UpdatedAt = now
}
