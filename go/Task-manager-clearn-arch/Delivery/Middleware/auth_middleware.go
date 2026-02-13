package middleware

import (
	"net/http"
	"strings"
	"task-manager-clean-arch/infrastructure"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(jwtSvc infrastructure.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims, err := jwtSvc.ValidateAccessToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		// Store the info in context for controllers
		c.Set("x-user-id", claims.ID)
		c.Set("x-user-role", claims.Role)

		c.Next()
	}
}
