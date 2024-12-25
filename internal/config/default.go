package config

import (
	"coupon_service/internal/mytypes"
	"time"
)

type environment string

const (
	Shutdown_Timeout = time.Duration(30) * time.Second
)

// WithDefaultEnv Option func to setup the default environment
func WithDefaultEnv(cfg *Config) error {
	cfg.API.Env = mytypes.DefaultEnvironment
	return nil
}

// WithDefaultShutdown_Timeout Option func to setup the default shutdown timeout
func WithDefaultShutdown_Timeout(cfg *Config) error {
	cfg.API.Shutdown_Timeout = mytypes.MyDuration(Shutdown_Timeout)
	return nil
}

// WithDefaultAllowOrigins Option func to setup the default allow origins
func WithDefaultAllowOrigins(cfg *Config) error {
	cfg.API.Allow_Origins = mytypes.DefaultAllowOrigins
	return nil
}

// NewDefault setup all default option func to create a new Config
func NewDefault(envFilePath string) (Config, error) {
	return New(envFilePath, WithDefaultEnv, WithDefaultShutdown_Timeout, WithDefaultAllowOrigins)
}
