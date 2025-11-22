package domain

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Email validation regex that prevents:
// - consecutive dots (..)
// - leading/trailing dots in local part
// - leading/trailing dots in domain
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

const (
	maxEmailLength = 254
	maxNameLength  = 255
)

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

	// Trim whitespace and validate name
	name = strings.TrimSpace(name)
	if err := isValidName(name); err != nil {
		return nil, err
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
	// Trim whitespace and validate name
	name = strings.TrimSpace(name)
	if err := isValidName(name); err != nil {
		return err
	}

	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

func isValidName(name string) error {
	if name == "" {
		return ErrInvalidName
	}

	if len(name) > maxNameLength {
		return ErrInvalidName
	}

	return nil
}

func isValidEmail(email string) bool {
	// Check length
	if len(email) > maxEmailLength {
		return false
	}

	// Check basic format
	if !emailRegex.MatchString(email) {
		return false
	}

	// Split into local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	local := parts[0]
	domain := parts[1]

	// Check for consecutive dots
	if strings.Contains(email, "..") {
		return false
	}

	// Check for leading/trailing dots in local part
	if strings.HasPrefix(local, ".") || strings.HasSuffix(local, ".") {
		return false
	}

	// Check for leading/trailing dots in domain
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	return true
}
