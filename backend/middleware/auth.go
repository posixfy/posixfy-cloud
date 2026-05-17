package middleware

import (
	"net/http"
	"strings"

	"posixfy-cloud/backend/auth"

	"github.com/gin-gonic/gin"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// Try Authorization header first
		header := c.GetHeader("Authorization")
		if header != "" && strings.HasPrefix(header, "Bearer ") {
			tokenStr = strings.TrimPrefix(header, "Bearer ")
		} else if t := c.Query("token"); t != "" {
			// Fallback to query param (for EventSource which can't set custom headers)
			tokenStr = t
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		claims, err := auth.ParseToken(secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) *auth.Claims {
	v, _ := c.Get("claims")
	return v.(*auth.Claims)
}
