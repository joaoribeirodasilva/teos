package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/hist/controllers"
)

func RegisterRoutes(router *server.Router) {
	router.Service.Http.Engine.GET("/hist/history", router.Services, controllers.HistoryList)
	router.Service.Http.Engine.GET("/hist/history/:id", router.Services, controllers.HistoryGet)
	//router.Services.Gin.POST("/hist/history", router.Variables, controllers.HistoryCreate)
}
