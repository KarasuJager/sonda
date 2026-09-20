package report

import (
	"encoding/json"
	"testing"

	"github.com/KarasuJager/sonda/internal/diff"
	"github.com/KarasuJager/sonda/internal/fingerprint"
)

func TestJSONReport(t *testing.T) {
	baseline := fingerprint.Fingerprint{
		StatusCode:     200,
		BodyLength:     30,
		SHA256:         "aaaa",
		DurationMillis: 4,
	}

	mutated := fingerprint.Fingerprint{
		StatusCode:     200,
		BodyLength:     11,
		SHA256:         "bbbb",
		DurationMillis: 1,
	}

	comparison := diff.Result{
		StatusChanged:     false,
		LengthChanged:     true,
		LengthDelta:       -19,
		HashChanged:       true,
		BodySimilarity:    0.0,
		TimingDeltaMillis: -3,
	}

	r := New(
		"http://127.0.0.1:8080/?q=admin",
		"http://127.0.0.1:8080/?q=SONDA_TEST",
		"q",
		"SONDA_TEST",
		baseline,
		mutated,
		comparison,
	)

	data, err := JSON(r)
	if err != nil {
		t.Fatalf("encode report: %v", err)
	}

	if !json.Valid(data) {
		t.Fatal("generated report is not valid JSON")
	}

	var decoded Report

	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decode generated report: %v", err)
	}

	if decoded.Mutation.Parameter != "q" {
		t.Fatalf(
			"expected mutation parameter q, got %q",
			decoded.Mutation.Parameter,
		)
	}

	if decoded.Diff.LengthDelta != -19 {
		t.Fatalf(
			"expected length delta -19, got %d",
			decoded.Diff.LengthDelta,
		)
	}

	if !decoded.Diff.HashChanged {
		t.Fatal("expected hash change")
	}
}
