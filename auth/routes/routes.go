package routes

import (
	"github.com/joaoribeirodasilva/teos/auth/controllers"
	"github.com/joaoribeirodasilva/teos/common/server"
)

func RegisterRoutes(router *server.Router) {
	router.Service.Http.Engine.POST("/auth/login", router.Services, controllers.AuthLogin, router.SendAuth)
	router.Service.Http.Engine.POST("/auth/forgot", router.Services, controllers.AuthForgot, router.SendAuth)
	router.Service.Http.Engine.PUT("/auth/reset/:key", router.Services, controllers.AuthReset, router.SendAuth)
	router.Service.Http.Engine.PATCH("/auth/reset/:key", router.Services, controllers.AuthReset, router.SendAuth)
	router.Service.Http.Engine.DELETE("/auth/logout", router.Services, router.IsLogged, controllers.AuthLogout, router.SendAuth)
}
