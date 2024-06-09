package cookie

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
)

type Cookie struct {
	Name     string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	server   *gin.Context
}

func New(
	server *gin.Context,
	Name string,
	Value string,
	MaxAge int,
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

	c := &Cookie{}

	tempKey := config.GetKey("COOKIE_NAME")
	if tempKey == nil || tempKey.Char == nil || *tempKey.Char == "" {
		return nil, errors.New("invalid cookie name")
	}
	c.Name = *tempKey.Char

	tempKey = config.GetKey("COOKIE_EXPIRE")
	if tempKey == nil || tempKey.Int == nil {
		return nil, errors.New("invalid cookie expire")
	}
	c.MaxAge = *tempKey.Int

	tempKey = config.GetKey("COOKIE_DOMAIN")
	if tempKey == nil || tempKey.Char == nil || *tempKey.Char == "" {
		return nil, errors.New("invalid cookie name")
	}
	c.Domain = *tempKey.Char

	tempKey = config.GetKey("COOKIE_SECURE")
	if tempKey == nil || tempKey.Int == nil {
		return nil, errors.New("invalid cookie secure")
	}
	tempSecure := false
	if *tempKey.Int != 0 {
		tempSecure = true
	}
	c.Secure = tempSecure

	tempKey = config.GetKey("COOKIE_HTTP_ONLY")
	if tempKey == nil || tempKey.Int == nil {
		return nil, errors.New("invalid cookie http only")
	}
	tempHttpOnly := false
	if *tempKey.Int != 0 {
		tempHttpOnly = true
	}
	c.HttpOnly = tempHttpOnly

	return c, nil
}

func (c *Cookie) SetCookie(server *gin.Context) {

	server.SetCookie(c.Name, c.Value, c.MaxAge, c.Path, c.Domain, c.Secure, c.HttpOnly)

}
