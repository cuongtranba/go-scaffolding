package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/yourusername/go-scaffolding/internal/user/domain"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Use SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the schema
	err = db.AutoMigrate(&UserModel{})
	require.NoError(t, err)

	return db
}

func TestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successfully creates user", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New().String(),
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		assert.NoError(t, err)

		// Verify user was created
		var model UserModel
		err = db.Where("id = ?", user.ID).First(&model).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Email, model.Email)
		assert.Equal(t, user.Name, model.Name)
	})

	t.Run("returns error on duplicate email", func(t *testing.T) {
		user1 := &domain.User{
			ID:        uuid.New().String(),
			Email:     "duplicate@example.com",
			Name:      "User One",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := &domain.User{
			ID:        uuid.New().String(),
			Email:     "duplicate@example.com",
			Name:      "User Two",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = repo.Create(ctx, user2)
		assert.ErrorIs(t, err, domain.ErrDuplicateEmail)
	})
}

func TestRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successfully retrieves user by ID", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New().String(),
			Email:     "getbyid@example.com",
			Name:      "Get By ID User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Email, retrieved.Email)
		assert.Equal(t, user.Name, retrieved.Name)
	})

	t.Run("returns ErrUserNotFound for non-existent ID", func(t *testing.T) {
		nonExistentID := uuid.New().String()

		retrieved, err := repo.GetByID(ctx, nonExistentID)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, retrieved)
	})
}

func TestRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successfully retrieves user by email", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New().String(),
			Email:     "getbyemail@example.com",
			Name:      "Get By Email User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		retrieved, err := repo.GetByEmail(ctx, user.Email)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Email, retrieved.Email)
		assert.Equal(t, user.Name, retrieved.Name)
	})

	t.Run("returns ErrUserNotFound for non-existent email", func(t *testing.T) {
		retrieved, err := repo.GetByEmail(ctx, "nonexistent@example.com")
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, retrieved)
	})
}

func TestRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successfully updates user", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New().String(),
			Email:     "update@example.com",
			Name:      "Original Name",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Update the user
		user.Name = "Updated Name"
		user.UpdatedAt = time.Now()

		err = repo.Update(ctx, user)
		assert.NoError(t, err)

		// Verify update
		retrieved, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", retrieved.Name)
	})

	t.Run("returns ErrUserNotFound for non-existent user", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New().String(),
			Email:     "nonexistent@example.com",
			Name:      "Non Existent",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Update(ctx, user)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}

func TestRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successfully deletes user", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New().String(),
			Email:     "delete@example.com",
			Name:      "Delete User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		err = repo.Delete(ctx, user.ID)
		assert.NoError(t, err)

		// Verify user is deleted
		retrieved, err := repo.GetByID(ctx, user.ID)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("returns ErrUserNotFound for non-existent user", func(t *testing.T) {
		nonExistentID := uuid.New().String()

		err := repo.Delete(ctx, nonExistentID)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}

func TestRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test users
	users := []*domain.User{
		{
			ID:        uuid.New().String(),
			Email:     "user1@example.com",
			Name:      "User One",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Email:     "user2@example.com",
			Name:      "User Two",
			CreatedAt: time.Now().Add(1 * time.Second),
			UpdatedAt: time.Now().Add(1 * time.Second),
		},
		{
			ID:        uuid.New().String(),
			Email:     "user3@example.com",
			Name:      "User Three",
			CreatedAt: time.Now().Add(2 * time.Second),
			UpdatedAt: time.Now().Add(2 * time.Second),
		},
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	t.Run("retrieves all users with no pagination", func(t *testing.T) {
		retrieved, err := repo.List(ctx, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, retrieved, 3)
	})

	t.Run("retrieves users with limit", func(t *testing.T) {
		retrieved, err := repo.List(ctx, 2, 0)
		assert.NoError(t, err)
		assert.Len(t, retrieved, 2)
	})

	t.Run("retrieves users with offset", func(t *testing.T) {
		retrieved, err := repo.List(ctx, 10, 2)
		assert.NoError(t, err)
		assert.Len(t, retrieved, 1)
	})

	t.Run("returns empty slice when no users found", func(t *testing.T) {
		retrieved, err := repo.List(ctx, 10, 100)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Len(t, retrieved, 0)
	})
}
