package mutation

import (
	"net/url"
	"testing"
)

func TestQueryParam(t *testing.T) {
	got, err := QueryParam(
		"http://example.test/search?q=admin&page=1",
		"q",
		"SONDA_TEST",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse mutated URL: %v", err)
	}

	if u.Query().Get("q") != "SONDA_TEST" {
		t.Fatalf("expected q=SONDA_TEST, got %q", u.Query().Get("q"))
	}

	if u.Query().Get("page") != "1" {
		t.Fatalf("unrelated parameter was modified")
	}
}

func TestQueryParamMissingParameter(t *testing.T) {
	_, err := QueryParam(
		"http://example.test/search?q=admin",
		"missing",
		"SONDA_TEST",
	)

	if err == nil {
		t.Fatal("expected error for missing parameter")
	}
}
