package authentication

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/utils/cookie"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
)

type Auth struct {
	config         *configuration.Configuration
	sessionDB      *redisdb.RedisDB
	userID         uint
	organizationID uint
	sessionID      uint
	email          string
	firstName      string
	surname        string
	avatarUrl      string
}

// Change to User
func New(config *configuration.Configuration, sessionDB *redisdb.RedisDB, userID uint, organizationID uint, sessionID uint, email string, firstName string, surname string, avatarUrl string) *Auth {

	return &Auth{
		config:         config,
		sessionDB:      sessionDB,
		userID:         userID,
		organizationID: organizationID,
		sessionID:      sessionID,
		email:          email,
		firstName:      firstName,
		surname:        surname,
		avatarUrl:      avatarUrl,
	}
}

func NewFromToken(config *configuration.Configuration, sessionDB *redisdb.RedisDB, user *token.User) *Auth {

	return &Auth{
		config:         config,
		sessionDB:      sessionDB,
		userID:         user.ID,
		organizationID: user.OrganizationID,
		sessionID:      user.SessionID,
		email:          user.Email,
		firstName:      user.Name,
		surname:        user.Surname,
		avatarUrl:      user.AvatarUrl,
	}
}

func (a *Auth) SetToken(c *gin.Context) error {

	tokenUser := token.User{
		ID:             a.userID,
		OrganizationID: a.organizationID,
		SessionID:      a.sessionID,
		Email:          a.email,
		Name:           a.firstName,
		Surname:        a.surname,
		AvatarUrl:      a.avatarUrl,
	}

	tk := token.New(a.config)
	if err := tk.Create(&tokenUser); err != nil {
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to create session token", err, nil)
	}

	ck := cookie.New(a.config, *tk.TokenString)
	if err := ck.Create(); err != nil {
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to create session cookie", err, nil)
	}

	sessionTTL, err := a.config.GetInt("AUTH_COOKIE_EXPIRE")
	if err != nil {
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to cookie ttl", err, nil)
	}

	stringId := strconv.Itoa(int(sessionTTL))
	if err := a.sessionDB.Set(stringId, *tk.TokenString, int(sessionTTL)); err != nil {
		return logger.Error(logger.LogStatusInternalServerError, nil, "failed to save session", err, nil)
	}

	ck.Set(c)

	return nil
}
