package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak: Anda tidak memiliki izin Admin!"})
			return
		}

		c.Next()
	}
}
