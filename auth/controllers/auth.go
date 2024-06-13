package controllers

import (
	"context"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joaoribeirodasilva/teos/auth/requests"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"github.com/joaoribeirodasilva/teos/common/utils/password"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
	"github.com/joaoribeirodasilva/teos/users/models"
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

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	email, passwd, ok := c.Request.BasicAuth()
	if !ok {
		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::AuthLogin", "", "invalid username or password")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	record := models.UserUser{}

	coll := vars.Db.Db.Collection("user_users")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthLogin", "", "failed to query database. ERR: %s", err.Error())
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if record.Password == nil {
		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::AuthLogin", "", "invalid password")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	if !password.Check(passwd, *record.Password) {
		appErr := service_log.Error(0, http.StatusForbidden, "CONTROLLER::AuthLogin", "", "invalid username or password")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	if record.Active == 0 || record.DeletedBy != nil || record.DeletedAt != nil {
		appErr := service_log.Error(0, http.StatusUnauthorized, "CONTROLLER::AuthLogin", "", "account disabled")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	sessionRecord := models.UserSession{}

	now := time.Now().UTC()
	sessionRecord.ID = primitive.NewObjectID()
	sessionRecord.UserUserID = record.ID
	sessionRecord.CreatedBy = record.ID
	sessionRecord.CreatedAt = now
	sessionRecord.UpdatedBy = record.ID
	sessionRecord.UpdatedAt = now

	collSession := vars.Db.Db.Collection("user_sessions")

	result, err := collSession.InsertOne(context.TODO(), sessionRecord)
	if err != nil {
		appErr := service_log.Error(0, http.StatusConflict, "CONTROLLER::AuthLogin", "", "failed to insert session into the database. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	sessionId := result.InsertedID.(primitive.ObjectID)
	tokenObject := token.New(vars.Configuration)
	if appErr := tokenObject.Create(&record, &sessionId); err != nil {
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempName := vars.Configuration.GetKey("COOKIE_NAME")
	if tempName == nil || tempName.String == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthLogin", "", "invalid cookie name")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempExpire := vars.Configuration.GetKey("COOKIE_EXPIRE")
	if tempExpire == nil || tempExpire.Int == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthLogin", "", "invalid cookie expire")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempDomain := vars.Configuration.GetKey("COOKIE_DOMAIN")
	if tempDomain == nil || tempDomain.String == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthLogin", "", "invalid cookie domain")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempHttpOnly := vars.Configuration.GetKey("COOKIE_HTTP_ONLY")
	if tempHttpOnly == nil || tempHttpOnly.Bool == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthLogin", "", "invalid cookie http only")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempSecure := vars.Configuration.GetKey("COOKIE_SECURE")
	if tempSecure == nil || tempSecure.Bool == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthLogin", "", "invalid cookie secure")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	expire := int(time.Now().Add(time.Second * time.Duration(*tempExpire.Int)).Unix())
	// c.SetCookie(*tempName.String, tokenObject.TokenString, *tempExpire.Int, "/", *tempDomain.String, *tempSecure.Bool, *tempHttpOnly.Bool)
	c.SetCookie(*tempName.String, tokenObject.TokenString, expire, "", "", false, false)

	c.Status(http.StatusOK)
}

func AuthForgot(c *gin.Context) {

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	request := requests.ForgotPassword{}
	if err := c.ShouldBind(&request); err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthForgot", "", "no email provided")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))
	_, err := mail.ParseAddress(request.Email)
	if err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthForgot", "", "invalid email provided")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	record := models.UserUser{}
	coll := vars.Db.Db.Collection("user_users")
	if err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: request.Email}}).Decode(&record); err != nil {
		if err != mongo.ErrNoDocuments {
			appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthForgot", "", "failed to query database. ERR: %s", err.Error())
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthForgot", "", "wrong username or password")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	if record.DeletedBy != nil || record.DeletedAt != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthForgot", "", "wrong username or password")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	uuid.EnableRandPool()
	resetKey := uuid.NewString()

	expire := time.Now()
	expire = expire.Add(time.Hour * 24)
	resetTypeId, err := primitive.ObjectIDFromHex("6669fdf9175f523b82a26a13")
	if err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthForgot", "", "invalid reset type id")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	resetRecord := models.UserReset{
		UserResetTypeID: resetTypeId,
		UserUserID:      record.ID,
		ResetKey:        resetKey,
		Used:            nil,
		Expire:          expire,
	}

	collUserResets := vars.Db.Db.Collection("user_resets")
	if _, err := collUserResets.InsertOne(context.TODO(), resetRecord); err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthForgot", "", "failed to insert password reset into database. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	//TODO: send email with the request link

	c.Status(http.StatusCreated)

}

func AuthReset(c *gin.Context) {

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	key := c.Param("key")
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
		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::AuthReset", "", "the password must have at least 6 characters")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	request.CheckPassword = strings.TrimSpace(request.CheckPassword)
	if request.CheckPassword != request.Password {
		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::AuthReset", "", "the passwords don't match")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	record := models.UserReset{}

	coll := vars.Db.Db.Collection("user_resets")
	if err := coll.FindOne(context.TODO(), bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "resetKey", Value: key}},
			bson.D{{Key: "used", Value: nil}},
		}},
	}).Decode(&record); err != nil {
		if err == mongo.ErrNoDocuments {
			appErr := service_log.Error(0, http.StatusNotFound, "CONTROLLER::AuthReset", "", "the reset was not found or it expired")
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthReset", "", "failed to qurey database. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	now := time.Now().UTC()
	if now.After(record.Expire) {
		appErr := service_log.Error(0, http.StatusUnauthorized, "CONTROLLER::AuthReset", "", "the reset request has expired")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	userRecord := models.UserUser{}

	usersColl := vars.Db.Db.Collection("user_users")
	if err := usersColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: record.UserUserID}}).Decode(&userRecord); err != nil {
		if err == mongo.ErrNoDocuments {
			appErr := service_log.Error(0, http.StatusNotFound, "CONTROLLER::AuthReset", "", "user not found")
			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
			return
		}
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthReset", "", "failed to qurey database. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	if userRecord.Active == 0 {
		appErr := service_log.Error(0, http.StatusUnauthorized, "CONTROLLER::AuthReset", "", "account disabled")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempPassword, appErr := password.Hash(request.Password)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	now = time.Now().UTC()
	userRecord.Password = &tempPassword
	userRecord.UpdatedBy = userRecord.ID
	userRecord.UpdatedAt = now

	_, err := usersColl.UpdateOne(context.TODO(), bson.D{
		{Key: "_id", Value: userRecord.ID},
	}, bson.D{
		{Key: "$set", Value: userRecord},
	})
	if err != nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthReset", "", "failed to update password into database. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
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
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthReset", "", "failed to update password reset into database. ERR: %s", err.Error())
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	}

	c.Status(http.StatusOK)

}

func AuthLogout(c *gin.Context) {

	vars, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tempName := vars.Configuration.GetKey("COOKIE_NAME")
	if tempName == nil || tempName.String == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthReset", "", "invalid cookie name")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempExpire := vars.Configuration.GetKey("COOKIE_EXPIRE")
	if tempExpire == nil || tempExpire.Int == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthReset", "", "invalid cookie expire")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempDomain := vars.Configuration.GetKey("COOKIE_DOMAIN")
	if tempDomain == nil || tempDomain.String == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthReset", "", "invalid cookie domain")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempHttpOnly := vars.Configuration.GetKey("COOKIE_HTTP_ONLY")
	if tempHttpOnly == nil || tempHttpOnly.Bool == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthReset", "", "invalid cookie http only")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	tempSecure := vars.Configuration.GetKey("COOKIE_SECURE")
	if tempSecure == nil || tempSecure.Bool == nil {
		appErr := service_log.Error(0, http.StatusInternalServerError, "CONTROLLER::AuthReset", "", "invalid cookie secure")
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	//TODO: delete the user

	c.SetCookie(*tempName.String, "", 0, "", "", false, false)

	c.Status(http.StatusOK)
}
