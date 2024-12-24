// internal/config/config_test.go
package config_test

import (
	"coupon_service/internal/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setEnv sets environment variables and returns a cleanup function.
func setEnv(t *testing.T, vars map[string]string) func() {
	for key, value := range vars {
		err := os.Setenv(key, value)
		if err != nil {
			t.Fatalf("Failed to set environment variable %s: %v", key, err)
		}
	}
	return func() {
		for key := range vars {
			err := os.Unsetenv(key)
			if err != nil {
				t.Fatalf("Failed to unset environment variable %s: %v", key, err)
			}
		}
	}
}

// createTempEnvFile creates a temporary .env file with the given content.
func createTempEnvFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	err := os.WriteFile(envPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary .env file: %v", err)
	}
	return envPath
}

// internal/config/config_test.go
func TestNewConfig_FromDotEnv(t *testing.T) {
	// Create a temporary .env file
	envContent := "API_PORT=9090\nAPI_ENV=production\n"
	envPath := createTempEnvFile(t, envContent)

	// Change the working directory to where the .env file is located
	originalWD, err := os.Getwd()
	assert.NoError(t, err, "Failed to get current working directory")
	defer func() {
		err := os.Chdir(originalWD)
		assert.NoError(t, err, "Failed to revert working directory")
	}()

	err = os.Chdir(filepath.Dir(envPath))
	assert.NoError(t, err, "Failed to change working directory")

	cfg, err := config.New()
	assert.NoError(t, err, "Expected no error when loading config from .env")

	assert.Equal(t, "9090", cfg.API.Port, "API.Port should be loaded from .env")
	assert.Equal(t, "production", cfg.API.Env, "API.Env should be loaded from .env")
}

func TestNewConfig_OverrideEnv(t *testing.T) {
	// Create a temporary .env file
	envContent := "API_PORT=8080\nAPI_ENV=development\n"
	envPath := createTempEnvFile(t, envContent)

	// Change the working directory to where the .env file is located
	originalWD, err := os.Getwd()
	assert.NoError(t, err, "Failed to get current working directory")
	defer func() {
		err := os.Chdir(originalWD)
		assert.NoError(t, err, "Failed to revert working directory")
	}()

	err = os.Chdir(filepath.Dir(envPath))
	assert.NoError(t, err, "Failed to change working directory")

	// Set environment variables that should override .env
	vars := map[string]string{
		"API_PORT": "6060",
	}
	cleanup := setEnv(t, vars)
	defer cleanup()

	cfg, err := config.New()
	assert.NoError(t, err, "Expected no error when loading config with override")

	assert.Equal(t, "6060", cfg.API.Port, "API.Port should be overridden by system environment")
	assert.Equal(t, "development", cfg.API.Env, "API.Env should be loaded from .env")
}

func TestNewConfig_MissingCriticalPort(t *testing.T) {
	// Ensure API_PORT is unset
	cleanup := setEnv(t, map[string]string{
		"API_ENV": "test",
	})
	defer cleanup()

	cfg, err := config.New()
	assert.Error(t, err, "Expected error due to missing API_PORT")
	assert.Contains(t, err.Error(), "critical environment variable API_PORT is missing", "Error message should indicate missing API_PORT")
	assert.Empty(t, cfg.API.Port, "API.Port should be empty when missing")
}

func TestNewConfig_DefaultValues(t *testing.T) {
	// Set only API_PORT to test default for API_ENV
	vars := map[string]string{
		"API_PORT": "5050",
	}
	cleanup := setEnv(t, vars)
	defer cleanup()

	cfg, err := config.New()
	assert.NoError(t, err, "Expected no error when loading config with default values")

	assert.Equal(t, "5050", cfg.API.Port, "API.Port should be loaded from system environment")
	assert.Equal(t, "dev", cfg.API.Env, "API.Env should use default value 'dev'")
}

func TestNewConfig_NoDotEnv(t *testing.T) {
	// Ensure no .env file is present by not creating one and ensuring environment variables are set
	vars := map[string]string{
		"API_PORT": "4040",
		"API_ENV":  "qa",
	}
	cleanup := setEnv(t, vars)
	defer cleanup()

	cfg, err := config.New()
	assert.NoError(t, err, "Expected no error when .env is missing but variables are set via system environment")

	assert.Equal(t, "4040", cfg.API.Port, "API.Port should be loaded from system environment")
	assert.Equal(t, "qa", cfg.API.Env, "API.Env should be loaded from system environment")
}
