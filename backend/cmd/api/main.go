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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	postgresRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/postgres"
	_ "github.com/gaston-garcia-cegid/gonsgarage/internal/apidocs" // anclas swag /health, /ready, /metrics
	"github.com/gaston-garcia-cegid/gonsgarage/internal/platform/slogsetup"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/sqlmigrate"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/handlers"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/middleware"
	redisRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/redis"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/appointment"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/auth"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/car"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/employee"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/repair"

	_ "github.com/gaston-garcia-cegid/gonsgarage/docs" // swagger (swag)
)

// apiVersion es la versión de contrato público documentada en CHANGELOG.md y en GET /health.
const apiVersion = "1.0.0"

// @title           GonsGarage API
// @version         1.0.0
// @description     API de gestión de taller: autenticación JWT, coches, citas, reparaciones y empleados.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     JWT: cabecera Authorization con valor Bearer seguido del token (rutas bajo /api/v1 salvo /health, /ready y /metrics).
func main() {
	slogsetup.InitFromEnv()
	slogsetup.BridgeStdLog()

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
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

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

	migrationsDir := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR"))
	if migrationsDir == "" {
		migrationsDir = "db/migrations"
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("ENABLE_SQL_MIGRATIONS")), "false") {
		log.Printf("ENABLE_SQL_MIGRATIONS=false, skipping golang-migrate SQL chain")
	} else {
		log.Printf("Applying golang-migrate SQL revisions from %s", migrationsDir)
		if err := sqlmigrate.Up(dsn, migrationsDir); err != nil {
			log.Fatalf("SQL migrations failed: %v", err)
		}
		log.Printf("SQL migrations completed (or no pending revisions)")
	}

	// Auto-migrate tables in correct order (dependencies first), unless disabled for SQL-only schema workflows.
	if strings.EqualFold(strings.TrimSpace(os.Getenv("SKIP_GORM_AUTOMIGRATE")), "true") {
		log.Printf("SKIP_GORM_AUTOMIGRATE=true, skipping GORM AutoMigrate and manual indexes")
	} else {
		log.Printf("Starting GORM database migration...")

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
				continue
			}
			log.Printf("Successfully migrated %T", model)
		}

		if err := createIndexes(db); err != nil {
			log.Printf("Warning: Failed to create some indexes: %v", err)
		}

		log.Printf("GORM database migration completed successfully")
	}

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
	repairRepo := postgresRepo.NewPostgresRepairRepository(db)
	appointmentRepo := postgresRepo.NewPostgresAppointmentRepository(db)
	log.Printf("Repositories initialized")

	// Initialize use cases
	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	const defaultJWTSecret = "your-super-secret-jwt-key"
	if jwtSecret == "" {
		jwtSecret = defaultJWTSecret
		log.Printf("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GIN_MODE")), "release") && jwtSecret == defaultJWTSecret {
		log.Fatal("Refusing to start with default JWT_SECRET when GIN_MODE=release; set JWT_SECRET in the environment.")
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GIN_MODE")), "release") && strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")) == "" {
		log.Printf("Warning: CORS_ALLOWED_ORIGINS is empty in release mode; only requests without Origin (e.g. curl) or same-site usage avoid browser CORS checks.")
	}

	authService := auth.NewAuthService(userRepo, jwtSecret, 24)
	employeeService := employee.NewEmployeeService(employeeRepo, cacheRepo)
	carService := car.NewCarService(carRepo, userRepo, cacheRepo)
	appointmentService := appointment.NewAppointmentService(appointmentRepo, userRepo, carRepo)
	repairService := repair.NewRepairService(repairRepo, carRepo, userRepo)

	log.Printf("Use cases initialized")

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	carHandler := handlers.NewCarHandler(carService)

	// Initialize appointment handler
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)
	repairHandler := handlers.NewRepairHandler(repairService)

	log.Printf("Handlers initialized")

	// Setup router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.SlogRequestLogger())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// CORS middleware
	router.Use(corsMiddleware())

	// Setup routes
	setupRoutes(router, authHandler, employeeHandler, carHandler, appointmentHandler, repairHandler, authMiddleware, db)

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

// parseCORSAllowedOrigins construye el conjunto de orígenes permitidos (URLs exactas, sin barra final obligatoria).
func parseCORSAllowedOrigins(raw string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, part := range strings.Split(raw, ",") {
		u := strings.TrimSpace(part)
		if u != "" {
			out[u] = struct{}{}
		}
	}
	return out
}

// corsMiddleware: en GIN_MODE=release solo refleja Origin si está en CORS_ALLOWED_ORIGINS (coma-separada).
// Fuera de release: permisivo (cualquier Origin no vacío) para desarrollo local.
func corsMiddleware() gin.HandlerFunc {
	release := strings.EqualFold(strings.TrimSpace(os.Getenv("GIN_MODE")), "release")
	allowed := parseCORSAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if release {
			if origin != "" {
				if _, ok := allowed[origin]; !ok {
					if c.Request.Method == http.MethodOptions {
						c.AbortWithStatus(http.StatusForbidden)
						return
					}
					// Petición “simple” con Origin no listado: sin cabecera CORS el navegador bloquea la lectura.
					c.Next()
					return
				}
				c.Header("Access-Control-Allow-Origin", origin)
			}
		} else {
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
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
	repairHandler *handlers.RepairHandler,
	authMiddleware *middleware.AuthMiddleware,
	db *gorm.DB,
) {
	// Liveness: proceso arriba (orquestadores suelen usar este endpoint).
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":     "ok",
			"message":    "GonsGarage API is running",
			"apiVersion": apiVersion,
		})
	})

	// Readiness: dependencias críticas (PostgreSQL).
	router.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "checks": gin.H{"database": "unavailable"}})
			return
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			if gin.Mode() == gin.ReleaseMode {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "checks": gin.H{"database": "unavailable"}})
			} else {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "checks": gin.H{"database": err.Error()}})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready", "checks": gin.H{"database": "ok"}})
	})

	// Métricas Prometheus (runtime Go + custom futuro). No exponer a Internet sin red/proxy restringido.
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API v1 routes
	api := router.Group("/api/v1")

	// Public auth routes (per-IP rate limit unless AUTH_RATE_LIMIT_DISABLED=true).
	authLimiter := middleware.IPRateLimiterFromEnv()
	authRate := middleware.RateLimitAuth(authLimiter)
	if strings.EqualFold(strings.TrimSpace(os.Getenv("AUTH_RATE_LIMIT_DISABLED")), "true") {
		authRate = func(c *gin.Context) { c.Next() }
	}
	auth := api.Group("/auth")
	auth.Use(authRate)
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Perfil JWT: registro explícito (evita conflictos con dos RouterGroup bajo "/auth" en algunas versiones de Gin).
	api.GET("/auth/me", ginAuthMiddleware(authMiddleware), authHandler.Me)

	// Protected routes
	protected := api.Group("/")
	protected.Use(ginAuthMiddleware(authMiddleware))
	{
		// Employee routes (admin / manager only)
		employees := protected.Group("/employees")
		employees.Use(middleware.RequireStaffManagers())
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
			cars.GET("/:id/repairs", repairHandler.ListRepairsByCar)
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

		repairs := protected.Group("/repairs")
		{
			repairs.POST("", repairHandler.CreateRepair)
			repairs.GET("/:id", repairHandler.GetRepair)
			repairs.PUT("/:id", repairHandler.UpdateRepair)
		}
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
