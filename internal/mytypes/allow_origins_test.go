// File: mytypes/allow_origins_test.go
package mytypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAllowOrigins_UnmarshalText tests the UnmarshalText method of the AllowOrigins type with various inputs.
func TestAllowOrigins_UnmarshalText(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expected       AllowOrigins
		expectingError bool
	}{

		{
			name:     "Single valid origin",
			input:    "https://example.com",
			expected: AllowOrigins{"https://example.com"},
		},
		{
			name:     "Multiple valid origins",
			input:    "https://example.com, https://api.example.com",
			expected: AllowOrigins{"https://example.com", "https://api.example.com"},
		},
		{
			name:     "Origins with extra whitespace",
			input:    "  https://example.com  ,https://api.example.com ",
			expected: AllowOrigins{"https://example.com", "https://api.example.com"},
		},
		// Invalid Inputs
		{
			name:           "Empty input",
			input:          "",
			expected:       nil,
			expectingError: false,
		},
		{
			name:           "Only commas",
			input:          ",,,",
			expected:       nil,
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ao AllowOrigins
			err := ao.UnmarshalText([]byte(tt.input))

			if tt.expectingError {
				assert.Error(t, err, "Expected an error for input: %s", tt.input)
				assert.Equal(t, tt.expected, ao, "AllowOrigins mismatch for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %s", tt.input)
				assert.Equal(t, tt.expected, ao, "AllowOrigins mismatch for input: %s", tt.input)
			}
		})
	}
}
