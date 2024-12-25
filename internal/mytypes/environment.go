package mytypes

import (
	"errors"
	"strings"
)

// Environment represents the type of environment the application is running in.
type Environment string

const (
	// Development environment
	Development Environment = "development"
	// Production environment
	Production Environment = "production"

	// DefaultEnvironment is the fallback environment if none is provided or invalid.
	DefaultEnvironment = Development
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It allows configuration libraries to parse environment strings directly into the Environment type.
func (e *Environment) UnmarshalText(text []byte) error {
	parsedEnv := strings.ToLower(strings.TrimSpace(string(text)))

	switch parsedEnv {
	case string(Development), string(Production):
		*e = Environment(parsedEnv)
		return nil
	default:
		// Assign default if invalid environment is provided
		*e = DefaultEnvironment
		return errors.New("invalid environment provided, defaulting to development")
	}
}
