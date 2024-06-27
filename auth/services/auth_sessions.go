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
	"github.com/joaoribeirodasilva/teos/common/utils/password"
	"gorm.io/gorm"
)

type AuthSessionsService struct {
	payload *payload.Payload
	db      *gorm.DB
	request *payload.HttpRequest
	history *redisdb.RedisDB
}

func NewAuthSessionsService(payload *payload.Payload) *AuthSessionsService {
	return &AuthSessionsService{
		payload: payload,
		db:      payload.Services.Db.GetDatabase(),
		request: payload.Http.Request,
		history: payload.Services.HistoryDb,
	}
}

// List returns a list of users from the collection
func (s *AuthSessionsService) List(filter string, args ...any) (*models.AuthSessions, *logger.HttpError) {

	model := models.AuthSession{}
	models := models.AuthSessions{}

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
func (s *AuthSessionsService) Get(model *models.AuthSession, filter string, args ...any) *logger.HttpError {

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
func (s *AuthSessionsService) Login(email string, passwd string) *logger.HttpError {

	user := models.User{}

	if err := s.db.Model(&user).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return logger.Error(logger.LogStatusForbidden, nil, "wrong username or password", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if !password.Check(passwd, user.Password) {
		err := errors.New("unknown user")
		return logger.Error(logger.LogStatusForbidden, nil, "wrong username or password", err, nil)
	}

	if user.EmailVerified == nil {
		err := errors.New("account not active")
		return logger.Error(logger.LogStatusInternalServerError, nil, "no email verification date", err, nil)
	}

	userOrg := models.UserOrganization{}
	if err := s.db.Model(&userOrg).Where("user_id = ?", user.ID).First(&userOrg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return logger.Error(logger.LogStatusForbidden, nil, "user does not belong to an organization", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if (userOrg.Active == 0 || user.Active == 0) && userOrg.OrganizationID != 1 {
		err := errors.New("account not active")
		return logger.Error(logger.LogStatusInternalServerError, nil, "user is not active in the organization", err, nil)
	}

	session := models.AuthSession{
		OrganizationID: userOrg.OrganizationID,
		UserID:         user.ID,
	}

	if err := s.Validate(&session); err != nil {
		return err
	}

	s.assign(&session, nil, payload.SVC_OPERATION_CREATE)

	if err := s.db.Create(&session).Error; err != nil {
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to create session in the database", err, nil)
	}

	s.request.Session.Auth.UserSession = &payload.SessionAuth{
		ID:             session.ID,
		UserID:         session.UserID,
		OrganizationID: session.OrganizationID,
		Email:          user.Email,
		Name:           user.FirstName,
		Surname:        user.Surname,
		AvatarUrl:      user.AvatarUrl,
	}

	s.payload.Services.SessionsDb.Set(
		fmt.Sprintf("%d", s.request.Session.Auth.UserSession.ID),
		s.request.Session.Auth.UserSession,
		s.request.Session.Auth.Ttl,
	)

	//auth := authentication.New(s.payload.Services.Configuration, s.sessionDb, user.ID, userOrg.ID, session.ID, user.Email, user.FirstName, user.Surname, user.AvatarUrl)

	return nil
}

func (s *AuthSessionsService) Forgot(email string) *logger.HttpError {
	//TODO:
	return nil
}

func (s *AuthSessionsService) Reset(password string, passwordConf string) *logger.HttpError {
	//TODO:
	return nil
}

// Delete deletes a user document or returns a logger.HttpError in case of error
func (s *AuthSessionsService) Logout() *logger.HttpError {

	//LOGOUT

	exists := &models.AuthSession{}
	if err := s.db.Where("id = ?", s.payload.Http.Request.Session.Auth.UserSession.ID).First(&exists).Error; err != nil {

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

	id := s.payload.Http.Request.Session.Auth.UserSession.ID
	if err := s.payload.Services.SessionsDb.Del(fmt.Sprintf("%d", id)); err != nil {
		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to delete session from redis",
			err,
			nil,
		)
	}

	now := time.Now().UTC()
	exists.UpdatedBy = s.payload.Http.Request.Session.Auth.UserSession.UserID
	exists.UpdatedAt = now
	exists.DeletedBy = &s.payload.Http.Request.Session.Auth.UserSession.UserID
	exists.DeletedAt = &now

	if err := s.db.Model(&exists).Save(exists).Error; err != nil {

		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to save document into database",
			err,
			nil,
		)
	}

	s.payload.Http.Request.Session.Auth.UserSession = payload.EmptySession()

	return nil
}

func (s *AuthSessionsService) Validate(model *models.AuthSession) *logger.HttpError {

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

	if err := validate.Var(model.UserID, "required,gte=1"); err != nil {
		fields := []string{"userId"}
		return logger.Error(logger.LogStatusBadRequest, &fields, "invalid userId ", err, nil)
	}
	userModel := models.User{}
	if err := s.db.Model(&userModel).Where("id = ?", model.UserID).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fields := []string{"userId"}
			return logger.Error(logger.LogStatusBadRequest, &fields, "invalid userId ", err, nil)
		}
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}
	return nil
}

func (s *AuthSessionsService) assign(to *models.AuthSession, from *models.AuthSession, operation payload.Operation) {

	now := time.Now().UTC()

	if operation == payload.SVC_OPERATION_CREATE {

		to.CreatedBy = s.request.Session.Auth.UserSession.UserID
		to.CreatedAt = now

	} else if operation == payload.SVC_OPERATION_DELETE {

		to.DeletedBy = &s.request.Session.Auth.UserSession.UserID
		to.DeletedAt = &now

	} else {

		to.OrganizationID = from.OrganizationID
		to.UserID = from.UserID
	}

	to.UpdatedBy = s.request.Session.Auth.UserSession.UserID
	to.UpdatedAt = now
}
