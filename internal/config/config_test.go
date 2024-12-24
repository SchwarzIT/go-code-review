// internal/config/config_test.go
package config_test

import (
	"coupon_service/internal/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// List of environment variables used in the tests
var envVars = []string{"API_PORT", "API_ENV"}

// Helper function to clear specified environment variables
func clearEnvVars(t *testing.T) {
	for _, key := range envVars {
		err := os.Unsetenv(key)
		if err != nil {
			t.Fatalf("Failed to unset environment variable %s: %v", key, err)
		}
	}
}

// Helper function to create a temporary .env file
func createTempEnvFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	err := os.WriteFile(envPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary .env file: %v", err)
	}
	return envPath
}
func TestNewConfig(t *testing.T) {
	// Define a list of subtests with their respective names and test functions
	subtests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "FromDotEnv",
			test: func(t *testing.T) {
				// Clear environment variables
				clearEnvVars(t)

				// Create a temporary .env file
				envContent := "API_PORT=9090\nAPI_ENV=production\n"
				envPath := createTempEnvFile(t, envContent)

				// Load the .env file
				cfg, err := config.New(envPath)
				assert.NoError(t, err, "Expected no error when loading config from .env")

				// Assert the values loaded from .env
				assert.Equal(t, "9090", cfg.API.Port, "API.Port should be loaded from .env")
				assert.Equal(t, "production", cfg.API.Env, "API.Env should be loaded from .env")

				// Cleanup: Unset environment variables set by .env
				t.Cleanup(func() {
					clearEnvVars(t)
				})
			},
		},
		{
			name: "FromSystemEnv",
			test: func(t *testing.T) {
				// Clear environment variables
				clearEnvVars(t)

				// Set environment variables using t.Setenv for automatic cleanup
				t.Setenv("API_PORT", "7070")
				t.Setenv("API_ENV", "staging")

				// Load config without a .env file
				cfg, err := config.New("")
				assert.NoError(t, err, "Expected no error when loading config from system environment")

				// Assert the values loaded from system environment
				assert.Equal(t, "7070", cfg.API.Port, "API.Port should be loaded from system environment")
				assert.Equal(t, "staging", cfg.API.Env, "API.Env should be loaded from system environment")
			},
		},
		{
			name: "OverrideEnv",
			test: func(t *testing.T) {
				// Clear environment variables
				clearEnvVars(t)

				// Create a temporary .env file
				envContent := "API_PORT=8080\nAPI_ENV=development\n"
				envPath := createTempEnvFile(t, envContent)

				// Set environment variables that should override .env
				t.Setenv("API_PORT", "6060")
				// Note: Not setting API_ENV to ensure it comes from .env

				// Load the .env file with overrides
				cfg, err := config.New(envPath)
				assert.NoError(t, err, "Expected no error when loading config with override")

				// Assert that API_PORT is overridden and API_ENV comes from .env
				assert.Equal(t, "6060", cfg.API.Port, "API.Port should be overridden by system environment")
				assert.Equal(t, "development", cfg.API.Env, "API.Env should be loaded from .env")
			},
		},
		{
			name: "MissingCriticalPort",
			test: func(t *testing.T) {
				// Clear environment variables
				clearEnvVars(t)

				// Set only API_ENV to test missing API_PORT
				t.Setenv("API_ENV", "test")

				// Attempt to load config without API_PORT
				cfg, err := config.New("")
				assert.Error(t, err, "Expected error due to missing API_PORT")
				assert.Contains(t, err.Error(), "critical environment variable API_PORT is missing", "Error message should indicate missing API_PORT")
				assert.Empty(t, cfg.API.Port, "API.Port should be empty when missing")
			},
		},
		{
			name: "DefaultValues",
			test: func(t *testing.T) {
				// Clear environment variables
				clearEnvVars(t)

				// Set only API_PORT to test default for API_ENV
				t.Setenv("API_PORT", "5050")
				// Do not set API_ENV to use default

				// Load config without a .env file
				cfg, err := config.New()
				assert.NoError(t, err, "Expected no error when loading config with default values")

				// Assert that API_ENV uses the default value
				assert.Equal(t, "5050", cfg.API.Port, "API.Port should be loaded from system environment")
				assert.Equal(t, "dev", cfg.API.Env, "API.Env should use default value 'dev'")
			},
		},
		{
			name: "NoDotEnv",
			test: func(t *testing.T) {
				// Clear environment variables
				clearEnvVars(t)

				// Set environment variables without a .env file
				t.Setenv("API_PORT", "4040")
				t.Setenv("API_ENV", "qa")

				// Load config without a .env file
				cfg, err := config.New("")
				assert.NoError(t, err, "Expected no error when .env is missing but variables are set via system environment")

				// Assert that values are loaded from system environment
				assert.Equal(t, "4040", cfg.API.Port, "API.Port should be loaded from system environment")
				assert.Equal(t, "qa", cfg.API.Env, "API.Env should be loaded from system environment")
			},
		},
	}

	// Iterate over each subtest and execute them in sequence
	for _, tc := range subtests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}
