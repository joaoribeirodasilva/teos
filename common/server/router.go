package server

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/conf"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/database"
	"github.com/joaoribeirodasilva/teos/common/logger"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/utils/token"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Router struct {
	Conf          *conf.Conf
	Gin           *gin.Engine
	Db            *database.Db
	Configuration *configuration.Configuration
	sessionsDB    *redisdb.RedisDB
	permissionsDB *redisdb.RedisDB
}

func NewRouter(gin *gin.Engine, conf *conf.Conf, db *database.Db, configuration *configuration.Configuration, sessionsDB *redisdb.RedisDB, permissionsDB *redisdb.RedisDB) *Router {

	r := &Router{}

	r.Conf = conf
	r.Gin = gin
	r.Db = db
	r.Configuration = configuration
	r.sessionsDB = sessionsDB
	r.permissionsDB = permissionsDB
	return r
}

func (r *Router) Variables(c *gin.Context) {

	//fmt.Println("router variables")
	c.Set("db", r.Db)
	c.Set("conf", r.Conf)
	c.Set("configuration", r.Configuration)
	c.Set("sessionsDb", r.sessionsDB)
	c.Set("permissionsDb", r.permissionsDB)

	tokenUser := token.User{}
	tokenUser.ID, _ = primitive.ObjectIDFromHex("666758b475cf5396aea26a13")
	tokenUser.Name = "Sistema"
	tokenUser.Surname = "Teos"
	tokenUser.Email = "admin@teos.com.br"
	tokenUser.SessionID = primitive.NewObjectID()

	c.Set("user", &tokenUser)
}

func (r *Router) IsLogged(c *gin.Context) {

	cookie, err := c.Cookie("teos_auth")
	if err != nil {
		httpErr := logger.Error(logger.LogStatusForbidden, nil, "no credentials", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr.Error())
		return
	}

	vars, httpErr := controllers.MustGetAll(c)
	if httpErr != nil {
		c.AbortWithStatusJSON(httpErr.Status, httpErr.Error())
		return
	}

	token := token.New(vars.Configuration)
	if !token.IsValid(cookie) {
		httpErr := logger.Error(logger.LogStatusForbidden, nil, "invalid token", err, nil)
		c.AbortWithStatusJSON(httpErr.Status, httpErr.Error())
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
