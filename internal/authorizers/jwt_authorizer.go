package authorizers

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Authorizer interface {
	CreateToken(id string) (string, error)
	VerifyToken(token string) (string, bool)
}

type JWTAuthorizer struct {
	JWTKey string
}

// CreateToken create the auth token from user ID
func (a *JWTAuthorizer) CreateToken(id string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(a.JWTKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyToken vefiry the auth token and grab user ID from token.
func (a *JWTAuthorizer) VerifyToken(token string) (string, bool) {
	claims := make(jwt.MapClaims)
	t, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(a.JWTKey), nil
	})
	if err != nil {
		log.Println(err)
		return "", false
	}

	if !t.Valid {
		return "", false
	}

	id, ok := claims["user_id"].(string)
	if !ok {
		return "", false
	}

	return id, true
}
