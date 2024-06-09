package routes

import (
	"github.com/joaoribeirodasilva/teos/auth/controllers"
	"github.com/joaoribeirodasilva/teos/common/server"
)

func RegisterRoutes(router *server.Router) {
	router.Gin.POST("/auth/login", router.Variables, controllers.AuthLogin)
	router.Gin.POST("/auth/forgot", router.Variables, controllers.AuthForgot)
	router.Gin.PUT("/auth/reset/:key", router.Variables, controllers.AuthReset)
	router.Gin.PATCH("/auth/reset/:key", router.Variables, controllers.AuthReset)
	router.Gin.DELETE("/auth/logout", router.Variables, controllers.AuthLogout)
}
