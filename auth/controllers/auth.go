package controllers

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaoribeirodasilva/teos/auth/requests"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/responses"
	"github.com/joaoribeirodasilva/teos/common/utils/cookie"
	"github.com/joaoribeirodasilva/teos/common/utils/password"
	"github.com/joaoribeirodasilva/teos/users/models"
	"gorm.io/gorm"
)

const (
	AUTH_COOKIE_NAME      = "auth"
	AUTH_COOKIE_EXPIRE    = 900
	AUTH_COOKIE_DOMAIN    = "localhost"
	AUTH_COOKIE_HTTP_ONLY = true
	AUTH_COOKIE_SECURE    = false
)

func AuthLogin(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	email, passwd, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	record := models.UserUser{}
	if err := vars.Db.Conn.Where("email = ?", email).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := responses.ResponseErrorMessage{
				Error: responses.ErrMessage{
					Message: "account not found",
				},
			}
			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	if record.Password == nil {
		response := responses.ResponseErrorMessage{
			Error: responses.ErrMessage{
				Message: "no password set, please reset your password",
			},
		}
		c.AbortWithStatusJSON(http.StatusNotFound, response)
	}

	if !password.Check(passwd, *record.Password) {
		response := responses.ResponseErrorMessage{
			Error: responses.ErrMessage{
				Message: "invalid user and/or password",
			},
		}
		c.AbortWithStatusJSON(http.StatusForbidden, response)
	}

	if record.Active == 0 {
		response := responses.ResponseErrorMessage{
			Error: responses.ErrMessage{
				Message: "account disabled",
			},
		}
		c.AbortWithStatusJSON(http.StatusForbidden, response)
	}

	tokenString := ""

	cookie, err := cookie.NewFromConfiguration(
		tokenString,
		vars.Configuration,
	)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	//TODO: Generate the token

	cookie.SetCookie(c)

	c.Status(http.StatusOK)
}

func AuthForgot(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	request := requests.ForgotPassword{}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	_, err = mail.ParseAddress(request.Email)
	if err != nil {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "email",
				Message: "invalid email address",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	record := models.UserUser{}
	if err := vars.Db.Conn.Where("email = ?", request.Email).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := responses.ResponseErrorMessage{
				Error: responses.ErrMessage{
					Message: "account not found",
				},
			}
			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	uuid.EnableRandPool()
	reset_key := uuid.NewString()

	expire := time.Now()
	expire = expire.Add(time.Hour * 24)

	reset := models.UserReset{
		UserResetTypeID: 1,
		UserUserID:      record.ID,
		ResetKey:        reset_key,
		Used:            nil,
		Expire:          expire.UTC(),
		CreatedBy:       1,
		UpdatedBy:       1,
	}

	if err := vars.Db.Conn.Create(&reset).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//TODO: send email with the request link

	c.Status(http.StatusCreated)

}

func AuthReset(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	key := c.Query("key")
	if key == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	request := requests.ResetPassword{}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	request.Password = strings.TrimSpace(request.Password)
	if request.Password == "" || len(request.Password) < 6 {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "password",
				Message: "the password must have at least 6 characters",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	request.CheckPassword = strings.TrimSpace(request.CheckPassword)
	if request.CheckPassword != request.Password {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "password",
				Message: "the passwords don't match",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	reset := models.UserReset{}
	if err := vars.Db.Conn.Where("reset_key = ?", key).First(&reset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	if now.After(reset.Expire) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// TODO: move password validation and generation into a separate
	//       code to use here and on user create and update
	// TODO: Get the user account
	// TODO: Update the user password

	user := models.UserUser{}

	if err := vars.Db.Conn.Where("id = ?", reset.UserUserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user.Active == 0 {
		response := responses.ResponseErrorMessage{
			Error: responses.ErrMessage{
				Message: "account disabled",
			},
		}
		c.AbortWithStatusJSON(http.StatusForbidden, response)
	}

	tempPassword, err := password.Hash(request.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user.Password = &tempPassword
	user.UpdatedBy = user.ID

	if err := vars.Db.Conn.Save(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)

}

func AuthLogout(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//TODO: delete the user

	c.SetCookie(AUTH_COOKIE_NAME, "", AUTH_COOKIE_EXPIRE, "/", AUTH_COOKIE_DOMAIN, AUTH_COOKIE_HTTP_ONLY, AUTH_COOKIE_SECURE)

	c.Status(http.StatusOK)
}
