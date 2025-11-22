package http

import (
	"time"

	"github.com/yourusername/go-scaffolding/internal/user/domain"
)

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

// UpdateUserRequest represents the request to update a user
// Note: Email updates are not allowed to maintain data integrity and prevent
// conflicts with existing email constraints in the database
type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Users  []UserResponse `json:"users"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ToUserResponse converts a domain user to a user response
func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUsersResponse converts a slice of domain users to user responses
func ToUsersResponse(users []*domain.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, ToUserResponse(user))
	}
	return responses
}
