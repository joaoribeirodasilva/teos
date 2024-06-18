package cookie

import (
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
)

type Cookie struct {
	Name     string
	Value    string
	MaxAge   int64
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	//server   *gin.Context
}

func New(
	server *gin.Context,
	Name string,
	Value string,
	MaxAge int64,
	Domain string,
	Secure int,
	HttpOnly int) *Cookie {

	c := &Cookie{
		Name:     Name,
		Value:    Value,
		MaxAge:   MaxAge,
		Path:     "/",
		Domain:   Domain,
		Secure:   false,
		HttpOnly: false,
	}

	tempSecure := false
	if Secure != 0 {
		tempSecure = true
	}
	c.Secure = tempSecure

	tempHttpOnly := false
	if HttpOnly != 0 {
		tempHttpOnly = true
	}
	c.HttpOnly = tempHttpOnly

	return c
}

func NewFromConfiguration(value string, config *configuration.Configuration) (*Cookie, error) {

	var err error
	c := &Cookie{}

	c.Name, err = config.GetString("COOKIE_NAME")
	if err != nil {
		return nil, err
	}

	c.MaxAge, err = config.GetInt("COOKIE_EXPIRE")
	if err != nil {
		return nil, err
	}

	c.Domain, err = config.GetString("COOKIE_DOMAIN")
	if err != nil {
		return nil, err
	}

	c.Secure, err = config.GetBool("COOKIE_SECURE")
	if err != nil {
		return nil, err
	}

	c.HttpOnly, err = config.GetBool("COOKIE_SECURE")
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Cookie) SetCookie(server *gin.Context) {

	server.SetCookie(c.Name, c.Value, int(c.MaxAge), c.Path, c.Domain, c.Secure, c.HttpOnly)

}
