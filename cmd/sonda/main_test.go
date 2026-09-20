package main

import "testing"

func TestParseHeaders(t *testing.T) {
	headers, err := parseHeaders([]string{
		"Authorization: Bearer test-token",
		"X-Sonda-Test: first",
		"X-Sonda-Test: second",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := headers.Get("Authorization"); got != "Bearer test-token" {
		t.Fatalf("unexpected Authorization header: %q", got)
	}

	values := headers.Values("X-Sonda-Test")
	if len(values) != 2 || values[0] != "first" || values[1] != "second" {
		t.Fatalf("unexpected repeated header values: %q", values)
	}
}

func TestParseHeadersRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "missing separator", value: "X-Sonda-Test"},
		{name: "empty name", value: ": value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseHeaders([]string{tt.value}); err == nil {
				t.Fatalf("expected %q to be rejected", tt.value)
			}
		})
	}
}
