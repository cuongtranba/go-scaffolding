package domain

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// User represents a user entity
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user with validation
func NewUser(email, name string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	if name == "" {
		return nil, ErrInvalidName
	}

	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateName updates the user's name
func (u *User) UpdateName(name string) error {
	if name == "" {
		return ErrInvalidName
	}

	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
