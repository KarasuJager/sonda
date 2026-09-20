package fingerprint

import (
	"net/http"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	got := Generate(http.StatusOK, nil, 1500*time.Microsecond)

	if got.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got.StatusCode)
	}

	if got.BodyLength != 0 {
		t.Fatalf("expected empty body, got length %d", got.BodyLength)
	}

	const emptySHA256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if got.SHA256 != emptySHA256 {
		t.Fatalf("unexpected SHA-256: %s", got.SHA256)
	}

	if got.DurationMillis != 1 {
		t.Fatalf("expected truncated duration of 1 ms, got %d", got.DurationMillis)
	}
}
