// internal/config/config_test.go
package config_test

import (
	"coupon_service/internal/config"
	"coupon_service/internal/mytypes"
	"os"
	"path/filepath"
	"testing"
	"time"

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
				clearEnvVars(t)

				envContent := "API_PORT=9090\nAPI_ENV=production\nAPI_TIME_ALIVE=2y\nAPI_ALLOW_ORIGINS=https://example.com,https://api.example.com\nAPI_SHUTDOWN_TIMEOUT=20s"
				envPath := createTempEnvFile(t, envContent)
				cfg, err := config.NewDefault(envPath)
				assert.NoError(t, err, "Expected no error when loading config from .env")

				assert.Equal(t, "9090", cfg.API.Port, "API.Port should be loaded from .env")
				assert.Equal(t, mytypes.Production, cfg.API.Env, "API.Env should be loaded from .env")
				assert.Len(t, cfg.API.Allow_Origins, 2, "API.AllowOrigin should be loaded from .env")
				assert.Equal(t, time.Duration(2)*time.Hour*mytypes.HoursInDay*mytypes.DaysInYear, cfg.API.Time_Alive.ParseTimeDuration(), "API.Time_Alive should be loaded from system environment")
				assert.Equal(t, time.Duration(20)*time.Second, cfg.API.Shutdown_Timeout.ParseTimeDuration(), "API.Time_Alive should be loaded from system environment")

				t.Cleanup(func() {
					clearEnvVars(t)
				})
			},
		},
		{
			name: "FromSystemEnv",
			test: func(t *testing.T) {
				clearEnvVars(t)

				t.Setenv("API_PORT", "7070")
				t.Setenv("API_ENV", "development")
				t.Setenv("API_TIME_ALIVE", "2y")

				cfg, err := config.NewDefault("")
				assert.NoError(t, err, "Expected no error when loading config from system environment")

				assert.Equal(t, "7070", cfg.API.Port, "API.Port should be loaded from system environment")
				assert.Equal(t, mytypes.Development, cfg.API.Env, "API.Env should be loaded from system environment")
				assert.Equal(t, time.Duration(2)*time.Hour*mytypes.HoursInDay*mytypes.DaysInYear, cfg.API.Time_Alive.ParseTimeDuration(), "API.Time_Alive should be loaded from system environment")
			},
		},
		{
			name: "OverrideEnv",
			test: func(t *testing.T) {
				clearEnvVars(t)

				envContent := "API_PORT=8080\nAPI_ENV=development\n"
				envPath := createTempEnvFile(t, envContent)

				t.Setenv("API_PORT", "6060")

				cfg, err := config.NewDefault(envPath)
				assert.NoError(t, err, "Expected no error when loading config with override")

				assert.Equal(t, "6060", cfg.API.Port, "API.Port should be overridden by system environment")
				assert.Equal(t, mytypes.Development, cfg.API.Env, "API.Env should be loaded from .env")
			},
		},
		{
			name: "MissingCriticalPort",
			test: func(t *testing.T) {
				clearEnvVars(t)

				t.Setenv("API_ENV", "development")

				cfg, err := config.NewDefault("")
				assert.Error(t, err, "Expected error due to missing API_PORT")
				assert.Contains(t, err.Error(), "critical environment variable API_PORT is missing", "Error message should indicate missing API_PORT")
				assert.Empty(t, cfg.API.Port, "API.Port should be empty when missing")
			},
		},
		{
			name: "InvalidEnvironment",
			test: func(t *testing.T) {
				clearEnvVars(t)

				t.Setenv("API_ENV", "matheus")

				t.Setenv("API_PORT", "6060")

				_, err := config.NewDefault("")
				assert.Error(t, err, "Expected error due to missing API_PORT")
				assert.Contains(t, err.Error(), "invalid environment provided, defaulting to development", "Error message should indicate invalid ENV")
			},
		},
		{
			name: "DefaultValues",
			test: func(t *testing.T) {
				clearEnvVars(t)

				t.Setenv("API_PORT", "5050")

				cfg, err := config.NewDefault("")
				assert.NoError(t, err, "Expected no error when loading config with default values")

				assert.Equal(t, "5050", cfg.API.Port, "API.Port should be loaded from system environment")
				assert.Equal(t, mytypes.Development, cfg.API.Env, "API.Env should use default value 'dev'")
			},
		},
		{
			name: "NoDotEnv",
			test: func(t *testing.T) {

				clearEnvVars(t)

				t.Setenv("API_PORT", "4040")
				t.Setenv("API_ENV", "development")

				cfg, err := config.NewDefault("")
				assert.NoError(t, err, "Expected no error when .env is missing but variables are set via system environment")

				assert.Equal(t, "4040", cfg.API.Port, "API.Port should be loaded from system environment")
				assert.Equal(t, mytypes.Development, cfg.API.Env, "API.Env should be loaded from system environment")
			},
		},
	}

	for _, tc := range subtests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.test(t)
		})
	}
}
