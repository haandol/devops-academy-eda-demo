package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authencate(authHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authHeader != c.Request.Header.Get("x-auth-token") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}
