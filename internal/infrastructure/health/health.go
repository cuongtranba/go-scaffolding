package health

import (
	"context"
	"sync"
	"time"
)

// Status represents health status
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

// CheckFunc is a function that performs a health check
type CheckFunc func(ctx context.Context) error

// CheckResult represents the result of a single health check
type CheckResult struct {
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}

// HealthResult represents overall health status
type HealthResult struct {
	Status Status                 `json:"status"`
	Checks map[string]CheckResult `json:"checks"`
}

// Checker manages health checks
type Checker struct {
	checks map[string]CheckFunc
	mu     sync.RWMutex
}

// NewChecker creates a new health checker
func NewChecker() *Checker {
	return &Checker{
		checks: make(map[string]CheckFunc),
	}
}

// AddCheck registers a health check with a name
func (c *Checker) AddCheck(name string, check CheckFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks[name] = check
}

// Check runs all health checks and returns the result
func (c *Checker) Check(ctx context.Context) HealthResult {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Set timeout for all checks
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := HealthResult{
		Status: StatusHealthy,
		Checks: make(map[string]CheckResult),
	}

	// Run all checks
	for name, check := range c.checks {
		checkResult := CheckResult{Status: StatusHealthy}

		if err := check(ctx); err != nil {
			checkResult.Status = StatusUnhealthy
			checkResult.Error = err.Error()
			result.Status = StatusUnhealthy
		}

		result.Checks[name] = checkResult
	}

	return result
}

// Liveness always returns healthy (app is running)
func (c *Checker) Liveness() HealthResult {
	return HealthResult{
		Status: StatusHealthy,
		Checks: map[string]CheckResult{
			"liveness": {Status: StatusHealthy},
		},
	}
}

// Readiness checks if app is ready to serve traffic
func (c *Checker) Readiness(ctx context.Context) HealthResult {
	return c.Check(ctx)
}
