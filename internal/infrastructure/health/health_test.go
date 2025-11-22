package health

import (
	"context"
	"errors"
	"testing"

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
