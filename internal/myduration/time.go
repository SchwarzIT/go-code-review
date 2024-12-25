package myduration

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	HoursInDay = 24
	DaysInWeek = 7
	DaysInYear = 365
)

// MyDuration is a custom duration type that handles extended time units
type MyDuration time.Duration

// ParseDuration convert MyDuration to time.Duration
func (md MyDuration) ParseTimeDuration() time.Duration {
	return time.Duration(md)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// This allows configuration libraries to parse duration strings directly into MyDuration.
func (md *MyDuration) UnmarshalText(text []byte) error {
	duration, err := ParseDuration(string(text))
	if err != nil {
		return err
	}
	*md = MyDuration(duration)
	return nil
}

// ParseMyDuration parses a duration string with units like s, m, h, d, w, y.
func ParseDuration(sdurantion string) (time.Duration, error) {
	re := regexp.MustCompile(`(-?\d+)([smhdwy])`)
	matches := re.FindAllStringSubmatch(sdurantion, -1)

	if matches == nil {
		return 0, fmt.Errorf("invalid duration format")
	}

	var totalDuration time.Duration
	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("invalid number: %v", match[1])
		}

		if value < 0 {
			return 0, fmt.Errorf("invalid value value is negative: %v", match[1])
		}

		unit := match[2]
		switch unit {
		case "s":
			totalDuration += time.Duration(value) * time.Second
		case "m":
			totalDuration += time.Duration(value) * time.Minute
		case "h":
			totalDuration += time.Duration(value) * time.Hour
		case "d":
			totalDuration += time.Duration(value) * time.Hour * HoursInDay
		case "w":
			totalDuration += time.Duration(value) * time.Hour * HoursInDay * DaysInWeek
		case "y":
			totalDuration += time.Duration(value) * time.Hour * HoursInDay * DaysInYear
		default:
			return 0, fmt.Errorf("unknown time unit: %v", unit)
		}
	}

	return totalDuration, nil
}
