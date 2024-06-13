package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/users/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Gin.GET("/users/users", router.Variables, router.IsLogged, controllers.UserUsersList)
	router.Gin.GET("/users/users/:id", router.Variables, router.IsLogged, controllers.UserUsersGet)
	router.Gin.POST("/users/users", router.Variables, controllers.UserUsersCreate)
	router.Gin.PUT("/users/users/:id", router.Variables, router.IsLogged, controllers.UserUsersUpdate)
	router.Gin.PATCH("/users/users/:id", router.Variables, router.IsLogged, controllers.UserUsersUpdate)
	router.Gin.DELETE("/users/users/:id", router.Variables, router.IsLogged, controllers.UserUsersDelete)

	router.Gin.GET("/users/groups", router.Variables, router.IsLogged, controllers.UserGroupsList)
	router.Gin.GET("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsGet)
	router.Gin.POST("/users/groups", router.Variables, router.IsLogged, controllers.UserGroupsCreate)
	router.Gin.PUT("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsUpdate)
	router.Gin.PATCH("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsUpdate)
	router.Gin.DELETE("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsDelete)

	router.Gin.GET("/users/roles", router.Variables, router.IsLogged, controllers.UserRolesList)
	router.Gin.GET("/users/roles/:id", router.Variables, router.IsLogged, controllers.UserRolesGet)
	router.Gin.POST("/users/roles", router.Variables, router.IsLogged, controllers.UserRolesCreate)
	router.Gin.PUT("/users/roles/:id", router.Variables, router.IsLogged, controllers.UserRolesUpdate)
	router.Gin.PATCH("/users/roles/:id", router.Variables, router.IsLogged, controllers.UserRolesUpdate)
	router.Gin.DELETE("/users/roles/:id", router.Variables, router.IsLogged, controllers.UserRolesDelete)

	router.Gin.GET("/users/rolegroups", router.Variables, router.IsLogged, controllers.UserRolesGroupsList)
	router.Gin.GET("/users/rolegroups/:id", router.Variables, router.IsLogged, controllers.UserRolesGroupsGet)
	router.Gin.POST("/users/rolegroups", router.Variables, router.IsLogged, controllers.UserRolesGroupsCreate)
	router.Gin.PUT("/users/rolegroups/:id", router.Variables, router.IsLogged, controllers.UserRolesGroupsUpdate)
	router.Gin.PATCH("/users/rolegroups/:id", router.Variables, router.IsLogged, controllers.UserRolesGroupsUpdate)
	router.Gin.DELETE("/users/rolegroups/:id", router.Variables, router.IsLogged, controllers.UserRolesGroupsDelete)

	router.Gin.GET("/users/resettypes", router.Variables, router.IsLogged, controllers.UserResetTypesList)
	router.Gin.GET("/users/resettypes/:id", router.Variables, router.IsLogged, controllers.UserResetTypesGet)
	router.Gin.POST("/users/resettypes", router.Variables, router.IsLogged, controllers.UserResetTypesCreate)
	router.Gin.PUT("/users/resettypes/:id", router.Variables, router.IsLogged, controllers.UserResetTypesUpdate)
	router.Gin.PATCH("/users/resettypes/:id", router.Variables, router.IsLogged, controllers.UserResetTypesUpdate)
	router.Gin.DELETE("/users/resettypes/:id", router.Variables, router.IsLogged, controllers.UserResetTypesDelete)

	router.Gin.GET("/users/resets", router.Variables, router.IsLogged, controllers.UserResetsList)
	router.Gin.GET("/users/resets/:id", router.Variables, router.IsLogged, controllers.UserResetsGet)
	router.Gin.PUT("/users/resets/:id", router.Variables, router.IsLogged, controllers.UserResetsUpdate)
	router.Gin.PATCH("/resets/:id", router.Variables, router.IsLogged, controllers.UserResetsUpdate)

	router.Gin.GET("/users/sessions", router.Variables, router.IsLogged, controllers.UserSessionsList)
	router.Gin.GET("/users/sessions/:id", router.Variables, router.IsLogged, controllers.UserSessionsGet)
	router.Gin.PUT("/users/sessions/:id", router.Variables, router.IsLogged, controllers.UserSessionsUpdate)
	router.Gin.PATCH("/users/sessions/:id", router.Variables, router.IsLogged, controllers.UserSessionsUpdate)
	router.Gin.DELETE("/users/sessions/:id", router.Variables, router.IsLogged, controllers.UserSessionsDelete)

}
