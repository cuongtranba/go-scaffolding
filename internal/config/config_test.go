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
