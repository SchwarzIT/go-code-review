package mytypes

import (
	"fmt"
	"net/url"
	"strings"
)

// AllowOrigins represents a list of allowed origins for CORS.
type AllowOrigins []string

var DefaultAllowOrigins = []string{"*"}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It parses a comma-separated string into an AllowOrigins slice,
// validating each origin as a proper URL.
func (ao *AllowOrigins) UnmarshalText(text []byte) error {
	input := strings.TrimSpace(string(text))

	if input == "" {
		return nil
	}

	parts := strings.Split(input, ",")

	var origins []string

	for idx, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			return fmt.Errorf("origin at position %d is empty", idx)
		}

		parsedURL, err := url.Parse(origin)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			return fmt.Errorf("invalid origin at position %d: %s", idx, origin)
		}

		origins = append(origins, origin)
	}

	*ao = origins
	return nil
}

// String returns the AllowOrigins as a comma-separated string.
func (ao AllowOrigins) String() string {
	return strings.Join(ao, ",")
}
