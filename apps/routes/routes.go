package routes

import (
	"github.com/joaoribeirodasilva/teos/apps/controllers"
	"github.com/joaoribeirodasilva/teos/common/server"
)

func RegisterRoutes(router *server.Router) {

	router.Services.Gin.GET("/apps/apps", router.Variables, router.IsLogged, controllers.AppAppsList)
	//router.Gin.GET("/apps/apps", router.Variables, controllers.AppAppsList)
	router.Services.Gin.GET("/apps/apps/:id", router.Variables, router.IsLogged, controllers.AppAppsGet)
	//router.Gin.GET("/apps/apps/:id", router.Variables, controllers.AppAppsGet)
	//router.Gin.POST("/apps/apps", router.Variables, router.IsLogged, controllers.AppAppsCreate)
	router.Services.Gin.POST("/apps/apps", router.Variables, controllers.AppAppsCreate)

	router.Services.Gin.PUT("/apps/apps/:id", router.Variables, router.IsLogged, controllers.AppAppsUpdate)
	//router.Gin.PUT("/apps/apps/:id", router.Variables, controllers.AppAppsUpdate)

	router.Services.Gin.PATCH("/apps/apps/:id", router.Variables, router.IsLogged, controllers.AppAppsUpdate)

	router.Services.Gin.DELETE("/apps/apps/:id", router.Variables, router.IsLogged, controllers.AppAppsDelete)
	//router.Gin.DELETE("/apps/apps/:id", router.Variables, controllers.AppAppsDelete)

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

	router.Services.Gin.GET("/apps/routesblocks", router.Variables, router.IsLogged, controllers.AppRoutesBlocksList)
	router.Services.Gin.GET("/apps/routesblocks/:id", router.Variables, router.IsLogged, controllers.AppRoutesBlocksGet)
	router.Services.Gin.POST("/apps/routesblocks", router.Variables, router.IsLogged, controllers.AppRoutesBlocksCreate)
	router.Services.Gin.PUT("/apps/routesblocks/:id", router.Variables, router.IsLogged, controllers.AppRoutesBlocksUpdate)
	router.Services.Gin.PATCH("/apps/routesblocks/:id", router.Variables, router.IsLogged, controllers.AppRoutesBlocksUpdate)
	router.Services.Gin.DELETE("/apps/routesblocks/:id", router.Variables, router.IsLogged, controllers.AppRoutesBlocksDelete)

}
