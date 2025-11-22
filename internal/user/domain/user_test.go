package domain

import (
	"testing"
	"time"

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

func TestNewUser_EmailEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		shouldErr bool
	}{
		// Valid emails
		{name: "valid simple", email: "test@example.com", shouldErr: false},
		{name: "valid with plus", email: "test+tag@example.com", shouldErr: false},
		{name: "valid with hyphen", email: "test-user@example.com", shouldErr: false},
		{name: "valid with underscore", email: "test_user@example.com", shouldErr: false},
		{name: "valid with dot", email: "test.user@example.com", shouldErr: false},
		{name: "valid subdomain", email: "test@mail.example.com", shouldErr: false},

		// Invalid emails - consecutive dots
		{name: "consecutive dots in local", email: "test..user@example.com", shouldErr: true},
		{name: "consecutive dots in domain", email: "test@example..com", shouldErr: true},

		// Invalid emails - leading/trailing dots
		{name: "leading dot in local", email: ".test@example.com", shouldErr: true},
		{name: "trailing dot in local", email: "test.@example.com", shouldErr: true},
		{name: "leading dot in domain", email: "test@.example.com", shouldErr: true},
		{name: "trailing dot in domain", email: "test@example.com.", shouldErr: true},

		// Invalid emails - format issues
		{name: "no at sign", email: "testexample.com", shouldErr: true},
		{name: "multiple at signs", email: "test@@example.com", shouldErr: true},
		{name: "no domain", email: "test@", shouldErr: true},
		{name: "no local part", email: "@example.com", shouldErr: true},
		{name: "no TLD", email: "test@example", shouldErr: true},

		// Invalid emails - length
		{name: "too long", email: "a" + string(make([]byte, 250)) + "@example.com", shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, "Test User")
			if tt.shouldErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidEmail)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
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

func TestUser_UpdateName_UpdatesTimestamp(t *testing.T) {
	user, err := NewUser("test@example.com", "Old Name")
	require.NoError(t, err)

	originalCreatedAt := user.CreatedAt
	originalUpdatedAt := user.UpdatedAt

	// Sleep briefly to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	err = user.UpdateName("New Name")
	require.NoError(t, err)

	// CreatedAt should not change
	assert.Equal(t, originalCreatedAt, user.CreatedAt)

	// UpdatedAt should change
	assert.NotEqual(t, originalUpdatedAt, user.UpdatedAt)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestNewUser_NameWhitespace(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
		expected  string
		shouldErr bool
	}{
		{name: "valid name", inputName: "John Doe", expected: "John Doe", shouldErr: false},
		{name: "leading whitespace", inputName: "  John Doe", expected: "John Doe", shouldErr: false},
		{name: "trailing whitespace", inputName: "John Doe  ", expected: "John Doe", shouldErr: false},
		{name: "both whitespace", inputName: "  John Doe  ", expected: "John Doe", shouldErr: false},
		{name: "tabs and spaces", inputName: "\t John Doe \t", expected: "John Doe", shouldErr: false},
		{name: "only whitespace", inputName: "   ", expected: "", shouldErr: true},
		{name: "only tabs", inputName: "\t\t", expected: "", shouldErr: true},
		{name: "newlines", inputName: "\n\n", expected: "", shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser("test@example.com", tt.inputName)
			if tt.shouldErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidName)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, user.Name)
			}
		})
	}
}

func TestNewUser_NameLength(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
		shouldErr bool
	}{
		{name: "valid short name", inputName: "Jo", shouldErr: false},
		{name: "valid max length", inputName: string(make([]byte, 255)), shouldErr: false},
		{name: "too long", inputName: string(make([]byte, 256)), shouldErr: true},
		{name: "way too long", inputName: string(make([]byte, 1000)), shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill with 'a' for valid names
			name := tt.inputName
			if !tt.shouldErr {
				name = ""
				for i := 0; i < len(tt.inputName); i++ {
					name += "a"
				}
			} else {
				name = ""
				for i := 0; i < len(tt.inputName); i++ {
					name += "a"
				}
			}

			user, err := NewUser("test@example.com", name)
			if tt.shouldErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidName)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(name), len(user.Name))
			}
		})
	}
}

func TestUser_UpdateName_Whitespace(t *testing.T) {
	user, err := NewUser("test@example.com", "Old Name")
	require.NoError(t, err)

	tests := []struct {
		name      string
		inputName string
		expected  string
		shouldErr bool
	}{
		{name: "leading whitespace", inputName: "  New Name", expected: "New Name", shouldErr: false},
		{name: "trailing whitespace", inputName: "New Name  ", expected: "New Name", shouldErr: false},
		{name: "only whitespace", inputName: "   ", expected: "", shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.UpdateName(tt.inputName)
			if tt.shouldErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidName)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, user.Name)
			}
		})
	}
}

func TestUser_UpdateName_Length(t *testing.T) {
	user, err := NewUser("test@example.com", "Old Name")
	require.NoError(t, err)

	// Test too long name
	tooLongName := ""
	for i := 0; i < 256; i++ {
		tooLongName += "a"
	}

	err = user.UpdateName(tooLongName)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidName)

	// Test max length name (should succeed)
	maxLengthName := ""
	for i := 0; i < 255; i++ {
		maxLengthName += "a"
	}

	err = user.UpdateName(maxLengthName)
	require.NoError(t, err)
	assert.Equal(t, maxLengthName, user.Name)
}
