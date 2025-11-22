package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/go-scaffolding/internal/user/domain"
	"github.com/yourusername/go-scaffolding/internal/user/ports/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	email := "test@example.com"
	name := "Test User"

	mockRepo.On("GetByEmail", ctx, email).Return(nil, domain.ErrUserNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.CreateUser(ctx, email, name)
	require.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, name, user.Name)

	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	existingUser := &domain.User{Email: "test@example.com"}

	mockRepo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)

	_, err := service.CreateUser(ctx, "test@example.com", "Test User")
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrDuplicateEmail)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	expectedUser := &domain.User{ID: "123", Email: "test@example.com", Name: "Test User"}

	mockRepo.On("GetByID", ctx, "123").Return(expectedUser, nil)

	user, err := service.GetUser(ctx, "123")
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	existingUser, _ := domain.NewUser("test@example.com", "Old Name")

	mockRepo.On("GetByID", ctx, existingUser.ID).Return(existingUser, nil)
	mockRepo.On("Update", ctx, existingUser).Return(nil)

	user, err := service.UpdateUser(ctx, existingUser.ID, "New Name")
	require.NoError(t, err)
	assert.Equal(t, "New Name", user.Name)

	mockRepo.AssertExpectations(t)
}
