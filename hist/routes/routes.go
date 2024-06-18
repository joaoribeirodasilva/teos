package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/hist/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Services.Gin.GET("/users/users", router.Variables, controllers.HistoriesList)
	router.Services.Gin.GET("/users/users/:id", router.Variables, controllers.HistoriesGet)
}
