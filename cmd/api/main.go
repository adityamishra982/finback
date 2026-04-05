package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/aditya/finback/internal/api/admin/users"
	"github.com/aditya/finback/internal/api/auth"
	"github.com/aditya/finback/internal/api/dashboard"
	"github.com/aditya/finback/internal/api/records"
	"github.com/aditya/finback/internal/config"
	"github.com/aditya/finback/internal/db"
	"github.com/aditya/finback/internal/middleware"
	"github.com/aditya/finback/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	config.LoadConfig()

	// Set Gin mode
	gin.SetMode(config.GetEnv("GIN_MODE", "debug"))

	// Initialize Database Connection
	database := db.ConnectPostgres()

	// Auto Migrate the schema (temporary location until migrations are handled separately)
	err := database.AutoMigrate(&models.User{}, &models.Record{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Set Trusted Proxies
	trustedProxies := config.GetEnv("TRUSTED_PROXIES", "")
	if trustedProxies == "" {
		r.SetTrustedProxies(nil)
	} else {
		r.SetTrustedProxies(strings.Split(trustedProxies, ","))
	}

	// Global Middleware
	r.Use(middleware.ErrorHandler())

	// Handlers
	authHandler := auth.NewHandler(database)
	recordHandler := records.NewHandler(database)
	dashboardHandler := dashboard.NewHandler(database)
	userHandler := users.NewHandler(database)

	port := config.GetEnv("PORT", "8080")

	// Root Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Finback API",
			"status":  "healthy",
			"version": "v1.0.0",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Public Routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Private Routes (Require Authentication)
		private := v1.Group("")
		private.Use(middleware.AuthMiddleware())
		{
			// Financial Records Management
			// Viewers can view, Analysts can view, Admins can do everything
			recordsGroup := private.Group("/records")
			{
				recordsGroup.GET("", recordHandler.ListRecords)
				recordsGroup.POST("", middleware.RequireRole(models.RoleAdmin), recordHandler.CreateRecord)
				recordsGroup.PUT("/:id", middleware.RequireRole(models.RoleAdmin), recordHandler.UpdateRecord)
				recordsGroup.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), recordHandler.DeleteRecord)
			}

			// Dashboard Summary APIs
			// Viewers, Analysts, and Admins can all see the dashboard summary
			dashboardGroup := private.Group("/dashboard")
			{
				dashboardGroup.GET("/summary", dashboardHandler.Summary)
			}

			// User & Role Management (Admin only)
			adminGroup := private.Group("/admin")
			adminGroup.Use(middleware.RequireRole(models.RoleAdmin))
			{
				adminGroup.GET("/users", userHandler.ListUsers)
				adminGroup.PATCH("/users/:id", userHandler.UpdateUser)
			}

			// Health Check
			private.GET("/health", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":  "up",
					"message": "API is running. Connected to PostgreSQL.",
				})
			})
		}
	}

	// Run the server
	log.Printf("Starting server on :%s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
