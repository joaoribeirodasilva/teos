package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/users/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Gin.GET("/users/users", router.Variables, controllers.UserUsersList)
	router.Gin.GET("/users/users/:id", router.Variables, controllers.UserUsersGet)
	router.Gin.POST("/users/users", router.Variables, controllers.UserUsersCreate)
	router.Gin.PUT("/users/users/:id", router.Variables, controllers.UserUsersUpdate)
	router.Gin.PATCH("/users/users/:id", router.Variables, controllers.UserUsersUpdate)
	router.Gin.DELETE("/users/users/:id", router.Variables, controllers.UserUsersDelete)

	router.Gin.GET("/users/groups", router.Variables, controllers.UserGroupsList)
	router.Gin.GET("/users/groups/:id", router.Variables, controllers.UserGroupsGet)
	router.Gin.POST("/users/groups", router.Variables, controllers.UserGroupsCreate)
	router.Gin.PUT("/users/groups/:id", router.Variables, controllers.UserGroupsUpdate)
	router.Gin.PATCH("/users/groups/:id", router.Variables, controllers.UserGroupsUpdate)
	router.Gin.DELETE("/users/groups/:id", router.Variables, controllers.UserGroupsDelete)

	router.Gin.GET("/users/roles", router.Variables, controllers.UserRolesList)
	router.Gin.GET("/users/roles/:id", router.Variables, controllers.UserRolesGet)
	router.Gin.POST("/users/roles", router.Variables, controllers.UserRolesCreate)
	router.Gin.PUT("/users/roles/:id", router.Variables, controllers.UserRolesUpdate)
	router.Gin.PATCH("/users/roles/:id", router.Variables, controllers.UserRolesUpdate)
	router.Gin.DELETE("/users/roles/:id", router.Variables, controllers.UserRolesDelete)

	router.Gin.GET("/users/rolegroups", router.Variables, controllers.UserRolesGroupsList)
	router.Gin.GET("/users/rolegroups/:id", router.Variables, controllers.UserRolesGroupsGet)
	router.Gin.POST("/users/rolegroups", router.Variables, controllers.UserRolesGroupsCreate)
	router.Gin.PUT("/users/rolegroups/:id", router.Variables, controllers.UserRolesGroupsUpdate)
	router.Gin.PATCH("/users/rolegroups/:id", router.Variables, controllers.UserRolesGroupsUpdate)
	router.Gin.DELETE("/users/rolegroups/:id", router.Variables, controllers.UserRolesGroupsDelete)

	router.Gin.GET("/users/resettypes", router.Variables, controllers.UserResetTypesList)
	router.Gin.GET("/users/resettypes/:id", router.Variables, controllers.UserResetTypesGet)
	router.Gin.POST("/users/resettypes", router.Variables, controllers.UserResetTypesCreate)
	router.Gin.PUT("/users/resettypes/:id", router.Variables, controllers.UserResetTypesUpdate)
	router.Gin.PATCH("/users/resettypes/:id", router.Variables, controllers.UserResetTypesUpdate)
	router.Gin.DELETE("/users/resettypes/:id", router.Variables, controllers.UserResetTypesDelete)

	router.Gin.GET("/users/resets", router.Variables, controllers.UserResetsList)
	router.Gin.GET("/users/resets/:id", router.Variables, controllers.UserResetsGet)
	router.Gin.PUT("/users/resets/:id", router.Variables, controllers.UserResetsUpdate)
	router.Gin.PATCH("/resets/:id", router.Variables, controllers.UserResetsUpdate)

	router.Gin.GET("/users/sessions", router.Variables, controllers.UserSessionsList)
	router.Gin.GET("/users/sessions/:id", router.Variables, controllers.UserSessionsGet)
	router.Gin.PUT("/users/sessions/:id", router.Variables, controllers.UserSessionsUpdate)
	router.Gin.PATCH("/users/sessions/:id", router.Variables, controllers.UserSessionsUpdate)
	router.Gin.DELETE("/users/sessions/:id", router.Variables, controllers.UserSessionsDelete)

}
