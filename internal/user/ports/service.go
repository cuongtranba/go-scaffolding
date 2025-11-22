package ports

import (
	"context"

	"github.com/yourusername/go-scaffolding/internal/user/domain"
)

//go:generate mockery --name=UserService --output=mocks --outpkg=mocks

// UserService defines the interface for user business logic
type UserService interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, email, name string) (*domain.User, error)

	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, id string) (*domain.User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// UpdateUser updates a user's information
	UpdateUser(ctx context.Context, id, name string) (*domain.User, error)

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id string) error

	// ListUsers retrieves users with pagination
	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
}
