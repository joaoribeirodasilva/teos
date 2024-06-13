package routes

import (
	"github.com/joaoribeirodasilva/teos/apps/controllers"
	"github.com/joaoribeirodasilva/teos/common/server"
)

func RegisterRoutes(router *server.Router) {

	router.Gin.GET("/apps/apps", router.Variables, router.IsLogged, controllers.AppAppsList)
	router.Gin.GET("/apps/apps/:id", router.Variables, router.IsLogged, controllers.AppAppsGet)
	router.Gin.POST("/apps/apps", router.Variables, router.IsLogged, controllers.AppAppsCreate)
	router.Gin.PUT("/apps/apps/:id", router.Variables, router.IsLogged, controllers.AppAppsUpdate)
	router.Gin.PATCH("/apps/apps/:id", router.Variables, router.IsLogged, controllers.AppAppsUpdate)
	router.Gin.DELETE("/apps/apps/:id", router.Variables, router.IsLogged, controllers.AppAppsDelete)

	router.Gin.GET("/apps/configurations", router.Variables, router.IsLogged, controllers.AppConfigurationsList)
	router.Gin.GET("/apps/configurations:id", router.Variables, router.IsLogged, controllers.AppConfigurationsGet)
	router.Gin.POST("/apps/configurations", router.Variables, router.IsLogged, controllers.AppConfigurationsCreate)
	router.Gin.PUT("/apps/configurations/:id", router.Variables, router.IsLogged, controllers.AppConfigurationsUpdate)
	router.Gin.PATCH("/apps/configurations/:id", router.Variables, router.IsLogged, controllers.AppConfigurationsUpdate)
	router.Gin.DELETE("/apps/configurations/:id", router.Variables, router.IsLogged, controllers.AppConfigurationsDelete)

	router.Gin.GET("/apps/routes", router.Variables, router.IsLogged, controllers.AppRoutesList)
	router.Gin.GET("/apps/routes/:id", router.Variables, router.IsLogged, controllers.AppRoutesGet)
	router.Gin.POST("/apps/routes", router.Variables, router.IsLogged, controllers.AppRoutesCreate)
	router.Gin.PUT("/apps/routes/:id", router.Variables, router.IsLogged, controllers.AppRoutesUpdate)
	router.Gin.PATCH("/apps/routes/:id", router.Variables, router.IsLogged, controllers.AppRoutesUpdate)
	router.Gin.DELETE("/apps/routes/:id", router.Variables, router.IsLogged, controllers.AppRoutesDelete)

	router.Gin.GET("/apps/routesblocks", router.Variables, router.IsLogged, controllers.AppRoutesBlocksList)
	router.Gin.GET("/apps/routesblocks/:id", router.Variables, router.IsLogged, controllers.AppRoutesBlocksGet)
	router.Gin.POST("/apps/routesblocks", router.Variables, router.IsLogged, controllers.AppRoutesBlocksCreate)
	router.Gin.PUT("/apps/routesblocks/:id", router.Variables, router.IsLogged, controllers.AppRoutesBlocksUpdate)
	router.Gin.PATCH("/apps/routesblocks/:id", router.Variables, router.IsLogged, controllers.AppRoutesBlocksUpdate)
	router.Gin.DELETE("/apps/routesblocks/:id", router.Variables, router.IsLogged, controllers.AppRoutesBlocksDelete)

}
