package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/web/data"
)

func AdminIndex(c *gin.Context) {

	data := data.PageData{
		Company: "Biqx Educação",
		Title:   "Dashboard",
		Js:      []string{},
	}

	c.HTML(http.StatusOK, "admin/index", data)

}
