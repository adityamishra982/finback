package middleware

import (
	"log"
	"net/http"

	"github.com/aditya/finback/internal/errors"
	"github.com/gin-gonic/gin"
)

// ErrorHandler is a Gin middleware that handles application errors.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors collected during the request
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Check if it's our custom AppError
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.Code, appErr)
				return
			}

			// For generic errors, return a 500 Internal Server Error
			log.Printf("Internal Server Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
		}
	}
}
