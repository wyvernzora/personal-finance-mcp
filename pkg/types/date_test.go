package types

import (
	"encoding/json"
	"testing"
	"time"
)

// TestDate_MarshalJSON_Zero verifies that the zero Date serializes to an empty JSON string.
func TestDate_MarshalJSON_Zero(t *testing.T) {
	var d Date
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("MarshalJSON failed for zero date: %v", err)
	}
	got := string(data)
	const want = `""`
	if got != want {
		t.Errorf("MarshalJSON zero: got %s, want %s", got, want)
	}
}

// TestDate_MarshalJSON_NonZero verifies that a non-zero Date serializes to "YYYY-MM-DD".
func TestDate_MarshalJSON_NonZero(t *testing.T) {
	tt := time.Date(3023, 7, 23, 15, 4, 5, 0, time.UTC)
	d := Date(tt)
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	got := string(data)
	const want = `"3023-07-23"`
	if got != want {
		t.Errorf("MarshalJSON non-zero: got %s, want %s", got, want)
	}
}

// TestDate_UnmarshalJSON_Empty verifies that an empty JSON string yields a zero Date.
func TestDate_UnmarshalJSON_Empty(t *testing.T) {
	var d Date
	err := json.Unmarshal([]byte(`""`), &d)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed for empty: %v", err)
	}
	if !time.Time(d).IsZero() {
		t.Errorf("UnmarshalJSON empty: got non-zero date %v", time.Time(d))
	}
}

// TestDate_UnmarshalJSON_NonZero verifies parsing a valid "YYYY-MM-DD" string.
func TestDate_UnmarshalJSON_NonZero(t *testing.T) {
	var d Date
	src := `"2136-04-17"`
	if err := json.Unmarshal([]byte(src), &d); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	got := time.Time(d)
	want := time.Date(2136, 4, 17, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("UnmarshalJSON non-zero: got %v, want %v", got, want)
	}
}

// TestDate_UnmarshalJSON_Invalid verifies that a malformed date string returns an error.
func TestDate_UnmarshalJSON_Invalid(t *testing.T) {
	var d Date
	err := json.Unmarshal([]byte(`"not-a-date"`), &d)
	if err == nil {
		t.Fatal("UnmarshalJSON invalid: expected error, got nil")
	}
}

// TestDate_JSON_RoundTrip verifies that marshaling then unmarshaling yields the original date.
func TestDate_JSON_RoundTrip(t *testing.T) {
	original := Date(time.Date(2500, 12, 31, 0, 0, 0, 0, time.UTC))
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	var parsed Date
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if !time.Time(parsed).Equal(time.Time(original)) {
		t.Errorf("roundtrip: got %v, want %v", parsed, original)
	}
}
