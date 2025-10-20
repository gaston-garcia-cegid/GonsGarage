package routes

import (
	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/handlers"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	authGroup := router.Group("/api/v1/auth")
	{
		// Public routes
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)

		// Protected routes (if you need them later)
		// protected := authGroup.Group("/")
		// protected.Use(middleware.AuthMiddleware())
		// {
		// 	protected.POST("/refresh", authHandler.RefreshToken)
		// 	protected.POST("/logout", authHandler.Logout)
		// }
	}
}

func SetupEmployeeRoutes(router *gin.Engine, employeeHandler *handlers.EmployeeHandler) {
	employeeGroup := router.Group("/api/v1/employees")
	employeeGroup.Use(middleware.AuthMiddleware()) // Protect all employee routes
	{
		employeeGroup.POST("/", employeeHandler.CreateEmployee)
		employeeGroup.GET("/", employeeHandler.ListEmployees)
		employeeGroup.GET("/:id", employeeHandler.GetEmployee)
		employeeGroup.PUT("/:id", employeeHandler.UpdateEmployee)
		employeeGroup.DELETE("/:id", employeeHandler.DeleteEmployee)
	}
}

func SetupAllRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, employeeHandler *handlers.EmployeeHandler) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "GonsGarage API is running",
		})
	})

	// Setup all route groups
	SetupAuthRoutes(router, authHandler)
	SetupEmployeeRoutes(router, employeeHandler)
}
