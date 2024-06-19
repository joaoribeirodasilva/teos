package server

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/structures"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
)

type Router struct {
	Services *structures.Services
}

func NewRouter(services *structures.Services) *Router {

	r := &Router{}

	r.Services = services

	return r
}

func (r *Router) Variables(c *gin.Context) {

	//fmt.Println("router variables")
	c.Set("services", r.Services)

	tokenUser := token.User{}
	tokenUser.ID = 1
	tokenUser.Name = "Sistema"
	tokenUser.Surname = "Teos"
	tokenUser.Email = "teos@teos.com.br"
	tokenUser.SessionID = 0
	tokenUser.OrganizationID = 1

	c.Set("user", &tokenUser)
}

func (r *Router) IsLogged(c *gin.Context) {

	cookie, err := c.Cookie("teos_auth")
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

	c.Set("user", token.User)

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
