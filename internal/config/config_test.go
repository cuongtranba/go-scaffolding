package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_FromFile(t *testing.T) {
	// Create temp config file
	configContent := `
app:
  name: test-app
  environment: development
  http_port: 8080
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	require.NoError(t, err)
	assert.Equal(t, "test-app", cfg.App.Name)
	assert.Equal(t, "development", cfg.App.Environment)
	assert.Equal(t, 8080, cfg.App.HTTPPort)
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("APP_NAME", "env-app")
	os.Setenv("APP_HTTP_PORT", "9000")
	defer os.Unsetenv("APP_NAME")
	defer os.Unsetenv("APP_HTTP_PORT")

	// Create minimal config file
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
app:
  environment: development
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	require.NoError(t, err)
	assert.Equal(t, "env-app", cfg.App.Name)
	assert.Equal(t, 9000, cfg.App.HTTPPort)
}

func TestPostgresConfig_ConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		config   PostgresConfig
		expected string
	}{
		{
			name: "full connection string",
			config: PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Password: "testpass",
				Database: "testdb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable",
		},
		{
			name: "connection string with SSL enabled",
			config: PostgresConfig{
				Host:     "db.example.com",
				Port:     5433,
				User:     "admin",
				Password: "secret123",
				Database: "production",
				SSLMode:  "require",
			},
			expected: "host=db.example.com port=5433 user=admin password=secret123 dbname=production sslmode=require",
		},
		{
			name: "connection string with empty password",
			config: PostgresConfig{
				Host:     "127.0.0.1",
				Port:     5432,
				User:     "postgres",
				Password: "",
				Database: "mydb",
				SSLMode:  "disable",
			},
			expected: "host=127.0.0.1 port=5432 user=postgres password= dbname=mydb sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ConnectionString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedisConfig_Address(t *testing.T) {
	tests := []struct {
		name     string
		config   RedisConfig
		expected string
	}{
		{
			name: "localhost with default port",
			config: RedisConfig{
				Host: "localhost",
				Port: 6379,
			},
			expected: "localhost:6379",
		},
		{
			name: "custom host and port",
			config: RedisConfig{
				Host: "redis.example.com",
				Port: 6380,
			},
			expected: "redis.example.com:6380",
		},
		{
			name: "IP address with custom port",
			config: RedisConfig{
				Host: "192.168.1.100",
				Port: 7000,
			},
			expected: "192.168.1.100:7000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.Address()
			assert.Equal(t, tt.expected, result)
		})
	}
}
