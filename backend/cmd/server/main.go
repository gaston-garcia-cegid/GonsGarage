package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	postgresRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/postgres"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/handlers"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/routes"
	redisRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/redis"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/usecases/auth"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/usecases/employee"
)

func main() {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://gonsgarage_user:secure_password@localhost:5432/gonsgarage?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(&domain.User{}, &domain.Employee{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Redis connection
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	userRepo := postgresRepo.NewPostgresUserRepository(db)
	employeeRepo := postgresRepo.NewPostgresEmployeeRepository(db)
	cacheRepo := redisRepo.NewRedisCacheRepository(rdb)

	// Initialize use cases
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key"
	}

	authUseCase := auth.NewAuthUseCase(userRepo, jwtSecret, 24)
	employeeUseCase := employee.NewEmployeeUseCase(employeeRepo, cacheRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	employeeHandler := handlers.NewEmployeeHandler(employeeUseCase)

	// Setup router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		routes.SetupRoutes(router, authHandler, employeeHandler, authService)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(router, authHandler, employeeHandler, authUseCase)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
