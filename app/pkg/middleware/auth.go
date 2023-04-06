package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authencate(authHeader string, skipPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authHeader != c.Request.Header.Get("x-auth-token") {
			for _, path := range skipPaths {
				if path == c.Request.URL.Path {
					c.Next()
					return
				}
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}
