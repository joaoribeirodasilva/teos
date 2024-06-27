package payload

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/models"
	"github.com/joaoribeirodasilva/teos/common/redisdb"
	"github.com/joaoribeirodasilva/teos/common/utils/permission_key"
)

type SessionAuth struct {
	ID             uint
	UserID         uint
	OrganizationID uint
	Email          string
	Name           string
	Surname        string
	AvatarUrl      string
}

func EmptySession() *SessionAuth {
	sa := SessionAuth{
		ID:             0,
		OrganizationID: 1,
		UserID:         1,
		Name:           "Sistema",
		Surname:        "Teos",
		Email:          "teos@teos.com.br",
		AvatarUrl:      "",
	}
	return &sa
}

type TokenClaims struct {
	iss string
	aud string
}

type HttpCookie struct {
	conf          *configuration.ConfigCookie
	ctx           *gin.Context
	tokenClaims   TokenClaims
	UserSession   *SessionAuth
	Ttl           int
	permissionsDb *redisdb.RedisDB
	sessionsDb    *redisdb.RedisDB
	method        string
	route         string
}

func NewHttpCookie(conf *configuration.ConfigCookie, ctx *gin.Context, permissionsDb *redisdb.RedisDB, sessionsDb *redisdb.RedisDB, method string, route string) *HttpCookie {

	// Apply to default user session
	return &HttpCookie{
		ctx:  ctx,
		conf: conf,
		tokenClaims: TokenClaims{
			iss: conf.Domain,
			aud: "users",
		},
		Ttl:           conf.MaxAge,
		permissionsDb: permissionsDb,
		sessionsDb:    sessionsDb,
		method:        method,
		route:         route,
	}
}

func (c *HttpCookie) Parse() error {

	if c.UserSession == nil {
		c.UserSession = EmptySession()
	}

	noCookie := false
	jwt, err := c.ctx.Cookie(c.conf.Name)
	if err != nil || jwt == "" {
		if !errors.Is(err, http.ErrNoCookie) {
			return err
		}
		noCookie = true
	}

	routeRey := permission_key.MakePermissionKey(c.method, c.route)
	routeVal, err := c.permissionsDb.Get(routeRey)
	if err != nil || routeVal == nil {
		return err
	}

	permission := models.SvcRoute{}
	if err := json.Unmarshal([]byte(*routeVal), &permission); err != nil {
		return err
	}

	if noCookie && permission.Open != 0 {
		return nil
	}

	if err := c.getTokenClaims(jwt); err != nil {
		return err
	}

	if permission.Open != 0 {
		return nil
	}

	sessionVal, err := c.sessionsDb.Get(fmt.Sprintf("%d", c.UserSession.ID))
	if err != nil || sessionVal == nil {
		return errors.New("invalid session")
	}

	session := models.AuthSession{}
	if err := json.Unmarshal([]byte(*sessionVal), &session); err != nil {
		return err
	}

	if session.ID != c.UserSession.ID {
		return errors.New("invalid session")
	}

	// TODO: validate if user has permission for this route

	return nil
}

func (c *HttpCookie) SetCookie() error {

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

	c.ctx.SetCookie(c.conf.Name, tokenStr, c.conf.MaxAge, "/", c.conf.Domain, c.conf.Secure, c.conf.HttpOnly)

	return nil
}

func (c *HttpCookie) SetEmptyCookie() {
	c.ctx.SetCookie(c.conf.Name, "", c.conf.MaxAge, "/", c.conf.Domain, c.conf.Secure, c.conf.HttpOnly)
}

func (c *HttpCookie) getTokenClaims(jwtToken string) error {

	retErr := fmt.Errorf("invalid auth token")

	token, err := c.parseToken(jwtToken)
	if err != nil {
		return err
	}

	claims, Ok := token.Claims.(jwt.MapClaims)
	/* 	valid := claims.Valid()
	   	audience := claims.VerifyAudience(c.tokenClaims.aud, true)
	   	issuer := claims.VerifyIssuer(c.tokenClaims.iss, true) */

	//fmt.Printf("ok: %s, aud: %t, issuer: %t\n", valid, audience, issuer)
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

	userSession := SessionAuth{}

	interfaceId, ok := sub["id"]
	if !ok {
		return retErr
	}
	tempFloat, ok := interfaceId.(float64)
	if !ok {
		return errors.New("invalid token (id)")
	}
	userSession.ID = uint(tempFloat)

	interfaceUserId, ok := sub["userId"]
	if !ok {
		return retErr
	}
	tempFloat, ok = interfaceUserId.(float64)
	if !ok {
		return errors.New("invalid token (user id)")
	}
	userSession.UserID = uint(tempFloat)

	interfaceOrganizationId, ok := sub["organizationId"]
	if !ok {
		return retErr
	}
	tempFloat, ok = interfaceOrganizationId.(float64)
	if !ok {
		return errors.New("invalid token (org id)")
	}
	userSession.OrganizationID = uint(tempFloat)

	interfaceEmail, ok := sub["email"]
	if !ok {
		return errors.New("invalid token (email)")
	}
	userSession.Email = interfaceEmail.(string)

	interfaceName, ok := sub["name"]
	if !ok {
		return errors.New("invalid token (name)")
	}

	userSession.Name = interfaceName.(string)

	interfaceSurname, ok := sub["surname"]
	if !ok {
		return errors.New("invalid token (surname)")
	}
	userSession.Surname = interfaceSurname.(string)

	interfaceAvatarUrl, ok := sub["avatarUrl"]
	if !ok {
		return errors.New("invalid token (avatar)")
	}
	userSession.AvatarUrl = interfaceAvatarUrl.(string)

	c.UserSession = &userSession

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
