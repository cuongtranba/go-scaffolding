package database

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/go-scaffolding/internal/config"
	"github.com/yourusername/go-scaffolding/internal/infrastructure/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewPostgresDB creates a new PostgreSQL database connection using GORM
func NewPostgresDB(cfg *config.Config, log *logger.Logger) (*gorm.DB, error) {
	dsn := cfg.Postgres.ConnectionString()

	// Configure GORM logger to use our zerolog logger
	gormLogger := gormlogger.New(
		&gormLoggerAdapter{logger: log},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  mapLogLevel(cfg.Postgres.LogLevel),
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime)

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// ClosePostgresDB closes the PostgreSQL database connection
func ClosePostgresDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// gormLoggerAdapter adapts our zerolog logger to GORM's logger interface
type gormLoggerAdapter struct {
	logger *logger.Logger
}

func (l *gormLoggerAdapter) Printf(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// mapLogLevel maps config log level to GORM log level
func mapLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "debug":
		return gormlogger.Info
	case "info":
		return gormlogger.Warn
	case "warn":
		return gormlogger.Error
	case "error":
		return gormlogger.Error
	default:
		return gormlogger.Warn
	}
}
