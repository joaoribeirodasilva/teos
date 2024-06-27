package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/web/data"
)

func AuthLogin(c *gin.Context) {

	data := data.PageData{
		Company: "Biqx Educação",
		Title:   "Login",
		Js: []string{
			"/assets/js/pages-auth.js",
		},
	}

	c.HTML(http.StatusOK, "pages/login", data)

}

func AuthForgot(c *gin.Context) {

	data := data.PageData{
		Company: "Biqx Educação",
		Title:   "Esqueceu Senha",
		Js: []string{
			"/assets/js/pages-auth.js",
		},
	}

	c.HTML(http.StatusOK, "pages/forgot", data)
}

func AuthSignup(c *gin.Context) {

	data := data.PageData{
		Company: "Biqx Educação",
		Title:   "Nova Conta",
		Js: []string{
			"/assets/js/pages-auth.js",
		},
	}

	c.HTML(http.StatusOK, "pages/signup", data)
}

func AuthReset(c *gin.Context) {

	data := data.PageData{
		Company: "Biqx Educação",
		Title:   "Reset de Senha",
		Js: []string{
			"/assets/js/pages-auth.js",
		},
	}

	c.HTML(http.StatusOK, "pages/reset", data)
}
