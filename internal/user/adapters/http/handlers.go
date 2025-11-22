package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-scaffolding/internal/user/domain"
	"github.com/yourusername/go-scaffolding/internal/user/ports"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService ports.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Email, req.Name)
	if err != nil {
		statusCode, errorMsg := mapDomainErrorToHTTP(err)
		c.JSON(statusCode, ErrorResponse{
			Error: errorMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, ToUserResponse(user))
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		statusCode, errorMsg := mapDomainErrorToHTTP(err)
		c.JSON(statusCode, ErrorResponse{
			Error: errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, ToUserResponse(user))
}

// GetUserByEmail handles GET /users/email/:email
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := h.userService.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		statusCode, errorMsg := mapDomainErrorToHTTP(err)
		c.JSON(statusCode, ErrorResponse{
			Error: errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, ToUserResponse(user))
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req.Name)
	if err != nil {
		statusCode, errorMsg := mapDomainErrorToHTTP(err)
		c.JSON(statusCode, ErrorResponse{
			Error: errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, ToUserResponse(user))
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		statusCode, errorMsg := mapDomainErrorToHTTP(err)
		c.JSON(statusCode, ErrorResponse{
			Error: errorMsg,
		})
		return
	}

	c.Status(http.StatusNoContent)
}

const (
	// MaxLimit defines the maximum number of users that can be fetched in a single request
	MaxLimit = 100
)

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters with defaults
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Enforce maximum limit to prevent database overload
	if limit > MaxLimit {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "limit cannot exceed 100",
		})
		return
	}

	users, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		statusCode, errorMsg := mapDomainErrorToHTTP(err)
		c.JSON(statusCode, ErrorResponse{
			Error: errorMsg,
		})
		return
	}

	response := ListUsersResponse{
		Users:  ToUsersResponse(users),
		Limit:  limit,
		Offset: offset,
	}

	c.JSON(http.StatusOK, response)
}

// mapDomainErrorToHTTP maps domain errors to HTTP status codes and messages
func mapDomainErrorToHTTP(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrInvalidEmail):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, domain.ErrInvalidName):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, domain.ErrDuplicateEmail):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
