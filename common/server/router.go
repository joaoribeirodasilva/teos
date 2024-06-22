package server

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/payload"
)

type Router struct {
	Service *payload.Payload
}

func NewRouter(request *payload.Payload) *Router {

	r := &Router{}

	r.Service = request

	return r
}

func (r *Router) Services(c *gin.Context) {

	session := &payload.SessionAuth{}
	session.ID = 0
	session.OrganizationID = 1
	session.UserID = 1
	session.Name = "Sistema"
	session.Surname = "Teos"
	session.Email = "teos@teos.com.br"
	session.AvatarUrl = ""

	r.Service.Http.Request.Session.Auth.UserSession = session

	c.Set("service", r.Service)
	c.Next()
}

func (r *Router) SendAuth(c *gin.Context) {

	serviceInterface := c.MustGet("service")
	services, ok := serviceInterface.(*payload.Payload)
	if !ok {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "invalid service pointer", nil, nil)
		c.AbortWithStatusJSON(int(httpErr.Status), httpErr)
	}

	if err := services.Http.Request.Session.Auth.SetCookie(); err != nil {
		httpErr := logger.Error(logger.LogStatusInternalServerError, nil, "failed to generate session cookie", err, nil)
		c.AbortWithStatusJSON(int(httpErr.Status), httpErr.Error())
		c.Abort()
		return
	}

	c.Next()
}

func (r *Router) IsLogged(c *gin.Context) {

	// set session
	// check route open
	// authorization

	/* 	cookie, err := c.Cookie("teos_auth")
	   	if err != nil {
	   		httpErr := logger.Error(logger.LogStatusForbidden, nil, "no credentials", err, nil)
	   		c.AbortWithStatusJSON(int(httpErr.Status), httpErr.Error())
	   		return
	   	}

	   	values, httpErr := controllers.MustGetAll(c)
	   	if httpErr != nil {
	   		c.AbortWithStatusJSON(int(httpErr.Status), httpErr.Error())
	   		return
	   	}

	   	token := token.New(values.Services.Configuration)
	   	if !token.IsValid(cookie) {
	   		httpErr := logger.Error(logger.LogStatusForbidden, nil, "invalid token", err, nil)
	   		c.AbortWithStatusJSON(int(httpErr.Status), httpErr.Error())
	   		return
	   	}

	   	c.Set("user", token.User) */

	c.Next()
}

func (r *Router) IsAdmin(c *gin.Context) {

	// fmt.Println("router is admin")
	/* 	a := c.MustGet("auth")
	   	user, ok := a.(*token.User)
	   	if !ok || !user.Admin {
	   		c.AbortWithStatus(http.StatusUnauthorized)
	   		c.Abort()
	   		return
	   	} */

	c.Next()
}
