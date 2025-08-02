package types

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		input    Money
		expected string
	}{
		{"zero", Money(0), "0.0000"},
		{"positive", Money(12345678), "1234.5678"},
		{"negative", Money(-567890), "-56.7890"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, err := json.Marshal(c.input)
			if err != nil {
				t.Fatalf("MarshalJSON error: %v", err)
			}
			got := string(data)
			if got != c.expected {
				t.Errorf("got %q, want %q", got, c.expected)
			}
		})
	}
}

func TestUnmarshalJSON_Number(t *testing.T) {
	var m Money
	err := json.Unmarshal([]byte("12.3400"), &m)
	if err != nil {
		t.Fatalf("Unmarshal number error: %v", err)
	}
	want := Money(123400)
	if m != want {
		t.Errorf("got %d, want %d", m, want)
	}
}

func TestUnmarshalJSON_QuotedString(t *testing.T) {
	var m Money
	err := json.Unmarshal([]byte(`"56.7800"`), &m)
	if err != nil {
		t.Fatalf("Unmarshal quoted string error: %v", err)
	}
	want := Money(567800)
	if m != want {
		t.Errorf("got %d, want %d", m, want)
	}
}

func TestUnmarshalJSON_Rounding(t *testing.T) {
	cases := []struct {
		input    string
		expected Money
	}{
		{"12.3456", Money(123456)},
		{"12.3444", Money(123444)},
		{"-1.2356", Money(-12356)},
		{"-1.2344", Money(-12344)},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			var m Money
			if err := json.Unmarshal([]byte(c.input), &m); err != nil {
				t.Errorf("Unmarshal %q error: %v", c.input, err)
				return
			}
			if m != c.expected {
				t.Errorf("Unmarshal %q = %d, want %d", c.input, m, c.expected)
			}
		})
	}
}

func TestUnmarshalJSON_Invalid(t *testing.T) {
	inputs := []string{
		`"abc"`,
		"xyz",
		`"12.3.4"`,
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			var m Money
			if err := json.Unmarshal([]byte(in), &m); err == nil {
				t.Errorf("expected error for input %q, got nil", in)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	a := Money(100)
	b := Money(250)
	if got := a.Add(b); got != Money(350) {
		t.Errorf("Add = %d, want %d", got, 350)
	}
}
