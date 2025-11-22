package health

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChecker_AddCheck(t *testing.T) {
	checker := NewChecker()
	checker.AddCheck("test", func(ctx context.Context) error {
		return nil
	})

	result := checker.Check(context.Background())
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Contains(t, result.Checks, "test")
	assert.Equal(t, StatusHealthy, result.Checks["test"].Status)
}

func TestChecker_UnhealthyCheck(t *testing.T) {
	checker := NewChecker()
	checker.AddCheck("failing", func(ctx context.Context) error {
		return errors.New("service unavailable")
	})

	result := checker.Check(context.Background())
	assert.Equal(t, StatusUnhealthy, result.Status)
	assert.Equal(t, StatusUnhealthy, result.Checks["failing"].Status)
	assert.Equal(t, "service unavailable", result.Checks["failing"].Error)
}

func TestChecker_MultipleChecks(t *testing.T) {
	checker := NewChecker()
	checker.AddCheck("pass", func(ctx context.Context) error {
		return nil
	})
	checker.AddCheck("fail", func(ctx context.Context) error {
		return errors.New("failed")
	})

	result := checker.Check(context.Background())
	assert.Equal(t, StatusUnhealthy, result.Status)
	assert.Equal(t, StatusHealthy, result.Checks["pass"].Status)
	assert.Equal(t, StatusUnhealthy, result.Checks["fail"].Status)
}

func TestChecker_Liveness(t *testing.T) {
	checker := NewChecker()

	// Liveness should always return healthy, even if checks are failing
	checker.AddCheck("failing", func(ctx context.Context) error {
		return errors.New("service unavailable")
	})

	result := checker.Liveness()
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Contains(t, result.Checks, "liveness")
	assert.Equal(t, StatusHealthy, result.Checks["liveness"].Status)
	assert.Empty(t, result.Checks["liveness"].Error)
}

func TestChecker_Readiness(t *testing.T) {
	checker := NewChecker()

	// Add multiple checks
	checker.AddCheck("database", func(ctx context.Context) error {
		return nil
	})
	checker.AddCheck("cache", func(ctx context.Context) error {
		return nil
	})

	result := checker.Readiness(context.Background())

	// Readiness should run all checks
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Contains(t, result.Checks, "database")
	assert.Contains(t, result.Checks, "cache")
	assert.Equal(t, StatusHealthy, result.Checks["database"].Status)
	assert.Equal(t, StatusHealthy, result.Checks["cache"].Status)
}

func TestChecker_CheckTimeout(t *testing.T) {
	checker := NewChecker()

	// Add a slow check that takes >5 seconds
	checker.AddCheck("slow", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
			return nil
		}
	})

	start := time.Now()
	result := checker.Check(context.Background())
	duration := time.Since(start)

	// Check should complete in less than 6 seconds (5s timeout + buffer)
	assert.Less(t, duration, 6*time.Second)

	// The slow check should have failed due to timeout
	assert.Equal(t, StatusUnhealthy, result.Status)
	assert.Contains(t, result.Checks, "slow")
	assert.Equal(t, StatusUnhealthy, result.Checks["slow"].Status)
	assert.Contains(t, result.Checks["slow"].Error, "context deadline exceeded")
}
