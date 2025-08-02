package types

import (
	"strings"
	"time"
)

// Date represents a calendar date (without time) in YYYY-MM-DD format.
// It wraps time.Time to provide custom JSON marshaling and unmarshaling.
// The zero value of Date serializes to an empty string.
type Date time.Time

// MarshalJSON implements json.Marshaler. It encodes the Date as a JSON string
// using the "YYYY-MM-DD" layout. Zero-value dates serialize to an empty string.
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte(`""`), nil
	}
	s := t.Format("2006-01-02")
	return []byte(`"` + s + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler. It decodes a JSON string
// formatted as "YYYY-MM-DD" into the Date. Empty or missing strings will
// result in a zero-value Date (time.Time{}).
func (d *Date) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		*d = Date(time.Time{})
		return nil
	}
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// Parse parses a string in "YYYY-MM-DD" format into a Date. Empty string returns a zero Date.
func ParseDate(s string) (Date, error) {
	if s == "" {
		return Date(time.Time{}), nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date(time.Time{}), err
	}
	return Date(t), nil
}

// String returns the Date formatted as "YYYY-MM-DD". Zero-value dates
// are returned as an empty string.
func (d Date) String() string {
	t := time.Time(d)
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}
