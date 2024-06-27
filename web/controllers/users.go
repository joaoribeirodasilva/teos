package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/web/data"
)

func UsersProfile(c *gin.Context) {

	data := data.PageData{
		Company: "Biqx Educação",
		Title:   "Usuário",
		Js: []string{
			"assets/js/pages-profile.js",
		},
		Css: []string{
			"/assets/vendor/css/pages/page-profile.css",
		},
	}

	c.HTML(http.StatusOK, "admin/user_profile", data)

}
