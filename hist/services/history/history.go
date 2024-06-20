package histories

import (
	"errors"
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

type HistoryService struct {
	services      *structures.RequestValues
	db            *gorm.DB
	user          *token.User
	query         *requests.QueryString
	sessionDb     *redisdb.RedisDB
	permissionsDb *redisdb.RedisDB
	historyDb     *redisdb.RedisDB
}

func New(services *structures.RequestValues) *HistoryService {
	s := &HistoryService{}
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
func (s *HistoryService) Create(model *models.History) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	s.assign(model, nil, services.SVC_OPERATION_CREATE)

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

func (s *HistoryService) assign(to *models.History, from *models.History, operation services.Operation) {

	now := time.Now().UTC()

	if operation == services.SVC_OPERATION_CREATE {

		to.CreatedBy = s.user.ID
		to.CreatedAt = now

	} else if operation == services.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.user.ID
		to.DeletedAt = &now

	} else {

		to.OrganizationID = from.OrganizationID
		to.Tablename = from.Tablename
		to.OriginalID = from.OriginalID
		to.JsonData = from.JsonData
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
