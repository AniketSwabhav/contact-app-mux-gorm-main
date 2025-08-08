package authorization

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

var secretKey = []byte("goTeam")

type Claims struct {
	UserID   uuid.UUID
	IsAdmin  bool
	IsActive bool
	jwt.StandardClaims
}

func (c *Claims) Coder() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secretKey)
}

func ValidateToken(_ http.ResponseWriter, r *http.Request) (*Claims, error) {

	authCookie, err := r.Cookie("auth")
	tokenString := authCookie.Value
	if err != nil {
		return nil, err
	}

	token, claim, err := checkToken(tokenString)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claim, nil
}

// Checks Token String
func checkToken(tokenString string) (*jwt.Token, *Claims, error) {

	var claim = &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	return token, claim, err
}
