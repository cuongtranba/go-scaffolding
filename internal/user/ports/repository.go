package ports

import (
	"context"

	"github.com/yourusername/go-scaffolding/internal/user/domain"
)

//go:generate mockery --name=UserRepository --output=mocks --outpkg=mocks

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*domain.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error

	// List retrieves users with pagination
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}
