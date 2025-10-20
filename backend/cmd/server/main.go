package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	postgresRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/postgres"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/handlers"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/routes"
	redisRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/redis"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/usecases/auth"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/usecases/employee"
)

func main() {
	log.Printf("/*************** Start Main ***************/")

	// Database connection with timeout
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://admindb:gonsgarage123@localhost:5432/gonsgarage?sslmode=disable"
	}
	log.Printf("Connecting to database: %s", dsn)
	log.Printf("/******************************************/")

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
	log.Printf("/******************************************/")

	// Test raw SQL
	log.Printf("Testing raw SQL...")
	var version string
	if err := db.Raw("SELECT version()").Scan(&version).Error; err != nil {
		log.Fatal("Failed to execute raw SQL:", err)
	}
	log.Printf("PostgreSQL version: %s", version)

	// Check schema
	var currentSchema string
	if err := db.Raw("SELECT current_schema()").Scan(&currentSchema).Error; err != nil {
		log.Fatal("Failed to get current schema:", err)
	}
	log.Printf("Current schema: %s", currentSchema)
	log.Printf("/******************************************/")

	// Auto-migrate with better error handling
	log.Printf("Starting database migration...")

	// Check if tables exist before migration
	if db.Migrator().HasTable(&domain.User{}) {
		log.Printf("User table already exists")
	} else {
		log.Printf("User table does not exist, will be created")
	}

	if db.Migrator().HasTable(&domain.Employee{}) {
		log.Printf("Employee table already exists")
	} else {
		log.Printf("Employee table does not exist, will be created")
	}

	// Migrate each table individually for better error tracking
	log.Printf("Migrating User table...")
	log.Printf("User struct: %+v", domain.User{})

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Printf("Failed to migrate User table: %v", err)
		log.Fatal("Failed to migrate User table:", err)
	}

	// Verificar se foi criada
	if db.Migrator().HasTable(&domain.User{}) {
		log.Printf("✅ User table created successfully")

		// Verificar colunas
		columns, _ := db.Migrator().ColumnTypes(&domain.User{})
		log.Printf("User table columns:")
		for _, col := range columns {
			log.Printf("  - %s (%s)", col.Name(), col.DatabaseTypeName())
		}
	} else {
		log.Fatal("❌ User table was not created")
	}

	log.Printf("User table migration completed")

	log.Printf("Migrating Employee table...")
	log.Printf("Employee struct: %+v", domain.Employee{})

	if err := db.AutoMigrate(&domain.Employee{}); err != nil {
		log.Printf("Failed to migrate Employee table: %v", err)
		log.Fatal("Failed to migrate Employee table:", err)
	}

	// Verificar se foi criada
	if db.Migrator().HasTable(&domain.Employee{}) {
		log.Printf("✅ Employee table created successfully")

		// Verificar colunas
		columns, _ := db.Migrator().ColumnTypes(&domain.Employee{})
		log.Printf("Employee table columns:")
		for _, col := range columns {
			log.Printf("  - %s (%s)", col.Name(), col.DatabaseTypeName())
		}
	} else {
		log.Fatal("❌ Employee table was not created")
	}
	log.Printf("Employee table migration completed")

	log.Printf("Database migration completed successfully")
	log.Printf("/******************************************/")

	// Redis connection with timeout
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	log.Printf("Connecting to Redis: %s", redisAddr)
	log.Printf("/******************************************/")

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test Redis connection with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Printf("Continuing without Redis cache...")
		rdb = nil
	} else {
		log.Printf("Redis connection established")
	}
	log.Printf("/******************************************/")

	// Initialize repositories
	userRepo := postgresRepo.NewPostgresUserRepository(db)
	employeeRepo := postgresRepo.NewPostgresEmployeeRepository(db)

	var cacheRepo ports.CacheRepository
	if err == nil {
		log.Printf("Redis connection established")
		cacheRepo = redisRepo.NewRedisCacheRepository(rdb)
	} else {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Printf("Continuing without Redis cache...")
		cacheRepo = &NullCacheRepository{}
	}

	log.Printf("Repositories initialized")
	log.Printf("/******************************************/")

	// Initialize use cases
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key"
		log.Printf("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	log.Printf("JWT secret initialized")
	log.Printf("/******************************************/")

	authUseCase := auth.NewAuthUseCase(userRepo, jwtSecret, 24)
	employeeUseCase := employee.NewEmployeeUseCase(employeeRepo, cacheRepo)
	log.Printf("Use cases initialized")
	log.Printf("/******************************************/")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	employeeHandler := handlers.NewEmployeeHandler(employeeUseCase)
	log.Printf("Handlers initialized")
	log.Printf("/******************************************/")

	// Setup router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	log.Printf("Router initialized")
	log.Printf("/******************************************/")

	// Setup routes
	routes.SetupAllRoutes(router, authHandler, employeeHandler)
	log.Printf("Routes set up")
	log.Printf("/******************************************/")

	log.Printf("Health check endpoint set up")
	log.Printf("/******************************************/")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server port set to %s", port)
	log.Printf("/******************************************/")

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	log.Printf("Server started successfully")
	log.Printf("/*************** End Main ***************/")
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
