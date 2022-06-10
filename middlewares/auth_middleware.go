package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/constants"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
)

func (m *Middlewares) IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// check if token exists in cookie
		cookie, err := c.Request.Cookie(string(constants.AUTH_TOKEN))
		if err != nil {
			// check in auth header as bearer
			authHeader := c.Request.Header.Get("Authorization")
			containsBearerToken := strings.HasPrefix(authHeader, "Bearer")
			if !containsBearerToken {
				api.SendErrorResponse(c, "Token is missing", http.StatusUnauthorized, nil)
				c.Abort()
				return
			} else {
				bearerToken := strings.Split(authHeader, "Bearer ")
				if len(bearerToken) == 2 {
					token = bearerToken[1]
				}
			}
		} else {
			token = cookie.Value
		}
		if token == "" {
			api.SendErrorResponse(c, "Not Authorized", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		// decode jwt token
		claims, err := tokens.DecodeJSONWebToken(token)

		if err != nil {
			api.SendErrorResponse(c, err.Error(), http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		if claims.Type != "access" {
			api.SendErrorResponse(c, "Refresh token is not acceptable", http.StatusNotAcceptable, nil)
			c.Abort()
			return
		}

		c.Set("id", claims.ID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// Checks user role for given route
func (m *Middlewares) UserRole(allowedRoles []constants.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {

		userRole, exists := c.Get("role")
		if !exists {
			userRole = "user"
		}
		for _, role := range allowedRoles {
			if string(role) == userRole {
				c.Next()
				return
			}
		}
		api.SendErrorResponse(c, "You don't have enough permission to access this route", http.StatusUnauthorized, nil)
		c.Abort()
	}
}
