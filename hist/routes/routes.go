package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/hist/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Gin.GET("/users/users", router.Variables, controllers.HistoriesList)
	router.Gin.GET("/users/users/:id", router.Variables, controllers.HistoriesGet)
	router.Gin.POST("/users/users", router.Variables, controllers.HistoriesCreate)
}
