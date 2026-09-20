package similarity

import (
	"math"
	"testing"
)

func TestDiceIdentical(t *testing.T) {
	got := Dice(
		[]byte("hello SONDA"),
		[]byte("hello SONDA"),
	)

	if got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
}

func TestDiceDifferent(t *testing.T) {
	got := Dice(
		[]byte("abcde"),
		[]byte("abXYZ"),
	)

	expected := 0.25

	if math.Abs(got-expected) > 0.0001 {
		t.Fatalf(
			"expected %.2f, got %.4f",
			expected,
			got,
		)
	}
}

func TestDiceNoSimilarity(t *testing.T) {
	got := Dice(
		[]byte("AAAA"),
		[]byte("ZZZZ"),
	)

	if got != 0.0 {
		t.Fatalf("expected 0.0, got %f", got)
	}
}
