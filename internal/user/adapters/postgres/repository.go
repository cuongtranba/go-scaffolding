package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/yourusername/go-scaffolding/internal/user/domain"
	"github.com/yourusername/go-scaffolding/internal/user/ports"
)

// userRepository implements ports.UserRepository using GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	model := ToUserModel(user)

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		// Check for unique constraint violation
		if isDuplicateEmailError(result.Error) {
			return domain.ErrDuplicateEmail
		}
		return result.Error
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var model UserModel

	result := r.db.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}

	return ToDomainUser(&model), nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel

	result := r.db.WithContext(ctx).Where("email = ?", email).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}

	return ToDomainUser(&model), nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	model := ToUserModel(user)

	result := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", user.ID).
		Updates(model)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserModel{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List retrieves users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var models []*UserModel

	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	return ToDomainUsers(models), nil
}

// isDuplicateEmailError checks if the error is a unique constraint violation on email
func isDuplicateEmailError(err error) bool {
	if err == nil {
		return false
	}

	// GORM wraps the error, so we check the error message
	errMsg := err.Error()

	// PostgreSQL unique constraint violation
	if contains(errMsg, "duplicate key value violates unique constraint") ||
		contains(errMsg, "UNIQUE constraint failed") {
		return contains(errMsg, "email")
	}

	return false
}

// contains checks if a string contains a substring (case-sensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
