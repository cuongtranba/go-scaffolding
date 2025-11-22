package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	userhttp "github.com/yourusername/go-scaffolding/internal/user/adapters/http"
	userPostgres "github.com/yourusername/go-scaffolding/internal/user/adapters/postgres"
	"github.com/yourusername/go-scaffolding/internal/user/ports"
	userservice "github.com/yourusername/go-scaffolding/internal/user/service"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get connection string")

	// Connect to database
	db, err := gorm.Open(gormpostgres.Open(connStr), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to database")

	// Run migrations
	err = db.AutoMigrate(&userPostgres.UserModel{})
	require.NoError(t, err, "Failed to run migrations")

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

func TestIntegration_UserAPI(t *testing.T) {
	// Setup test database
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create service and handler
	repo := userPostgres.NewUserRepository(db)
	usersvc := userservice.NewUserService(repo)
	router := setupTestRouter(usersvc)

	// Test data
	userEmail := "integration@example.com"
	userName := "Integration Test User"
	updatedName := "Updated Name"

	var userID string

	t.Run("CreateUser", func(t *testing.T) {
		reqBody := map[string]string{
			"email": userEmail,
			"name":  userName,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, userEmail, resp["email"])
		assert.Equal(t, userName, resp["name"])
		assert.NotEmpty(t, resp["id"])

		userID = resp["id"].(string)
	})

	t.Run("GetUserByID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/"+userID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, userID, resp["id"])
		assert.Equal(t, userEmail, resp["email"])
		assert.Equal(t, userName, resp["name"])
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/email/"+userEmail, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, userID, resp["id"])
		assert.Equal(t, userEmail, resp["email"])
	})

	t.Run("UpdateUser", func(t *testing.T) {
		reqBody := map[string]string{
			"name": updatedName,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/users/"+userID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, userID, resp["id"])
		assert.Equal(t, updatedName, resp["name"])
	})

	t.Run("ListUsers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users?limit=10&offset=0", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		users := resp["users"].([]interface{})
		assert.Len(t, users, 1)

		firstUser := users[0].(map[string]interface{})
		assert.Equal(t, userID, firstUser["id"])
		assert.Equal(t, updatedName, firstUser["name"])
	})

	t.Run("DeleteUser", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/"+userID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify user is deleted
		req = httptest.NewRequest(http.MethodGet, "/users/"+userID, nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		// Create first user
		reqBody := map[string]string{
			"email": "duplicate@example.com",
			"name":  "First User",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Try to create second user with same email
		reqBody["name"] = "Second User"
		body, _ = json.Marshal(reqBody)

		req = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "email already exists", resp["error"])
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		reqBody := map[string]string{
			"email": "invalid-email",
			"name":  "Test User",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("MaxLimitExceeded", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users?limit=101", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "limit cannot exceed 100", resp["error"])
	})
}

func TestIntegration_DatabasePersistence(t *testing.T) {
	// Setup test database
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create service and handler
	repo := userPostgres.NewUserRepository(db)
	usersvc := userservice.NewUserService(repo)
	router := setupTestRouter(usersvc)

	t.Run("DataPersistsAcrossRequests", func(t *testing.T) {
		// Create user
		reqBody := map[string]string{
			"email": "persist@example.com",
			"name":  "Persist User",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createResp)
		userID := createResp["id"].(string)

		// Verify data persists by querying database directly
		var count int64
		err := db.Model(&userPostgres.UserModel{}).Where("id = ?", userID).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// Verify data can be retrieved via API
		req = httptest.NewRequest(http.MethodGet, "/users/"+userID, nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var getResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &getResp)
		assert.Equal(t, createResp["email"], getResp["email"])
		assert.Equal(t, createResp["name"], getResp["name"])
	})

	t.Run("UpdatesArePersistedCorrectly", func(t *testing.T) {
		// Create user
		reqBody := map[string]string{
			"email": "update@example.com",
			"name":  "Original Name",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var createResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createResp)
		userID := createResp["id"].(string)
		originalUpdatedAt := createResp["updated_at"].(string)

		// Wait a moment to ensure timestamp difference
		time.Sleep(10 * time.Millisecond)

		// Update user
		reqBody = map[string]string{
			"name": "Updated Name",
		}
		body, _ = json.Marshal(reqBody)

		req = httptest.NewRequest(http.MethodPut, "/users/"+userID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var updateResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &updateResp)

		// Verify updated_at timestamp changed
		assert.NotEqual(t, originalUpdatedAt, updateResp["updated_at"])

		// Verify update persisted in database
		var model userPostgres.UserModel
		err := db.Where("id = ?", userID).First(&model).Error
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", model.Name)
	})
}

// setupTestRouter creates a Gin router with user routes for testing
func setupTestRouter(userService ports.UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	userhttp.RegisterUserRoutes(router, userService)
	return router
}
