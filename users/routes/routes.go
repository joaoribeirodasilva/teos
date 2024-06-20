package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/users/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Services.Gin.GET("/users/users", router.Variables, router.IsLogged, controllers.UsersList)
	router.Services.Gin.GET("/users/users/:id", router.Variables, router.IsLogged, controllers.UsersGet)
	router.Services.Gin.POST("/users/users", router.Variables, router.IsLogged, controllers.UsersCreate)
	router.Services.Gin.PUT("/users/users/:id", router.Variables, router.IsLogged, controllers.UsersUpdate)
	router.Services.Gin.PATCH("/users/users/:id", router.Variables, router.IsLogged, controllers.UsersUpdate)
	router.Services.Gin.DELETE("/users/users/:id", router.Variables, router.IsLogged, controllers.UsersDelete)

	router.Services.Gin.GET("/users/groups", router.Variables, router.IsLogged, controllers.UserGroupsList)
	router.Services.Gin.GET("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsGet)
	router.Services.Gin.POST("/users/groups", router.Variables, router.IsLogged, controllers.UserGroupsCreate)
	router.Services.Gin.PUT("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsUpdate)
	router.Services.Gin.PATCH("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsUpdate)
	router.Services.Gin.DELETE("/users/groups/:id", router.Variables, router.IsLogged, controllers.UserGroupsDelete)

	router.Services.Gin.GET("/users/organizations", router.Variables, router.IsLogged, controllers.UserOrganizationsList)
	router.Services.Gin.GET("/users/organizations/:id", router.Variables, router.IsLogged, controllers.UserOrganizationsGet)
	router.Services.Gin.POST("/users/organizations", router.Variables, router.IsLogged, controllers.UserOrganizationsCreate)
	router.Services.Gin.PUT("/users/organizations/:id", router.Variables, router.IsLogged, controllers.UserOrganizationsUpdate)
	router.Services.Gin.PATCH("/users/organizations/:id", router.Variables, router.IsLogged, controllers.UserOrganizationsUpdate)
	router.Services.Gin.DELETE("/users/organizations/:id", router.Variables, router.IsLogged, controllers.UserOrganizationsDelete)

	router.Services.Gin.GET("/users/permissions", router.Variables, router.IsLogged, controllers.UserPermissionsList)
	router.Services.Gin.GET("/users/permissions/:id", router.Variables, router.IsLogged, controllers.UserPermissionsGet)
	router.Services.Gin.POST("/users/permissions", router.Variables, router.IsLogged, controllers.UserPermissionsCreate)
	router.Services.Gin.PUT("/users/permissions/:id", router.Variables, router.IsLogged, controllers.UserPermissionsUpdate)
	router.Services.Gin.PATCH("/users/permissions/:id", router.Variables, router.IsLogged, controllers.UserPermissionsUpdate)
	router.Services.Gin.DELETE("/users/permissions/:id", router.Variables, router.IsLogged, controllers.UserPermissionsDelete)

}
