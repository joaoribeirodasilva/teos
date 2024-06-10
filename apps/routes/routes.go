package routes

import (
	"github.com/joaoribeirodasilva/teos/apps/controllers"
	"github.com/joaoribeirodasilva/teos/common/server"
)

func RegisterRoutes(router *server.Router) {

	router.Gin.GET("/apps/apps", router.Variables, controllers.AppAppsList)
	router.Gin.GET("/apps/apps/:id", router.Variables, controllers.AppAppsGet)
	router.Gin.POST("/apps/apps", router.Variables, controllers.AppAppsCreate)
	router.Gin.PUT("/apps/apps/:id", router.Variables, controllers.AppAppsUpdate)
	router.Gin.PATCH("/apps/apps/:id", router.Variables, controllers.AppAppsUpdate)
	router.Gin.DELETE("/apps/apps/:id", router.Variables, controllers.AppAppsDelete)

	router.Gin.GET("/apps/configurations", router.Variables, controllers.AppConfigurationsList)
	router.Gin.GET("/apps/configurations:id", router.Variables, controllers.AppConfigurationsGet)
	router.Gin.POST("/apps/configurations", router.Variables, controllers.AppConfigurationsCreate)
	router.Gin.PUT("/apps/configurations/:id", router.Variables, controllers.AppConfigurationsUpdate)
	router.Gin.PATCH("/apps/configurations/:id", router.Variables, controllers.AppConfigurationsUpdate)
	router.Gin.DELETE("/apps/configurations/:id", router.Variables, controllers.AppConfigurationsDelete)

	router.Gin.GET("/apps/routes", router.Variables, controllers.AppRoutesList)
	router.Gin.GET("/apps/routes/:id", router.Variables, controllers.AppRoutesGet)
	router.Gin.POST("/apps/routes", router.Variables, controllers.AppRoutesCreate)
	router.Gin.PUT("/apps/routes/:id", router.Variables, controllers.AppRoutesUpdate)
	router.Gin.PATCH("/apps/routes/:id", router.Variables, controllers.AppRoutesUpdate)
	router.Gin.DELETE("/apps/routes/:id", router.Variables, controllers.AppRoutesDelete)

	router.Gin.GET("/apps/routesblocks", router.Variables, controllers.AppRoutesBlocksList)
	router.Gin.GET("/apps/routesblocks/:id", router.Variables, controllers.AppRoutesBlocksGet)
	router.Gin.POST("/apps/routesblocks", router.Variables, controllers.AppRoutesBlocksCreate)
	router.Gin.PUT("/apps/routesblocks/:id", router.Variables, controllers.AppRoutesBlocksUpdate)
	router.Gin.PATCH("/apps/routesblocks/:id", router.Variables, controllers.AppRoutesBlocksUpdate)
	router.Gin.DELETE("/apps/routesblocks/:id", router.Variables, controllers.AppRoutesBlocksDelete)

}
