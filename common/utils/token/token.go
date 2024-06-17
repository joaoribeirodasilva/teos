package token

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joaoribeirodasilva/teos/common/configuration"
	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"github.com/joaoribeirodasilva/teos/common/service_log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID
	SessionID primitive.ObjectID
	Email     string
	Name      string
	Surname   string
}

type Token struct {
	conf        *configuration.Configuration
	User        *User
	token       *jwt.Token
	TokenString string
}

const (
	iss = "api.claudiametzgercorretora.com.br"
	aud = "users"
)

func New(conf *configuration.Configuration) *Token {

	t := &Token{}

	t.conf = conf

	return t
}

func (t *Token) Create(user *User, sessionId *primitive.ObjectID) *service_errors.Error {

	now := time.Now()

	tempKey := t.conf.GetKey("COOKIE_EXPIRE")
	if tempKey == nil || tempKey.Int == nil {
		return service_log.Error(0, http.StatusInternalServerError, "COMMON::TOKEN::Create", "", "invalid token expiration")
	}
	expires := now.Add(time.Duration(time.Second * time.Duration(*tempKey.Int)))

	tempSecret := t.conf.GetKey("SECRET_KEY")
	if tempSecret == nil || tempSecret.String == nil {
		return service_log.Error(0, http.StatusInternalServerError, "COMMON::TOKEN::Create", "", "invalid secret")
	}

	sub := make(map[string]interface{})
	sub["id"] = user.ID.Hex()
	sub["sessionId"] = sessionId.Hex()
	sub["email"] = user.Email
	sub["name"] = user.Name
	sub["surename"] = user.Surname

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": iss,
		"sub": sub,
		"aud": aud,
		"iat": now.Unix(),
		"exp": expires.Unix(),
	})

	tokenStr, err := token.SignedString([]byte(*tempSecret.String))
	if err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "COMMON::TOKEN::Create", "", "failed to encrypt token")
	}

	t.TokenString = tokenStr

	return nil
}

func (t *Token) IsValid(tokenString string) bool {

	if err := t.parseToken(tokenString); err != nil {
		return false
	}

	claims, Ok := t.token.Claims.(jwt.MapClaims)
	if !Ok || claims.Valid() != nil || !claims.VerifyAudience(aud, true) || !claims.VerifyIssuer(iss, true) {
		return false
	}

	// defer func() {
	// 	recover()
	// }()

	sub := claims["sub"].(map[string]interface{})

	iid, ok := sub["id"]
	if !ok {
		return false
	}
	strId := iid.(string)
	mid, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		return false
	}

	isessionid, ok := sub["sessionId"]
	if !ok {
		return false
	}

	strSessionId := isessionid.(string)
	sessionid, err := primitive.ObjectIDFromHex(strSessionId)
	if err != nil {
		return false
	}

	iemail, ok := sub["email"]
	if !ok {
		return false
	}
	email := iemail.(string)

	inome, ok := sub["name"]
	if !ok {
		return false
	}
	nome := inome.(string)

	isurename, ok := sub["surename"]
	if !ok {
		return false
	}
	surename := isurename.(string)

	t.User = &User{
		ID:        mid,
		SessionID: sessionid,
		Email:     email,
		Name:      nome,
		Surname:   surename,
	}

	return true
}

func (t *Token) parseToken(jwtToken string) *service_errors.Error {

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("")
		}
		return []byte(t.conf.Secret), nil
	})

	if err != nil {
		return service_log.Error(0, http.StatusInternalServerError, "COMMON::TOKEN::parseToken", "", "invalid token")
	}

	t.token = token

	return nil
}
