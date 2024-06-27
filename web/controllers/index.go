package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/web/data"
)

func IndexIndex(c *gin.Context) {

	data := data.PageData{
		Company: "Biqx Educação",
		Title:   "Início",
		Js: []string{
			"/assets/js/front-page-landing.js",
		},
	}

	c.HTML(http.StatusOK, "index/index", data)

}
