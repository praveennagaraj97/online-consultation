package tokens

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/constants"
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
		return nil, errors.New("Provided token is expired or invalid")
	}

	claims := token.Claims.(*SignedClaims)

	return claims, nil
}

// returns access, refresh and error
func GetAccessAndRefreshTokenFromRequest(c *gin.Context) (string, string, error) {
	var token string
	var refresh_token string

	// Get auth token either from cookie nor Header
	cookie, err := c.Request.Cookie(string(constants.AUTH_TOKEN))
	if err != nil {
		// check in auth header as bearer
		authHeader := c.Request.Header.Get("Authorization")
		containsBearerToken := strings.HasPrefix(authHeader, "Bearer")
		if !containsBearerToken {
			token = ""
		} else {
			token = strings.Split(authHeader, "Bearer ")[1]
		}
	} else {
		token = cookie.Value
	}

	// Get refresh token either from cookie nor params
	refreshCookie, err := c.Request.Cookie(string(constants.REFRESH_TOKEN))
	if err != nil {
		// check in auth header as bearer
		refresh_token = c.Request.URL.Query().Get("refresh_token")
		if refresh_token == "" {
			return "", "", errors.New("refresh token is missing")
		}

	} else {
		refresh_token = refreshCookie.Value
	}

	return token, refresh_token, nil
}
