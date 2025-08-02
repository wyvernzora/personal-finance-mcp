package types

import (
	"fmt"
	"math"
	"strconv"
)

const PRECISION = 4  // Number of decimal places to keep (e.g., 4 means ten-thousandths)
const FACTOR = 10000 // 10^PRECISION, used to scale amounts to integer representation

// Money represents a fixed-precision decimal value as an integer number of (1/FACTOR) units.
// For PRECISION=4, a Money value of 12345 represents 1.2345 units.
// Provides safe arithmetic and JSON serialization with fixed precision.
type Money int64

// MarshalJSON implements json.Marshaler, formatting the Money value as a JSON number
// with exactly PRECISION decimals (e.g. "12.3400" for PRECISION=4).
func (m Money) MarshalJSON() ([]byte, error) {
	f := float64(m) / FACTOR
	s := strconv.FormatFloat(f, 'f', PRECISION, 64)
	return []byte(s), nil
}

// UnmarshalJSON implements json.Unmarshaler. It accepts a JSON number or quoted string,
// parses it to a float with up to PRECISION decimals, and rounds to the nearest fixed unit.
func (m *Money) UnmarshalJSON(data []byte) error {
	s := string(data)
	// Strip possible quotes
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		unquoted, err := strconv.Unquote(s)
		if err != nil {
			return fmt.Errorf("invalid quoted money %q: %w", s, err)
		}
		s = unquoted
	}
	// Allow both integer and float inputs
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid money value %q: %w", s, err)
	}
	units := int64(math.Round(f * FACTOR))
	*m = Money(units)
	return nil
}

// Add returns the sum of two Money values, preserving precision.
func (m Money) Add(o Money) Money {
	return m + o
}
