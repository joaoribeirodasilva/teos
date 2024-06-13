package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
)

type Router struct {
	Conf          *conf.Conf
	Gin           *gin.Engine
	Db            *database.Db
	Configuration *configuration.Configuration
}

func NewRouter(gin *gin.Engine, conf *conf.Conf, db *database.Db, configuration *configuration.Configuration) *Router {

	r := &Router{}

	r.Conf = conf
	r.Gin = gin
	r.Db = db
	r.Configuration = configuration

	return r
}

func (r *Router) Variables(c *gin.Context) {

	//fmt.Println("router variables")
	c.Set("db", r.Db)
	c.Set("conf", r.Conf)
	c.Set("configuration", r.Configuration)
}

func (r *Router) IsLogged(c *gin.Context) {

	cookie, err := c.Cookie("teos_auth")
	if err != nil {
		appErr := service_errors.New(0, http.StatusForbidden, "ROUTER", "IsLogged", "", "ERR: %s", err.Error()).LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	vars, appErr := controllers.MustGetAll(c)
	if appErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token := token.New(vars.Configuration)
	if !token.IsValid(cookie) {
		appErr := service_errors.New(0, http.StatusForbidden, "ROUTER", "IsLogged", "", "invalid token").LogError()
		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
		return
	}

	c.Set("user", token.User)

	//fmt.Printf("Cookie: %s\n", cookie)

	// fmt.Println("router is logged")
	/* 	auth := token.New(r.conf)
	   	if !auth.IsValid(c.GetHeader("Authorization")) {
	   		c.AbortWithStatus(http.StatusForbidden)
	   		c.Abort()
	   		return
	   	}

	   	c.Set("auth", auth.User) */

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
