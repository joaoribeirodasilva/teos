package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/users/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Gin.GET("/users", router.Variables, controllers.UserUsersList)
	router.Gin.GET("/users/:id", router.Variables, controllers.UserUsersGet)
	router.Gin.POST("/users", router.Variables, controllers.UserUsersCreate)
	router.Gin.PUT("/users/:id", router.Variables, controllers.UserUsersUpdate)
	router.Gin.PATCH("/users/:id", router.Variables, controllers.UserUsersUpdate)
	router.Gin.DELETE("/users/:id", router.Variables, controllers.UserUsersDelete)

	router.Gin.GET("/groups", router.Variables, controllers.UserGroupsList)
	router.Gin.GET("/groups/:id", router.Variables, controllers.UserGroupsGet)
	router.Gin.POST("/groups", router.Variables, controllers.UserGroupsCreate)
	router.Gin.PUT("/groups/:id", router.Variables, controllers.UserGroupsUpdate)
	router.Gin.PATCH("/groups/:id", router.Variables, controllers.UserGroupsUpdate)
	router.Gin.DELETE("/groups/:id", router.Variables, controllers.UserGroupsDelete)

	router.Gin.GET("/roles", router.Variables, controllers.UserRolesList)
	router.Gin.GET("/roles/:id", router.Variables, controllers.UserRolesGet)
	router.Gin.POST("/roles", router.Variables, controllers.UserRolesCreate)
	router.Gin.PUT("/roles/:id", router.Variables, controllers.UserRolesUpdate)
	router.Gin.PATCH("/roles/:id", router.Variables, controllers.UserRolesUpdate)
	router.Gin.DELETE("/roles/:id", router.Variables, controllers.UserRolesDelete)

	router.Gin.GET("/rolegroups", router.Variables, controllers.UserRolesGroupsList)
	router.Gin.GET("/rolegroups/:id", router.Variables, controllers.UserRolesGroupsGet)
	router.Gin.POST("/rolegroups", router.Variables, controllers.UserRolesGroupsCreate)
	router.Gin.PUT("/rolegroups/:id", router.Variables, controllers.UserRolesGroupsUpdate)
	router.Gin.PATCH("/rolegroups/:id", router.Variables, controllers.UserRolesGroupsUpdate)
	router.Gin.DELETE("/rolegroups/:id", router.Variables, controllers.UserRolesGroupsDelete)
}
