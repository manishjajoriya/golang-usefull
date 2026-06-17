package middleware

import (
	"NoRethink/internal/auth"
	"NoRethink/internal/config"
	"NoRethink/internal/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwt *config.JwtConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if accessCookie, err := c.Request.Cookie(util.AccessTokenKey); err == nil && accessCookie != nil {
			isValid, userId, appErr := auth.ValidateAccessToken(accessCookie.Value, jwt.AccessSecret)
			if isValid && appErr == nil {
				c.Set(util.UserIDKey, userId)

				c.Next()
				return
			}
		}
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "missing authorization header",
				"success": false,
			})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid authorization format",
				"success": false,
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid authorization format",
				"success": false,
			})
			return
		}

		isValid, userId, appErr := auth.ValidateAccessToken(tokenString, jwt.AccessSecret)

		if appErr != nil || !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"success": false,
			})
			return
		}

		c.Set(util.UserIDKey, userId)

		c.Next()
	}
}
