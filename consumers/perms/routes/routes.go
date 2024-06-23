package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/consumers/perms/controllers"
)

func RegisterRoutes(router *server.Router) {
	router.Service.Http.Engine.GET("/svc/permissions", router.Services, controllers.PermissionsList, router.SendAuth)
	router.Service.Http.Engine.GET("/svc/permissions/:id", router.Services, controllers.PermissionsGet, router.SendAuth)
	router.Service.Http.Engine.GET("/svc/permissions/reload", router.Services, controllers.PermissionsReload, router.SendAuth)

	router.Service.Http.Engine.GET("/svc/routes", router.Services, controllers.RoutesList, router.SendAuth)
	router.Service.Http.Engine.GET("/svc/routes/:id", router.Services, controllers.RoutesGet, router.SendAuth)
	router.Service.Http.Engine.GET("/svc/routes/reload", router.Services, controllers.RoutesReload, router.SendAuth)
}
