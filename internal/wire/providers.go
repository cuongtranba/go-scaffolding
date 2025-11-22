package wire

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/yourusername/go-scaffolding/internal/config"
	"github.com/yourusername/go-scaffolding/internal/infrastructure/database"
	"github.com/yourusername/go-scaffolding/internal/infrastructure/health"
	"github.com/yourusername/go-scaffolding/internal/infrastructure/logger"
	"github.com/yourusername/go-scaffolding/internal/user/adapters/http"
	"github.com/yourusername/go-scaffolding/internal/user/adapters/postgres"
	"github.com/yourusername/go-scaffolding/internal/user/ports"
	"github.com/yourusername/go-scaffolding/internal/user/service"
	"gorm.io/gorm"
)

// ProviderSet is the Wire provider set that includes all dependencies
var ProviderSet = wire.NewSet(
	// Infrastructure
	ProvideConfig,
	ProvideLogger,
	ProvideHealthChecker,
	ProvidePostgresDB,

	// User domain
	ProvideUserRepository,
	ProvideUserService,

	// HTTP server
	ProvideGinEngine,
)

// ProvideConfig provides the application configuration
func ProvideConfig(configPath string) (*config.Config, error) {
	return config.Load(configPath)
}

// ProvideLogger provides the logger instance
func ProvideLogger(cfg *config.Config) *logger.Logger {
	return logger.New(cfg.App.LogLevel, os.Stdout)
}

// ProvideHealthChecker provides the health checker instance
func ProvideHealthChecker() *health.Checker {
	return health.NewChecker()
}

// ProvidePostgresDB provides the PostgreSQL database connection
func ProvidePostgresDB(cfg *config.Config, log *logger.Logger) (*gorm.DB, func(), error) {
	db, err := database.NewPostgresDB(cfg, log)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := database.ClosePostgresDB(db); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}

	return db, cleanup, nil
}

// ProvideUserRepository provides the user repository implementation
func ProvideUserRepository(db *gorm.DB) ports.UserRepository {
	return postgres.NewUserRepository(db)
}

// ProvideUserService provides the user service implementation
func ProvideUserService(repo ports.UserRepository) ports.UserService {
	return service.NewUserService(repo)
}

// ProvideGinEngine provides the configured Gin engine with all routes
func ProvideGinEngine(cfg *config.Config, userService ports.UserService, healthChecker *health.Checker) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health check routes
	router.GET("/health/live", func(c *gin.Context) {
		result := healthChecker.Liveness()
		status := 200
		if result.Status != "healthy" {
			status = 503
		}
		c.JSON(status, result)
	})

	router.GET("/health/ready", func(c *gin.Context) {
		result := healthChecker.Readiness(c.Request.Context())
		status := 200
		if result.Status != "healthy" {
			status = 503
		}
		c.JSON(status, result)
	})

	// Register user routes
	http.RegisterUserRoutes(router, userService)

	return router
}
