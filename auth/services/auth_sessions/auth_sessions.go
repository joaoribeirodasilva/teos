package auth_sessions

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joaoribeirodasilva/teos/common/authentication"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/requests"
	"github.com/joaoribeirodasilva/teos/common/services"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"github.com/joaoribeirodasilva/teos/common/utils/password"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
	"gorm.io/gorm"
)

type AuthSessionsService struct {
	services      *structures.RequestValues
	db            *gorm.DB
	user          *token.User
	query         *requests.QueryString
	sessionDb     *redisdb.RedisDB
	permissionsDb *redisdb.RedisDB
	historyDb     *redisdb.RedisDB
}

func New(services *structures.RequestValues) *AuthSessionsService {
	s := &AuthSessionsService{}
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
func (s *AuthSessionsService) Login(email string, passwd string) (*authentication.Auth, *logger.HttpError) {

	user := models.User{}

	if err := s.db.Model(&user).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, logger.Error(logger.LogStatusForbidden, nil, "wrong username or password", err, nil)
		}
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if !password.Check(passwd, user.Password) {
		err := errors.New("unknown user")
		return nil, logger.Error(logger.LogStatusForbidden, nil, "wrong username or password", err, nil)
	}

	if user.EmailVerified == nil {
		err := errors.New("account not active")
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "no email verification date", err, nil)
	}

	userOrg := models.UserOrganization{}
	if err := s.db.Model(&userOrg).Where("user_id = ?", user.ID).First(&userOrg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, logger.Error(logger.LogStatusForbidden, nil, "user does not belong to the organization", err, nil)
		}
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
	}

	if userOrg.Active == 0 {
		err := errors.New("account not active")
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "user is not active in the organization", err, nil)
	}

	session := models.AuthSession{
		OrganizationID: userOrg.ID,
		UserID:         user.ID,
	}

	if err := s.Validate(&session); err != nil {
		return nil, err
	}

	s.assign(&session, nil, services.SVC_OPERATION_CREATE)

	if err := s.db.Create(&session).Error; err != nil {
		return nil, logger.Error(logger.LogStatusInternalServerError, nil, "failed to create session in the database", err, nil)
	}

	auth := authentication.New(s.services.Services.Configuration, s.sessionDb, user.ID, userOrg.ID, session.ID, user.Email, user.FirstName, user.Surname, user.AvatarUrl)

	return auth, nil
}

// Delete deletes a user document or returns a logger.HttpError in case of error
func (s *AuthSessionsService) Delete(id uint) *logger.HttpError {

	//LOGOUT

	exists := &models.AuthSession{}
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

	strId := strconv.Itoa(int(id))
	if err := s.sessionDb.Del(strId); err != nil {
		return logger.Error(
			logger.LogStatusInternalServerError,
			nil,
			"failed to delete session from redis",
			err,
			nil,
		)
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

func (s *AuthSessionsService) assign(to *models.AuthSession, from *models.AuthSession, operation services.Operation) {

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
	}

	to.UpdatedBy = s.user.ID
	to.UpdatedAt = now
}
