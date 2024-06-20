package sessions

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
	"github.com/joaoribeirodasilva/teos/hist/services/histories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

const (
	collectionName = "user_sessions"
)

type SessionsService struct {
	services      *structures.RequestValues
	db            *gorm.DB
	user          *token.User
	query         *requests.QueryString
	sessionDb     *redisdb.RedisDB
	permissionsDb *redisdb.RedisDB
	historyDb     *redisdb.RedisDB
}

func New(services *structures.RequestValues) *SessionsService {
	s := &SessionsService{}
	s.services = services
	s.db = services.Services.Db.GetDatabase()
	s.user = services.User
	s.query = &services.Query
	s.sessionDb = services.Services.SessionsDB
	s.permissionsDb = services.Services.PermissionsDB
	s.historyDb = services.Services.HistoryDB
	return s
}

// List returns a list of user sessions from the collection
func (s *SessionsService) List(filter string, args ...any) (*models.Users, *logger.HttpError) {

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

// Get returns a single user sessions from the collection
func (s *SessionsService) Get(model *models.User, filter string, args ...any) *logger.HttpError {

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

// Create creates a new user session document or returns a logger.HttpError in case of error
func (s *SessionsService) Create(model *models.UserSession) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	s.assign(
		nil,
		model,
		services.SVC_OPERATION_CREATE,
	)

	if _, err := s.coll.InsertOne(s.context, model); err != nil {

		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to create document",
			err,
			nil,
		)
	}

	return nil
}

// Create updates a user session document or returns a logger.HttpError in case of error
func (s *SessionsService) Update(model *models.UserSessionModel) *logger.HttpError {

	if err := s.Validate(model); err != nil {
		return err
	}

	exists := &models.UserSessionModel{}

	if err := s.Get(
		bson.D{
			{Key: "_id", Value: s.query.ID},
		},
		exists,
	); err != nil {

		return err
	}

	historySvc := histories.New(s.services)
	history := &models.HistHistoryModel{
		AppAppID:   s.services.Services.Configuration.GetAppID(),
		Collection: collectionName,
		OriginalID: exists.ID,
		Data:       exists,
	}
	if err := historySvc.Create(history); err != nil {
		return err
	}

	s.assign(
		exists,
		model,
		services.SVC_OPERATION_UPDATE,
	)

	if _, err := s.coll.UpdateOne(
		s.context,
		bson.D{{Key: "_id", Value: s.query.ID}},
		bson.D{{Key: "$set", Value: model}},
	); err != nil {
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to update document", err, nil)
	}

	return nil
}

// Delete deletes a user session document or returns a logger.HttpError in case of error
func (s *SessionsService) Delete(model *models.UserSessionModel) *logger.HttpError {

	exists := &models.UserSessionModel{}

	if err := s.Get(
		bson.D{
			{Key: "_id", Value: s.query.ID},
		},
		exists,
	); err != nil {

		return err
	}

	historySvc := histories.New(s.services)
	history := &models.HistHistoryModel{
		AppAppID:   s.services.Services.Configuration.GetAppID(),
		Collection: collectionName,
		OriginalID: exists.ID,
		Data:       exists,
	}
	if err := historySvc.Create(history); err != nil {
		return err
	}

	s.assign(
		exists,
		model,
		services.SVC_OPERATION_DELETE,
	)

	if _, err := s.coll.UpdateOne(
		s.context,
		bson.D{{Key: "_id", Value: s.query.ID}},
		bson.D{{Key: "$set", Value: model}},
	); err != nil {
		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to update document",
			err,
			nil,
		)
	}

	return nil
}

func (m *SessionsService) Validate(model *models.UserSessionModel) *logger.HttpError {

	validate := validator.New()

	if err := validate.Var(model.UserUserID, "required"); err != nil {
		fields := []string{"userUserId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid userUserId ", err, nil)
	}

	// TODO: Validate related
	/* 	user := NewUserUserModel(m.ctx)
	   	if appErr := m.FindByID(m.UserUserID, user); appErr != nil {
	   		return appErr
	   	} */

	return nil
}

func (s *SessionsService) assign(to *models.UserSessionModel, from *models.UserSessionModel, operation services.Operation) {

	now := time.Now().UTC()

	if operation == services.SVC_OPERATION_CREATE {

		to.ID = primitive.NewObjectID()
		to.CreatedBy = s.user.ID
		to.CreatedAt = now

	} else if operation == services.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.user.ID
		to.DeletedAt = &now

	} else {

		to.UserUserID = from.UserUserID
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
