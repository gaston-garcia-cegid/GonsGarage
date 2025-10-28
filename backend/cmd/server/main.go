package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	postgresRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/postgres"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/handlers"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/middleware"
	redisRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/redis"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/appointment"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/auth"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/car"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/employee"
)

// @title GonsGarage API
// @version 1.0
// @description Auto repair shop management API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	log.Printf("/*************** Start Main ***************/")

	// Database connection with timeout
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://admindb:gonsgarage123@localhost:5432/gonsgarage?sslmode=disable"
	}
	log.Printf("Connecting to database: %s", dsn)

	// Configure GORM with better logging and timeout
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test database connection with timeout
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Printf("Database connection established")

	// Check if we need to reset database (for development)
	resetDB := os.Getenv("RESET_DATABASE")
	if resetDB == "true" {
		log.Printf("RESET_DATABASE=true, dropping existing tables...")
		if err := dropAllTables(db); err != nil {
			log.Printf("Warning: Failed to drop tables: %v", err)
		}
	}

	// Auto-migrate tables in correct order (dependencies first)
	log.Printf("Starting database migration...")

	// Migrate tables one by one in dependency order
	models := []interface{}{
		&domain.User{},
		&domain.Employee{},
		&domain.Car{},
		&domain.Repair{},
		&domain.Appointment{},
	}

	for _, model := range models {
		log.Printf("Migrating %T...", model)
		if err := db.AutoMigrate(model); err != nil {
			log.Printf("Migration error for %T: %v", model, err)
			// Continue with other migrations instead of failing completely
			continue
		}
		log.Printf("Successfully migrated %T", model)
	}

	// Create indexes manually if they don't exist
	if err := createIndexes(db); err != nil {
		log.Printf("Warning: Failed to create some indexes: %v", err)
	}

	log.Printf("Database migration completed successfully")

	// Redis connection with timeout
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test Redis connection with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cacheRepo ports.CacheRepository
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Printf("Continuing without Redis cache...")
		cacheRepo = &NullCacheRepository{}
	} else {
		log.Printf("Redis connection established")
		cacheRepo = redisRepo.NewRedisCacheRepository(rdb)
	}

	// Initialize repositories
	userRepo := postgresRepo.NewPostgresUserRepository(db)
	employeeRepo := postgresRepo.NewPostgresEmployeeRepository(db)
	carRepo := postgresRepo.NewPostgresCarRepository(db)
	appointmentRepo := postgresRepo.NewPostgresAppointmentRepository(db)
	log.Printf("Repositories initialized")

	// Initialize use cases
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key"
		log.Printf("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	authService := auth.NewAuthService(userRepo, jwtSecret, 24)
	employeeService := employee.NewEmployeeService(employeeRepo, cacheRepo)
	carService := car.NewCarService(carRepo, userRepo, cacheRepo)
	appointmentService := appointment.NewAppointmentService(appointmentRepo, userRepo, cacheRepo)

	log.Printf("Use cases initialized")

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	carHandler := handlers.NewCarHandler(carService)

	// Initialize appointment handler
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)

	log.Printf("Handlers initialized")

	// Setup router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// CORS middleware
	router.Use(corsMiddleware())

	// Setup routes
	setupRoutes(router, authHandler, employeeHandler, carHandler, appointmentHandler, authMiddleware)

	log.Printf("Routes set up")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Drop all tables for clean migration
func dropAllTables(db *gorm.DB) error {
	tables := []string{
		"appointments",
		"repairs",
		"cars",
		"employees",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
			log.Printf("Failed to drop table %s: %v", table, err)
		} else {
			log.Printf("Dropped table: %s", table)
		}
	}
	return nil
}

// Create indexes manually
func createIndexes(db *gorm.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_employees_user_id ON employees(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON employees(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_cars_owner_id ON cars(owner_id)",
		"CREATE INDEX IF NOT EXISTS idx_cars_license_plate ON cars(license_plate)",
		"CREATE INDEX IF NOT EXISTS idx_cars_deleted_at ON cars(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_repairs_car_id ON repairs(car_id)",
		"CREATE INDEX IF NOT EXISTS idx_repairs_technician_id ON repairs(technician_id)",
		"CREATE INDEX IF NOT EXISTS idx_repairs_deleted_at ON repairs(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_appointments_customer_id ON appointments(customer_id)",
		"CREATE INDEX IF NOT EXISTS idx_appointments_car_id ON appointments(car_id)",
		"CREATE INDEX IF NOT EXISTS idx_appointments_deleted_at ON appointments(deleted_at)",
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			log.Printf("Failed to create index: %s, error: %v", idx, err)
		}
	}
	return nil
}

// CORS middleware function
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow specific origins or all origins for development
		if origin == "http://localhost:3000" || origin == "http://localhost:3001" || os.Getenv("GIN_MODE") != "release" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Setup all routes
func setupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	employeeHandler *handlers.EmployeeHandler,
	carHandler *handlers.CarHandler,
	appointmentHandler *handlers.AppointmentHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "GonsGarage API is running",
		})
	})

	// API v1 routes
	api := router.Group("/api/v1")

	// Public auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(ginAuthMiddleware(authMiddleware))
	{
		// Employee routes
		employees := protected.Group("/employees")
		{
			employees.POST("", employeeHandler.CreateEmployee)
			employees.GET("", employeeHandler.ListEmployees)
			employees.GET("/:id", employeeHandler.GetEmployee)
			employees.PUT("/:id", employeeHandler.UpdateEmployee)
			employees.DELETE("/:id", employeeHandler.DeleteEmployee)
		}

		// Car routes would go here
		cars := protected.Group("/cars")
		{
			cars.POST("", carHandler.CreateCar)
			cars.GET("", carHandler.ListCars)
			cars.GET("/:id", carHandler.GetCar)
			cars.PUT("/:id", carHandler.UpdateCar)
			cars.DELETE("/:id", carHandler.DeleteCar)
		}

		// Appointment routes would go here
		appointments := protected.Group("/appointments")
		{
			appointments.POST("", appointmentHandler.CreateAppointment)
			appointments.GET("", appointmentHandler.ListAppointments)
			appointments.GET("/:id", appointmentHandler.GetAppointment)
			appointments.PUT("/:id", appointmentHandler.UpdateAppointment)
			appointments.DELETE("/:id", appointmentHandler.DeleteAppointment)
		}

		// Repair routes would go here
		// repairs := protected.Group("/repairs")
		// {
		// 	repairs.POST("", repairHandler.CreateRepair)
		// 	repairs.GET("", repairHandler.ListRepairs)
		// 	repairs.GET("/:id", repairHandler.GetRepair)
		// 	repairs.PUT("/:id", repairHandler.UpdateRepair)
		// 	repairs.DELETE("/:id", repairHandler.DeleteRepair)
		// }
	}
}

// Convert auth middleware to Gin middleware
// ginAuthMiddleware converts the auth middleware to work properly with Gin
func ginAuthMiddleware(authMiddleware *middleware.AuthMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header using Gin's method
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if token has Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate token using jwt
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(authMiddleware.GetJWTSecret()), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Extract user ID
		var userIDStr string
		if uid, exists := claims["userID"]; exists {
			if uidStr, ok := uid.(string); ok {
				userIDStr = uidStr
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID in token"})
				c.Abort()
				return
			}
		} else if sub, exists := claims["sub"]; exists {
			if subStr, ok := sub.(string); ok {
				userIDStr = subStr
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid sub in token"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user identifier in token"})
			c.Abort()
			return
		}

		// Validate UUID format
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID format"})
			c.Abort()
			return
		}

		// Store user info in Gin context (much cleaner than http.Request context)
		c.Set("userID", userID.String())

		if email, exists := claims["email"]; exists {
			if emailStr, ok := email.(string); ok {
				c.Set("userEmail", emailStr)
			}
		}

		if role, exists := claims["role"]; exists {
			if roleStr, ok := role.(string); ok {
				c.Set("userRole", roleStr)
			}
		}

		log.Printf("✅ Authentication successful for user: %s", userID.String())

		// Continue to next handler
		c.Next()
	}
}

// Helper function to get user ID from Gin context
func getUserIDFromGinContext(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}

	userIDString, ok := userIDStr.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID is not a string")
	}

	return uuid.Parse(userIDString)
}

// ResponseWriter wrapper to capture status codes
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// NullCacheRepository é uma implementação dummy quando Redis não está disponível
type NullCacheRepository struct{}

func (n *NullCacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	return nil // Não faz nada
}

func (n *NullCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	return nil // Não faz nada
}

func (n *NullCacheRepository) Delete(ctx context.Context, key string) error {
	return nil // Não faz nada
}
