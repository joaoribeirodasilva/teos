package routes

import (
	"github.com/joaoribeirodasilva/teos/common/server"
	"github.com/joaoribeirodasilva/teos/users/controllers"
)

func RegisterRoutes(router *server.Router) {

	router.Service.Http.Engine.GET("/users/users", router.Services, router.IsLogged, controllers.UsersList)
	router.Service.Http.Engine.GET("/users/users/:id", router.Services, router.IsLogged, controllers.UsersGet)
	router.Service.Http.Engine.POST("/users/users", router.Services, router.IsLogged, controllers.UsersCreate)
	router.Service.Http.Engine.PUT("/users/users/:id", router.Services, router.IsLogged, controllers.UsersUpdate)
	router.Service.Http.Engine.PATCH("/users/users/:id", router.Services, router.IsLogged, controllers.UsersUpdate)
	router.Service.Http.Engine.DELETE("/users/users/:id", router.Services, router.IsLogged, controllers.UsersDelete)

	router.Service.Http.Engine.GET("/users/groups", router.Services, router.IsLogged, controllers.UserGroupsList)
	router.Service.Http.Engine.GET("/users/groups/:id", router.Services, router.IsLogged, controllers.UserGroupsGet)
	router.Service.Http.Engine.POST("/users/groups", router.Services, router.IsLogged, controllers.UserGroupsCreate)
	router.Service.Http.Engine.PUT("/users/groups/:id", router.Services, router.IsLogged, controllers.UserGroupsUpdate)
	router.Service.Http.Engine.PATCH("/users/groups/:id", router.Services, router.IsLogged, controllers.UserGroupsUpdate)
	router.Service.Http.Engine.DELETE("/users/groups/:id", router.Services, router.IsLogged, controllers.UserGroupsDelete)

	router.Service.Http.Engine.GET("/users/organizations", router.Services, router.IsLogged, controllers.UserOrganizationsList)
	router.Service.Http.Engine.GET("/users/organizations/:id", router.Services, router.IsLogged, controllers.UserOrganizationsGet)
	router.Service.Http.Engine.POST("/users/organizations", router.Services, router.IsLogged, controllers.UserOrganizationsCreate)
	router.Service.Http.Engine.PUT("/users/organizations/:id", router.Services, router.IsLogged, controllers.UserOrganizationsUpdate)
	router.Service.Http.Engine.PATCH("/users/organizations/:id", router.Services, router.IsLogged, controllers.UserOrganizationsUpdate)
	router.Service.Http.Engine.DELETE("/users/organizations/:id", router.Services, router.IsLogged, controllers.UserOrganizationsDelete)

	router.Service.Http.Engine.GET("/users/permissions", router.Services, router.IsLogged, controllers.UserPermissionsList)
	router.Service.Http.Engine.GET("/users/permissions/:id", router.Services, router.IsLogged, controllers.UserPermissionsGet)
	router.Service.Http.Engine.POST("/users/permissions", router.Services, router.IsLogged, controllers.UserPermissionsCreate)
	router.Service.Http.Engine.PUT("/users/permissions/:id", router.Services, router.IsLogged, controllers.UserPermissionsUpdate)
	router.Service.Http.Engine.PATCH("/users/permissions/:id", router.Services, router.IsLogged, controllers.UserPermissionsUpdate)
	router.Service.Http.Engine.DELETE("/users/permissions/:id", router.Services, router.IsLogged, controllers.UserPermissionsDelete)

}
