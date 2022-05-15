package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var trustedDomains map[string]bool = map[string]bool{
	"http://localhost:4000": true,
	"http://localhost:3000": true,
}

func (m *Middlewares) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Server", "online-consultation")

		if trustedDomains[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
