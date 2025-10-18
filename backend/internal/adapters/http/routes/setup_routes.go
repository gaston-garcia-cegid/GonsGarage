package routes

import (
	"github.com/gaston-garcia-cegid/gonsgarage/backend/internal/adapters/http/handlers"
	"github.com/gaston-garcia-cegid/gonsgarage/backend/internal/adapters/http/middleware"
	"github.com/gaston-garcia-cegid/gonsgarage/backend/internal/core/ports"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	employeeHandler *handlers.EmployeeHandler,
	authService ports.AuthService,
) {
	api := router.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// Employee routes
		employees := protected.Group("/employees")
		{
			employees.POST("/", middleware.RoleMiddleware("admin", "manager"), employeeHandler.CreateEmployee)
			employees.GET("/", employeeHandler.ListEmployees)
			employees.GET("/:id", employeeHandler.GetEmployee)
			employees.PUT("/:id", middleware.RoleMiddleware("admin", "manager"), employeeHandler.UpdateEmployee)
			employees.DELETE("/:id", middleware.RoleMiddleware("admin"), employeeHandler.DeleteEmployee)
		}
	}
}
