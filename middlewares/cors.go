package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var trustedDomains map[string]bool = map[string]bool{
	"https://online-consultation.vercel.app": true,
	"https://gmg-web.vercel.app":             true,
}

var allowedHeaders = []string{
	"User-Agent",
	"Pragma",
	"Host",
	"Connection",
	"Cache-Control",
	"Accept-Language",
	"Accept-Encoding",
	"Accept",
	"Time-Zone",
	"Content-Type",
	"Authorization",
}

func (m *Middlewares) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))

		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Server", "online-consultation")

		if trustedDomains[origin] || strings.Contains(origin, "http://localhost") {
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
