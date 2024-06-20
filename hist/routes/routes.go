package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/hist/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Services.Gin.GET("/hist/history", router.Variables, controllers.HistoryList)
	router.Services.Gin.GET("/hist/history/:id", router.Variables, controllers.HistoryGet)
	//router.Services.Gin.POST("/hist/history", router.Variables, controllers.HistoryCreate)
}
