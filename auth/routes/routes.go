package routes

import (
	"github.com/joaoribeirodasilva/teos/auth/controllers"
	"github.com/joaoribeirodasilva/teos/common/server"
)

func RegisterRoutes(router *server.Router) {
	router.Services.Gin.POST("/auth/login", router.Variables, controllers.AuthLogin)
	router.Services.Gin.POST("/auth/forgot", router.Variables, controllers.AuthForgot)
	router.Services.Gin.PUT("/auth/reset/:key", router.Variables, controllers.AuthReset)
	router.Services.Gin.PATCH("/auth/reset/:key", router.Variables, controllers.AuthReset)
	router.Services.Gin.DELETE("/auth/logout", router.Variables, router.IsLogged, controllers.AuthLogout)
}
