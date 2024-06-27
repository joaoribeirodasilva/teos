package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/web/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Service.Http.Engine.Static("/assets", "./assets")
	router.Service.Http.Engine.LoadHTMLGlob("templates/**/*.templ")

	router.Service.Http.Engine.GET("/login", controllers.AuthLogin)
	router.Service.Http.Engine.GET("/forgot", controllers.AuthForgot)
	router.Service.Http.Engine.GET("/signup", controllers.AuthSignup)
	router.Service.Http.Engine.GET("/reset", controllers.AuthReset)

	router.Service.Http.Engine.GET("/", controllers.IndexIndex)
	router.Service.Http.Engine.GET("/admin", controllers.AdminIndex)

	router.Service.Http.Engine.GET("/admin/users/profile", controllers.UsersProfile)
}
