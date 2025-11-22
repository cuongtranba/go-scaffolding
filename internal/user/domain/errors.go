package domain

import "errors"

var (
	// ErrUserNotFound indicates user was not found
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidEmail indicates email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")

	// ErrInvalidName indicates name is invalid
	ErrInvalidName = errors.New("name must be non-empty and not exceed 255 characters")

	// ErrDuplicateEmail indicates email already exists
	ErrDuplicateEmail = errors.New("email already exists")
)
