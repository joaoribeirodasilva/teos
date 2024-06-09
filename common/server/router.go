package server

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/database"
)

type Router struct {
	Conf *conf.Conf
	Gin  *gin.Engine
	Db   *database.Db
}

func NewRouter(gin *gin.Engine, conf *conf.Conf, db *database.Db) *Router {

	r := &Router{}

	r.Conf = conf
	r.Gin = gin
	r.Db = db

	return r
}

func (r *Router) Variables(c *gin.Context) {

	//fmt.Println("router variables")
	c.Set("db", r.Db)
	c.Set("conf", r.Conf)
}

func (r *Router) IsLogged(c *gin.Context) {

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
