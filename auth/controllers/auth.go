package controllers

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaoribeirodasilva/teos/auth/requests"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/utils/password"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	AUTH_COOKIE_NAME      = "auth"
	AUTH_COOKIE_EXPIRE    = 900
	AUTH_COOKIE_DOMAIN    = "localhost"
	AUTH_COOKIE_HTTP_ONLY = true
	AUTH_COOKIE_SECURE    = false
)

func AuthLogin(c *gin.Context) {

	values, httpErr := controllers.MustGetAll(c)
	if httpErr != nil {
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	email, passwd, ok := c.Request.BasicAuth()
	if !ok {
		err := errors.New("invalid username or password")
		httpErr := logger.Error(logger.LogStatusBadRequest, nil, "server didn't receive authentication", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	record := models.UserUserModel{}

	coll := values.Services.Db.GetDatabase().Collection("user_users")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
			c.AbortWithStatusJSON(httpErr.Status, httpErr)
			return
		}
		httpErr := logger.Error(logger.LogStatusForbidden, nil, "invalid username or password", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	if record.Password == nil {
		err := errors.New("")
		httpErr := logger.Error(logger.LogStatusForbidden, nil, "invalid password", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	if !password.Check(passwd, *record.Password) {
		err := errors.New("")
		httpErr := logger.Error(logger.LogStatusForbidden, nil, "invalid username or password", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	if record.Active == 0 || record.DeletedBy != nil || record.DeletedAt != nil {
		err := errors.New("")
		httpErr := logger.Error(logger.LogStatusUnauthorized, nil, "account disabled", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	sessionRecord := models.UserSessionModel{}

	now := time.Now().UTC()
	sessionRecord.ID = primitive.NewObjectID()
	sessionRecord.UserUserID = record.ID
	sessionRecord.CreatedBy = record.ID
	sessionRecord.CreatedAt = now
	sessionRecord.UpdatedBy = record.ID
	sessionRecord.UpdatedAt = now

	collSession := values.Services.Db.GetDatabase().Collection("user_sessions")

	result, err := collSession.InsertOne(context.TODO(), sessionRecord)
	if err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "failed to insert session into the database", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	sessionId := result.InsertedID.(primitive.ObjectID)
	tokenObject := token.New(values.Services.Configuration)

	tokenUser := token.User{
		ID:        record.ID,
		SessionID: sessionRecord.ID,
		Email:     record.Email,
		Name:      record.FirstName,
		Surname:   record.Surname,
	}

	if err := tokenObject.Create(&tokenUser, &sessionId); err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "failed to generate session token", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	cookieName, err := values.Services.Configuration.GetString("COOKIE_NAME")
	if err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "invalid cookie name", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	cookieExpire, err := values.Services.Configuration.GetInt("COOKIE_EXPIRE")
	if err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "invalid cookie expire", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	cookieDomain, err := values.Services.Configuration.GetString("COOKIE_DOMAIN")
	if err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "invalid cookie domain", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	cookieHttpOnly, err := values.Services.Configuration.GetBool("COOKIE_HTTP_ONLY")
	if err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "invalid cookie http only", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	cookieSecure, err := values.Services.Configuration.GetBool("COOKIE_SECURE")
	if err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "invalid cookie secure", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	expire := int(time.Now().Add(time.Second * time.Duration(cookieExpire)).Unix())
	// c.SetCookie(*tempName.String, tokenObject.TokenString, *tempExpire.Int, "/", *tempDomain.String, *tempSecure.Bool, *tempHttpOnly.Bool)
	c.SetCookie(cookieName, tokenObject.TokenString, expire, "", cookieDomain, cookieSecure, cookieHttpOnly)

	sessionDb := values.Services.SessionsDB
	sessionDb.Set(sessionRecord.ID.Hex(), tokenUser, int(cookieExpire))

	c.Status(http.StatusOK)
}

func AuthForgot(c *gin.Context) {

	values, httpErr := controllers.MustGetAll(c)
	if httpErr != nil {
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	request := requests.ForgotPassword{}
	if err := c.ShouldBind(&request); err != nil {
		fields := []string{"password"}
		httpErr := logger.Error(logger.LogStatusBadRequest, &fields, "no email provided", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	_, err := mail.ParseAddress(request.Email)
	if err != nil {
		fields := []string{"password"}
		httpErr := logger.Error(logger.LogStatusBadRequest, &fields, "invalid email provided", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	record := models.UserUserModel{}
	coll := values.Services.Db.GetDatabase().Collection("user_users")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: request.Email}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
			c.AbortWithStatusJSON(httpErr.Status, httpErr)
			return
		}
		httpErr := logger.Error(logger.LogStatusBadRequest, nil, "wrong username or password", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	if record.DeletedBy != nil || record.DeletedAt != nil {
		err := errors.New("")
		httpErr := logger.Error(logger.LogStatusUnauthorized, nil, "account disabled", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	uuid.EnableRandPool()
	resetKey := uuid.NewString()

	expire := time.Now()
	expire = expire.Add(time.Hour * 24)
	resetTypeId, err := primitive.ObjectIDFromHex("6669fdf9175f523b82a26a13")
	if err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "invalid reset type id", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	resetRecord := models.UserResetModel{
		UserResetTypeID: resetTypeId,
		UserUserID:      record.ID,
		ResetKey:        resetKey,
		Used:            nil,
		Expire:          expire,
	}

	collUserResets := values.Services.Db.GetDatabase().Collection("user_resets")
	if _, err := collUserResets.InsertOne(context.TODO(), resetRecord); err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "failed to insert password reset into database", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	//TODO: send email with the request link

	c.Status(http.StatusCreated)

}

func AuthReset(c *gin.Context) {

	values, httpErr := controllers.MustGetAll(c)
	if httpErr != nil {
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	key := c.Param("key")
	if key == "" {
		err := errors.New("")
		fields := []string{"key"}
		httpErr := logger.Error(logger.LogStatusBadRequest, &fields, "invalid reset key", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	request := requests.ResetPassword{}
	if err := c.ShouldBind(&request); err != nil {
		httpErr := logger.Error(logger.LogStatusBadRequest, nil, "invalid JSON body", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	request.Password = strings.TrimSpace(request.Password)
	if request.Password == "" || len(request.Password) < 6 {
		err := errors.New("")
		fields := []string{"password"}
		httpErr := logger.Error(logger.LogStatusBadRequest, &fields, "the password must have at least 6 characters", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	request.CheckPassword = strings.TrimSpace(request.CheckPassword)
	if request.CheckPassword != request.Password {
		err := errors.New("")
		fields := []string{"password"}
		httpErr := logger.Error(logger.LogStatusBadRequest, &fields, "the passwords don't match", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	record := models.UserResetModel{}

	coll := values.Services.Db.GetDatabase().Collection("user_resets")
	if err := coll.FindOne(context.TODO(), bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "resetKey", Value: key}},
			bson.D{{Key: "used", Value: nil}},
		}},
	}).Decode(&record); err != nil {
		if err == mongo.ErrNoDocuments {
			httpErr := logger.Error(logger.LogStatusBadRequest, nil, "the reset was not found or it expired", err, nil)
			c.AbortWithStatusJSON(httpErr.Status, httpErr)
			return
		}

		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	now := time.Now().UTC()
	if now.After(record.Expire) {
		err := errors.New("")
		httpErr := logger.Error(logger.LogStatusUnauthorized, nil, "the reset request has expired", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	userRecord := models.UserUserModel{}

	usersColl := values.Services.Db.GetDatabase().Collection("user_users")
	if err := usersColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: record.UserUserID}}).Decode(&userRecord); err != nil {
		if err == mongo.ErrNoDocuments {
			httpErr := logger.Error(logger.LogStatusNotFound, nil, "user not found", err, nil)
			c.AbortWithStatusJSON(httpErr.Status, httpErr)
			return
		}

		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "failed to query database", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	if userRecord.Active == 0 {
		err := errors.New("")
		httpErr := logger.Error(logger.LogStatusUnauthorized, nil, "account disabled", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	tempPassword, err := password.Hash(request.Password)
	if err != nil {
		httpErr := logger.Error(logger.LogStatusUnauthorized, nil, "account disabled", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	now = time.Now().UTC()
	userRecord.Password = &tempPassword
	userRecord.UpdatedBy = userRecord.ID
	userRecord.UpdatedAt = now

	_, err = usersColl.UpdateOne(context.TODO(), bson.D{
		{Key: "_id", Value: userRecord.ID},
	}, bson.D{
		{Key: "$set", Value: userRecord},
	})
	if err != nil {
		httpErr := logger.Error(logger.LogStatusUnauthorized, nil, "failed to update password into database", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	now = time.Now().UTC()
	record.Used = &now
	record.UpdatedBy = userRecord.ID
	record.UpdatedAt = now
	_, err = coll.UpdateOne(context.TODO(), bson.D{
		{Key: "_id", Value: record.ID},
	}, bson.D{
		{Key: "$set", Value: record},
	})
	if err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "failed to update password reset into database", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	c.Status(http.StatusOK)

}

func AuthLogout(c *gin.Context) {

	values, httpErr := controllers.MustGetAll(c)
	if httpErr != nil {
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	cookieName, err := values.Services.Configuration.GetString("COOKIE_NAME")
	if err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "invalid cookie name", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	//TODO: delete the user session

	c.SetCookie(cookieName, "", 0, "", "", false, false)

	sessionDb := values.Services.SessionsDB
	sessionId := values.User.SessionID
	sessionDb.Del(sessionId.Hex())

	c.Status(http.StatusOK)
}
