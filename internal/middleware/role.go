package middleware

import (
	"net/http"

	"github.com/aditya/finback/internal/models"
	"github.com/gin-gonic/gin"
)

// RequireRole ensures the authenticated user has one of the required roles.
func RequireRole(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User role not found in context"})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User role is invalid type"})
			return
		}

		currentRole := models.Role(roleStr)

		// Check if the current role is allowed
		authorized := false
		for _, role := range allowedRoles {
			if currentRole == role {
				authorized = true
				break
			}
		}

		if !authorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied for your role"})
			return
		}

		c.Next()
	}
}
