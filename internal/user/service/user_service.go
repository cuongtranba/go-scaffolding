package service

import (
	"context"
	"errors"

	"github.com/yourusername/go-scaffolding/internal/user/domain"
	"github.com/yourusername/go-scaffolding/internal/user/ports"
)

// UserService implements the UserService port
type UserService struct {
	repo ports.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo ports.UserRepository) ports.UserService {
	return &UserService{
		repo: repo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	// Check if email already exists
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrDuplicateEmail
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	// Create new user
	user, err := domain.NewUser(email, name)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(ctx context.Context, id, name string) (*domain.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update name
	if err := user.UpdateName(name); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// ListUsers retrieves users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return s.repo.List(ctx, limit, offset)
}
