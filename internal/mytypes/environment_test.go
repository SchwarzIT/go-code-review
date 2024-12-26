// File: mytypes/environment_test.go
package mytypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnvironment_UnmarshalText tests the UnmarshalText method of the Environment type with various inputs.
func TestEnvironment_UnmarshalText(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expected       Environment
		expectingError bool
	}{
		// Valid environment strings
		{"Valid development", "development", Development, false},
		{"Valid production", "production", Production, false},
		{"Valid uppercase DEVELOPMENT", "DEVELOPMENT", Development, false},
		{"Valid uppercase PRODUCTION", "PRODUCTION", Production, false},
		{"Valid mixed case", "DeVeLoPmEnT", Development, false},
		{"Valid with surrounding whitespace", "  production  ", Production, false},

		// Invalid environment strings
		{"Invalid environment", "staging1", DefaultEnvironment, true},
		{"Empty string", "", DefaultEnvironment, true},
		{"Unknown environment", "test", DefaultEnvironment, true},
		{"Numeric input", "123", DefaultEnvironment, true},
		{"Special characters", "@production!", DefaultEnvironment, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var env Environment
			err := env.UnmarshalText([]byte(tt.input))

			if tt.expectingError {
				assert.Error(t, err, "Expected an error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %s", tt.input)
			}

			assert.Equal(t, tt.expected, env, "Parsed environment mismatch for input: %s", tt.input)
		})
	}
}

// BenchmarkEnvironment_UnmarshalText benchmarks the UnmarshalText method of the Environment type.
func BenchmarkEnvironment_UnmarshalText(b *testing.B) {
	inputs := []string{
		"development",
		"production",
		"DeVeLoPmEnT",
		"PRODUCTION",
		"  production  ",
		"staging1",
	}

	for _, input := range inputs {
		b.Run(input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var env Environment
				err := env.UnmarshalText([]byte(input))
				if err != nil && input == "staging" {
					// Expected error for invalid input; continue
					continue
				} else if err != nil {
					b.Fatalf("Unexpected error for input %s: %v", input, err)
				}
			}
		})
	}
}
