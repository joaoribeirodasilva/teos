package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/users/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Services.Gin.GET("/users/users", router.Variables, router.IsLogged, controllers.UserUsersList)
	router.Services.Gin.GET("/users/users/:id", router.Variables, router.IsLogged, controllers.UserUsersGet)
	router.Services.Gin.POST("/users/users", router.Variables, controllers.UserUsersCreate)
	router.Services.Gin.PUT("/users/users/:id", router.Variables, router.IsLogged, controllers.UserUsersUpdate)
	router.Services.Gin.PATCH("/users/users/:id", router.Variables, router.IsLogged, controllers.UserUsersUpdate)
	router.Services.Gin.DELETE("/users/users/:id", router.Variables, router.IsLogged, controllers.UserUsersDelete)

	router.Services.Gin.GET("/users/groups", router.Variables, router.IsLogged, controllers.UserGroupsList)
	router.Services.Gin.GET("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsGet)
	router.Services.Gin.POST("/users/groups", router.Variables, router.IsLogged, controllers.UserGroupsCreate)
	router.Services.Gin.PUT("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsUpdate)
	router.Services.Gin.PATCH("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsUpdate)
	router.Services.Gin.DELETE("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsDelete)

	router.Services.Gin.GET("/users/roles", router.Variables, router.IsLogged, controllers.UserRolesList)
	router.Services.Gin.GET("/users/roles/:id", router.Variables, router.IsLogged, controllers.UserRolesGet)
	router.Services.Gin.POST("/users/roles", router.Variables, router.IsLogged, controllers.UserRolesCreate)
	router.Services.Gin.PUT("/users/roles/:id", router.Variables, router.IsLogged, controllers.UserRolesUpdate)
	router.Services.Gin.PATCH("/users/roles/:id", router.Variables, router.IsLogged, controllers.UserRolesUpdate)
	router.Services.Gin.DELETE("/users/roles/:id", router.Variables, router.IsLogged, controllers.UserRolesDelete)

	router.Services.Gin.GET("/users/rolegroups", router.Variables, router.IsLogged, controllers.UserRolesGroupsList)
	router.Services.Gin.GET("/users/rolegroups/:id", router.Variables, router.IsLogged, controllers.UserRolesGroupsGet)
	router.Services.Gin.POST("/users/rolegroups", router.Variables, router.IsLogged, controllers.UserRolesGroupsCreate)
	router.Services.Gin.PUT("/users/rolegroups/:id", router.Variables, router.IsLogged, controllers.UserRolesGroupsUpdate)
	router.Services.Gin.PATCH("/users/rolegroups/:id", router.Variables, router.IsLogged, controllers.UserRolesGroupsUpdate)
	router.Services.Gin.DELETE("/users/rolegroups/:id", router.Variables, router.IsLogged, controllers.UserRolesGroupsDelete)

	router.Services.Gin.GET("/users/resettypes", router.Variables, router.IsLogged, controllers.UserResetTypesList)
	router.Services.Gin.GET("/users/resettypes/:id", router.Variables, router.IsLogged, controllers.UserResetTypesGet)
	router.Services.Gin.POST("/users/resettypes", router.Variables, router.IsLogged, controllers.UserResetTypesCreate)
	router.Services.Gin.PUT("/users/resettypes/:id", router.Variables, router.IsLogged, controllers.UserResetTypesUpdate)
	router.Services.Gin.PATCH("/users/resettypes/:id", router.Variables, router.IsLogged, controllers.UserResetTypesUpdate)
	router.Services.Gin.DELETE("/users/resettypes/:id", router.Variables, router.IsLogged, controllers.UserResetTypesDelete)

	router.Services.Gin.GET("/users/resets", router.Variables, router.IsLogged, controllers.UserResetsList)
	router.Services.Gin.GET("/users/resets/:id", router.Variables, router.IsLogged, controllers.UserResetsGet)
	router.Services.Gin.PUT("/users/resets/:id", router.Variables, router.IsLogged, controllers.UserResetsUpdate)
	router.Services.Gin.PATCH("/resets/:id", router.Variables, router.IsLogged, controllers.UserResetsUpdate)

	router.Services.Gin.GET("/users/sessions", router.Variables, router.IsLogged, controllers.UserSessionsList)
	router.Services.Gin.GET("/users/sessions/:id", router.Variables, router.IsLogged, controllers.UserSessionsGet)

}
