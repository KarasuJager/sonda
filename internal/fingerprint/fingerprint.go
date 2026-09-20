package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Fingerprint struct {
	StatusCode     int
	BodyLength     int
	SHA256         string
	DurationMillis int64
}

func Generate(statusCode int, body []byte, duration time.Duration) Fingerprint {
	hash := sha256.Sum256(body)

	return Fingerprint{
		StatusCode:     statusCode,
		BodyLength:     len(body),
		SHA256:         hex.EncodeToString(hash[:]),
		DurationMillis: duration.Milliseconds(),
	}
}
