package payload

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
)

type SessionAuth struct {
	ID             uint
	UserID         uint
	Email          string
	Name           string
	Surname        string
	OrganizationID uint
	AvatarUrl      string
}

type TokenClaims struct {
	iss string
	aud string
}

type HttpCookie struct {
	conf        *configuration.ConfigCookie
	ctx         *gin.Context
	tokenClaims TokenClaims
	UserSession *SessionAuth
}

func NewHttpCookie(conf *configuration.ConfigCookie, ctx *gin.Context) *HttpCookie {

	// Apply to default user session
	return &HttpCookie{
		ctx:  ctx,
		conf: conf,
		tokenClaims: TokenClaims{
			iss: conf.Domain,
			aud: "users",
		},
	}
}

func (c *HttpCookie) Parse() error {

	jwt, err := c.ctx.Cookie("gin_cookie")
	if err != nil {
		return err
	}

	if err := c.getTokenClaims(jwt); err != nil {
		return err
	}

	return nil
}

func (c *HttpCookie) Set() error {

	sub := make(map[string]interface{})
	sub["id"] = c.UserSession.ID
	sub["userId"] = c.UserSession.UserID
	sub["organizationId"] = c.UserSession.OrganizationID
	sub["email"] = c.UserSession.Email
	sub["name"] = c.UserSession.Name
	sub["surname"] = c.UserSession.Surname
	sub["avatarUrl"] = c.UserSession.AvatarUrl

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": c.tokenClaims.iss,
		"sub": sub,
		"aud": c.tokenClaims.aud,
		"iat": time.Now().Unix(),
	})

	tokenStr, err := token.SignedString([]byte(c.conf.Secret))
	if err != nil {
		return fmt.Errorf("failed to encrypt token")
	}

	c.ctx.SetCookie("gin_cookie", tokenStr, c.conf.MaxAge, "/", c.conf.Domain, c.conf.Secure, c.conf.HttpOnly)

	return nil
}

func (c *HttpCookie) getTokenClaims(jwtToken string) error {

	retErr := fmt.Errorf("invalid auth token")

	token, err := c.parseToken(jwtToken)
	if err != nil {
		return err
	}

	claims, Ok := token.Claims.(jwt.MapClaims)
	if !Ok || claims.Valid() != nil || !claims.VerifyAudience(c.tokenClaims.aud, true) || !claims.VerifyIssuer(c.tokenClaims.iss, true) {
		return retErr
	}

	tempSub, ok := claims["sub"]
	if !ok {
		return retErr
	}

	sub, ok := tempSub.(map[string]interface{})
	if !ok {
		return retErr
	}

	interfaceId, ok := sub["id"]
	if !ok {
		return retErr
	}
	c.UserSession.ID = interfaceId.(uint)

	interfaceUserId, ok := sub["userId"]
	if !ok {
		return retErr
	}
	c.UserSession.ID = interfaceUserId.(uint)

	interfaceOrganizationId, ok := sub["organizationId"]
	if !ok {
		return retErr
	}
	c.UserSession.OrganizationID = interfaceOrganizationId.(uint)

	interfaceEmail, ok := sub["email"]
	if !ok {
		return retErr
	}
	c.UserSession.Email = interfaceEmail.(string)

	interfaceName, ok := sub["name"]
	if !ok {
		return retErr
	}
	c.UserSession.Name = interfaceName.(string)

	interfaceSurname, ok := sub["surname"]
	if !ok {
		return retErr
	}
	c.UserSession.Surname = interfaceSurname.(string)

	interfaceAvatarUrl, ok := sub["avatarUrl"]
	if !ok {
		return retErr
	}
	c.UserSession.AvatarUrl = interfaceAvatarUrl.(string)

	return nil

}

func (c *HttpCookie) parseToken(jwtToken string) (*jwt.Token, error) {

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("")
		}
		return []byte(c.conf.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil

}
