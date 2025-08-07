package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	 "github.com/FransHutagalung/rsbackend/controllers"
	 "github.com/FransHutagalung/rsbackend/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Add CORS middleware (basic version)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Go Backend Starter is running",
		})
	})

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// Protected routes (authentication required)
		protected := api.Group("")
		protected.Use(middleware.JWTMiddleware())
		{
			// User profile routes
			user := protected.Group("/user")
			{
				user.GET("/profile", controllers.GetProfile)
				// Add more user routes here
			}

			// Admin routes (admin role required)
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "Admin access granted",
						"data":    "This is admin-only content",
					})
				})
				// Add more admin routes here
			}
		}
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Route not found",
			"message": "The requested endpoint does not exist",
		})
	})
}