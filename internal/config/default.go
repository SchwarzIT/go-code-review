package config

import (
	"coupon_service/internal/mytypes"
	"time"
)

type environment string

const (
	SHUTDOWN_TIMEOUT = time.Duration(30) * time.Second
)

// WithDefaultEnv Option func to setup the default environment
func WithDefaultEnv(cfg *Config) error {
	cfg.API.ENV = mytypes.DefaultEnvironment
	return nil
}

// WithDefaultSHUTDOWN_TIMEOUT Option func to setup the default shutdown timeout
func WithDefaultSHUTDOWN_TIMEOUT(cfg *Config) error {
	cfg.API.SHUTDOWN_TIMEOUT = mytypes.MyDuration(SHUTDOWN_TIMEOUT)
	return nil
}

// WithDefaultAllowOrigins Option func to setup the default allow origins
func WithDefaultAllowOrigins(cfg *Config) error {
	cfg.API.ALLOW_ORIGINS = mytypes.DefaultAllowOrigins
	return nil
}

// NewDefault setup all default option func to create a new Config
func NewDefault(envFilePath string) (Config, error) {
	return New(envFilePath, WithDefaultEnv, WithDefaultSHUTDOWN_TIMEOUT, WithDefaultAllowOrigins)
}
