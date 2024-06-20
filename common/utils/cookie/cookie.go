package cookie

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
)

type Cookie struct {
	config   *configuration.Configuration
	value    string
	name     string
	maxAge   int
	domain   string
	path     string
	secure   bool
	httpOnly bool
	//server   *gin.Context
}

func New(config *configuration.Configuration, tokenString string) *Cookie {

	c := &Cookie{
		config: config,
		value:  tokenString,
	}

	return c
}

func (c *Cookie) Create() error {
	var err error

	c.name, err = c.config.GetString("AUTH_COOKIE_NAME")
	if err != nil {
		return err
	}

	tempAge, err := c.config.GetInt("AUTH_COOKIE_EXPIRE")
	if err != nil {
		return err
	}
	c.maxAge = int(tempAge)

	c.domain, err = c.config.GetString("AUTH_COOKIE_DOMAIN")
	if err != nil {
		return err
	}

	c.path, err = c.config.GetString("AUTH_COOKIE_PATH")
	if err != nil {
		c.path = ""
	}

	c.secure, err = c.config.GetBool("AUTH_COOKIE_SECURE")
	if err != nil {
		return err
	}

	c.httpOnly, err = c.config.GetBool("AUTH_COOKIE_SECURE")
	if err != nil {
		return err
	}

	c.maxAge = int(time.Now().Add(time.Second * time.Duration(c.maxAge)).Unix())

	return nil
}

func (c *Cookie) Set(server *gin.Context) {

	server.SetCookie(c.name, c.value, c.maxAge, c.path, c.domain, c.secure, c.httpOnly)

}

func (c *Cookie) SetEmpty(server *gin.Context) error {

	name, err := c.config.GetString("AUTH_COOKIE_NAME")
	if err != nil {
		return err
	}

	server.SetCookie(name, "", 0, "", "", false, false)

	return nil
}
