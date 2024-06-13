package password

import (
	"net/http"

	"github.com/joaoribeirodasilva/teos/common/service_errors"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, *service_errors.Error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", service_errors.New(0, http.StatusUnauthorized, "CONTROLLER", "AuthReset", "", "bad password").LogError()
	}
	return string(bytes), nil

}

func Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
