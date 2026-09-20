package diff

import (
	"math"
	"testing"

	"github.com/KarasuJager/sonda/internal/fingerprint"
)

func TestCompare(t *testing.T) {
	before := fingerprint.Fingerprint{
		StatusCode:     200,
		BodyLength:     5,
		SHA256:         "aaaa",
		DurationMillis: 10,
	}

	after := fingerprint.Fingerprint{
		StatusCode:     200,
		BodyLength:     5,
		SHA256:         "bbbb",
		DurationMillis: 15,
	}

	result := Compare(
		before,
		after,
		[]byte("abcde"),
		[]byte("abXYZ"),
	)

	if result.StatusChanged {
		t.Fatal("status should not have changed")
	}

	if result.StatusBefore != 200 || result.StatusAfter != 200 {
		t.Fatalf(
			"unexpected statuses: %d -> %d",
			result.StatusBefore,
			result.StatusAfter,
		)
	}

	if result.LengthChanged {
		t.Fatal("length should not have changed")
	}

	if result.LengthBefore != 5 || result.LengthAfter != 5 || result.LengthDelta != 0 {
		t.Fatalf(
			"unexpected lengths: %d -> %d (delta %+d)",
			result.LengthBefore,
			result.LengthAfter,
			result.LengthDelta,
		)
	}

	if !result.HashChanged {
		t.Fatal("hash change was not detected")
	}

	if math.Abs(result.BodySimilarity-0.25) > 0.0001 {
		t.Fatalf(
			"expected similarity 0.25, got %.4f",
			result.BodySimilarity,
		)
	}

	if result.TimingDeltaMillis != 5 {
		t.Fatalf(
			"expected timing delta 5 ms, got %d",
			result.TimingDeltaMillis,
		)
	}
}
