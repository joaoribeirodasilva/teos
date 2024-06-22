package routes

import (
	"github.com/joaoribeirodasilva/teos/apps/controllers"
	"github.com/joaoribeirodasilva/teos/common/server"
)

func RegisterRoutes(router *server.Router) {

	router.Service.Http.Engine.GET("/apps/apps", router.Services, router.IsLogged, controllers.ApplicationsList)
	router.Service.Http.Engine.GET("/apps/apps/:id", router.Services, router.IsLogged, controllers.ApplicationsGet)
	router.Service.Http.Engine.POST("/apps/apps", router.Services, controllers.ApplicationsCreate)
	router.Service.Http.Engine.PUT("/apps/apps/:id", router.Services, router.IsLogged, controllers.ApplicationsUpdate)
	router.Service.Http.Engine.PATCH("/apps/apps/:id", router.Services, router.IsLogged, controllers.ApplicationsUpdate)
	router.Service.Http.Engine.DELETE("/apps/apps/:id", router.Services, router.IsLogged, controllers.ApplicationsDelete)

	router.Service.Http.Engine.GET("/apps/configurations", router.Services, router.IsLogged, controllers.AppConfigurationsList)
	router.Service.Http.Engine.GET("/apps/configurations:id", router.Services, router.IsLogged, controllers.AppConfigurationsGet)
	router.Service.Http.Engine.POST("/apps/configurations", router.Services, router.IsLogged, controllers.AppConfigurationsCreate)
	router.Service.Http.Engine.PUT("/apps/configurations/:id", router.Services, router.IsLogged, controllers.AppConfigurationsUpdate)
	router.Service.Http.Engine.PATCH("/apps/configurations/:id", router.Services, router.IsLogged, controllers.AppConfigurationsUpdate)
	router.Service.Http.Engine.DELETE("/apps/configurations/:id", router.Services, router.IsLogged, controllers.AppConfigurationsDelete)

	router.Service.Http.Engine.GET("/apps/routes", router.Services, router.IsLogged, controllers.AppRoutesList)
	router.Service.Http.Engine.GET("/apps/routes/:id", router.Services, router.IsLogged, controllers.AppRoutesGet)
	router.Service.Http.Engine.POST("/apps/routes", router.Services, router.IsLogged, controllers.AppRoutesCreate)
	router.Service.Http.Engine.PUT("/apps/routes/:id", router.Services, router.IsLogged, controllers.AppRoutesUpdate)
	router.Service.Http.Engine.PATCH("/apps/routes/:id", router.Services, router.IsLogged, controllers.AppRoutesUpdate)
	router.Service.Http.Engine.DELETE("/apps/routes/:id", router.Services, router.IsLogged, controllers.AppRoutesDelete)

	router.Service.Http.Engine.GET("/apps/environments", router.Services, router.IsLogged, controllers.AppEnvironmentsList)
	router.Service.Http.Engine.GET("/apps/environments/:id", router.Services, router.IsLogged, controllers.AppEnvironmentsGet)
	router.Service.Http.Engine.POST("/apps/environments", router.Services, router.IsLogged, controllers.AppEnvironmentsCreate)
	router.Service.Http.Engine.PUT("/apps/environments/:id", router.Services, router.IsLogged, controllers.AppEnvironmentsUpdate)
	router.Service.Http.Engine.PATCH("/apps/environments/:id", router.Services, router.IsLogged, controllers.AppEnvironmentsUpdate)
	router.Service.Http.Engine.DELETE("/apps/environments/:id", router.Services, router.IsLogged, controllers.AppEnvironmentsDelete)

	router.Service.Http.Engine.GET("/apps/routemethods", router.Services, router.IsLogged, controllers.AppRouteMethodsList)
	router.Service.Http.Engine.GET("/apps/routemethods/:id", router.Services, router.IsLogged, controllers.AppRouteMethodsGet)
	router.Service.Http.Engine.POST("/apps/routemethods", router.Services, router.IsLogged, controllers.AppRouteMethodsCreate)
	router.Service.Http.Engine.PUT("/apps/routemethods/:id", router.Services, router.IsLogged, controllers.AppRouteMethodsUpdate)
	router.Service.Http.Engine.PATCH("/apps/routemethods/:id", router.Services, router.IsLogged, controllers.AppRouteMethodsUpdate)
	router.Service.Http.Engine.DELETE("/apps/routemethods/:id", router.Services, router.IsLogged, controllers.AppRouteMethodsDelete)

}
