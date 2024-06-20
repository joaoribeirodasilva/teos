package routes

import (
	"github.com/joaoribeirodasilva/teos/apps/controllers"
	"github.com/joaoribeirodasilva/teos/common/server"
)

func RegisterRoutes(router *server.Router) {

	router.Services.Gin.GET("/apps/apps", router.Variables, router.IsLogged, controllers.ApplicationsList)
	router.Services.Gin.GET("/apps/apps/:id", router.Variables, router.IsLogged, controllers.ApplicationsGet)
	router.Services.Gin.POST("/apps/apps", router.Variables, controllers.ApplicationsCreate)
	router.Services.Gin.PUT("/apps/apps/:id", router.Variables, router.IsLogged, controllers.ApplicationsUpdate)
	router.Services.Gin.PATCH("/apps/apps/:id", router.Variables, router.IsLogged, controllers.ApplicationsUpdate)
	router.Services.Gin.DELETE("/apps/apps/:id", router.Variables, router.IsLogged, controllers.ApplicationsDelete)

	router.Services.Gin.GET("/apps/configurations", router.Variables, router.IsLogged, controllers.AppConfigurationsList)
	router.Services.Gin.GET("/apps/configurations:id", router.Variables, router.IsLogged, controllers.AppConfigurationsGet)
	router.Services.Gin.POST("/apps/configurations", router.Variables, router.IsLogged, controllers.AppConfigurationsCreate)
	router.Services.Gin.PUT("/apps/configurations/:id", router.Variables, router.IsLogged, controllers.AppConfigurationsUpdate)
	router.Services.Gin.PATCH("/apps/configurations/:id", router.Variables, router.IsLogged, controllers.AppConfigurationsUpdate)
	router.Services.Gin.DELETE("/apps/configurations/:id", router.Variables, router.IsLogged, controllers.AppConfigurationsDelete)

	router.Services.Gin.GET("/apps/routes", router.Variables, router.IsLogged, controllers.AppRoutesList)
	router.Services.Gin.GET("/apps/routes/:id", router.Variables, router.IsLogged, controllers.AppRoutesGet)
	router.Services.Gin.POST("/apps/routes", router.Variables, router.IsLogged, controllers.AppRoutesCreate)
	router.Services.Gin.PUT("/apps/routes/:id", router.Variables, router.IsLogged, controllers.AppRoutesUpdate)
	router.Services.Gin.PATCH("/apps/routes/:id", router.Variables, router.IsLogged, controllers.AppRoutesUpdate)
	router.Services.Gin.DELETE("/apps/routes/:id", router.Variables, router.IsLogged, controllers.AppRoutesDelete)

	router.Services.Gin.GET("/apps/environments", router.Variables, router.IsLogged, controllers.AppEnvironmentsList)
	router.Services.Gin.GET("/apps/environments/:id", router.Variables, router.IsLogged, controllers.AppEnvironmentsGet)
	router.Services.Gin.POST("/apps/environments", router.Variables, router.IsLogged, controllers.AppEnvironmentsCreate)
	router.Services.Gin.PUT("/apps/environments/:id", router.Variables, router.IsLogged, controllers.AppEnvironmentsUpdate)
	router.Services.Gin.PATCH("/apps/environments/:id", router.Variables, router.IsLogged, controllers.AppEnvironmentsUpdate)
	router.Services.Gin.DELETE("/apps/environments/:id", router.Variables, router.IsLogged, controllers.AppEnvironmentsDelete)

	router.Services.Gin.GET("/apps/routemethods", router.Variables, router.IsLogged, controllers.AppRouteMethodsList)
	router.Services.Gin.GET("/apps/routemethods/:id", router.Variables, router.IsLogged, controllers.AppRouteMethodsGet)
	router.Services.Gin.POST("/apps/routemethods", router.Variables, router.IsLogged, controllers.AppRouteMethodsCreate)
	router.Services.Gin.PUT("/apps/routemethods/:id", router.Variables, router.IsLogged, controllers.AppRouteMethodsUpdate)
	router.Services.Gin.PATCH("/apps/routemethods/:id", router.Variables, router.IsLogged, controllers.AppRouteMethodsUpdate)
	router.Services.Gin.DELETE("/apps/routemethods/:id", router.Variables, router.IsLogged, controllers.AppRouteMethodsDelete)

}
