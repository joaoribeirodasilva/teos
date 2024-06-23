package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/consumers/logs/controllers"
)

func RegisterRoutes(router *server.Router) {
	router.Service.Http.Engine.GET("/svc/logs", router.Services, controllers.LogsList, router.SendAuth)
	router.Service.Http.Engine.GET("/svc/logs/stats", router.Services, controllers.LogsStats, router.SendAuth)
	router.Service.Http.Engine.GET("/svc/logs/save", router.Services, controllers.LogsSave, router.SendAuth)
	router.Service.Http.Engine.GET("/svc/logs/:id", router.Services, controllers.LogsGet, router.SendAuth)
}
