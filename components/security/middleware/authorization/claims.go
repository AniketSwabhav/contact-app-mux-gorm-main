package authorization

import (
	"contact_app_mux_gorm_main/components/apperror"
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

func ValidateToken(_ http.ResponseWriter, r *http.Request, claim *Claims) error {

	authCookie, err := r.Cookie("auth")
	if err != nil {
		return apperror.NewUnAuthorizedError("missing or invalid auth cookie")
	}

	tokenString := authCookie.Value
	if tokenString == "" {
		return apperror.NewUnAuthorizedError("empty token")
	}

	token, err := checkToken(tokenString, claim)
	if err != nil {
		return apperror.NewUnAuthorizedError("invalid token: " + err.Error())
	}

	if !token.Valid {
		return apperror.NewInValidTokenError("invalid token")
	}

	return nil
}

// Checks Token String
func checkToken(tokenString string, claim *Claims) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	return token, err
}

// func ValidateToken(_ http.ResponseWriter, r *http.Request) (*Claims, error) {

// 	authCookie, err := r.Cookie("auth")
// 	if err != nil {
// 		return nil, apperror.NewUnAuthorizedError("missing or invalid auth cookie")
// 	}

// 	tokenString := authCookie.Value
// 	if tokenString == "" {
// 		return nil, apperror.NewUnAuthorizedError("empty token")
// 	}

// 	token, claim, err := checkToken(tokenString)
// 	if err != nil {
// 		return nil, apperror.NewUnAuthorizedError("invalid token: " + err.Error())
// 	}

// 	if !token.Valid {
// 		return nil, apperror.NewInValidTokenError("invalid token")
// 	}

// 	return claim, nil
// }

// // Checks Token String
// func checkToken(tokenString string) (*jwt.Token, *Claims, error) {

// 	var claim = &Claims{}
// 	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})

// 	return token, claim, err
// }
