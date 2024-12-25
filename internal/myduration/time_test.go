package myduration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestParseMyDuration tests the ParseMyDuration function with various inputs.
func TestParseMyDuration(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expected       time.Duration
		expectingError bool
	}{
		// Valid duration strings
		{"Single seconds", "10s", 10 * time.Second, false},
		{"Single minutes", "5m", 5 * time.Minute, false},
		{"Single hours", "2h", 2 * time.Hour, false},
		{"Single days", "1d", 24 * time.Hour, false},
		{"Single weeks", "3w", 3 * 7 * 24 * time.Hour, false},
		{"Single years", "1y", 365 * 24 * time.Hour, false},
		{"Mixed hours and minutes", "1h30m", time.Hour + 30*time.Minute, false},
		{"Mixed days and hours", "2d4h", 2*24*time.Hour + 4*time.Hour, false},
		{"Complex mixed units", "1w2d3h4m5s", (1 * 7 * 24 * time.Hour) + (2 * 24 * time.Hour) + (3 * time.Hour) + (4 * time.Minute) + (5 * time.Second), false},
		// Invalid duration strings
		{"Empty string", "", 0, true},
		{"Missing unit", "10", 0, true},
		{"Unknown unit", "10x", 0, true},
		{"Negative duration", "-5m", 0, true}, // Assuming negative durations are invalid

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDuration(tt.input)
			if tt.expectingError {
				assert.Error(t, err, "Expected an error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %s", tt.input)
				assert.Equal(t, tt.expected, result, "Parsed duration mismatch for input: %s", tt.input)
			}
		})
	}
}

// BenchmarkParseMyDuration benchmarks the ParseMyDuration function.
func BenchmarkParseMyDuration(b *testing.B) {
	input := "1w2d3h4m5s"
	for i := 0; i < b.N; i++ {
		_, err := ParseDuration(input)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
