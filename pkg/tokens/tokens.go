package tokens

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
)

type SignedClaims struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Role string `json:"role"`
	jwt.StandardClaims
}

var secretKey []byte = []byte(env.GetEnvVariable("ACCESS_SECRET"))

func GenerateTokenWithExpiryTimeAndType(id string, expiry int64, tokenType string, role string) (string, error) {
	claims := &SignedClaims{
		Type: tokenType,
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateNoExpiryTokenWithCustomType(id, tokenType string, role string) (string, error) {
	claims := &SignedClaims{
		Type: tokenType,
		ID:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return token, nil
}

// Decode JWT Token
func DecodeJSONWebToken(tokenString string) (*SignedClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&SignedClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

	if err != nil || !token.Valid {
		return nil, errors.New("token_is_expired")
	}

	claims := token.Claims.(*SignedClaims)

	return claims, nil
}
