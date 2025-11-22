package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("test@example.com", "Test User")
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestNewUser_InvalidEmail(t *testing.T) {
	_, err := NewUser("invalid-email", "Test User")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidEmail)
}

func TestNewUser_EmptyName(t *testing.T) {
	_, err := NewUser("test@example.com", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidName)
}

func TestUser_UpdateName(t *testing.T) {
	user, err := NewUser("test@example.com", "Old Name")
	require.NoError(t, err)

	err = user.UpdateName("New Name")
	require.NoError(t, err)
	assert.Equal(t, "New Name", user.Name)
}

func TestUser_UpdateName_Empty(t *testing.T) {
	user, err := NewUser("test@example.com", "Test User")
	require.NoError(t, err)

	err = user.UpdateName("")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidName)
}
