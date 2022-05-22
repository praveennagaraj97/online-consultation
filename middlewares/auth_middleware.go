package middlewares

import (
	"fmt"
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
				api.SendErrorResponse(c, "token_not_found", http.StatusUnauthorized, nil)
				c.Abort()
				return
			} else {
				token = strings.Split(authHeader, "Bearer ")[1]
			}
		} else {
			token = cookie.Value
		}
		if token == "" {
			api.SendErrorResponse(c, "un_authorized", http.StatusUnauthorized, nil)
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
			api.SendErrorResponse(c, "refresh_token_not_acceptable", http.StatusUnauthorized, nil)
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
			fmt.Println(role, userRole)
			if string(role) == userRole {
				c.Next()
				return
			}
		}
		api.SendErrorResponse(c, "Method Not Allowed", http.StatusMethodNotAllowed, nil)
		c.Abort()
	}
}
