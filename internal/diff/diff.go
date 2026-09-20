package diff

import (
	"github.com/KarasuJager/sonda/internal/fingerprint"
	"github.com/KarasuJager/sonda/internal/similarity"
)

type Result struct {
	StatusChanged bool
	StatusBefore  int
	StatusAfter   int

	LengthChanged bool
	LengthBefore  int
	LengthAfter   int
	LengthDelta   int

	HashChanged bool

	BodySimilarity float64

	TimingDeltaMillis int64
}

func Compare(
	before fingerprint.Fingerprint,
	after fingerprint.Fingerprint,
	beforeBody []byte,
	afterBody []byte,
) Result {
	return Result{
		StatusChanged: before.StatusCode != after.StatusCode,
		StatusBefore:  before.StatusCode,
		StatusAfter:   after.StatusCode,

		LengthChanged: before.BodyLength != after.BodyLength,
		LengthBefore:  before.BodyLength,
		LengthAfter:   after.BodyLength,
		LengthDelta:   after.BodyLength - before.BodyLength,

		HashChanged: before.SHA256 != after.SHA256,

		BodySimilarity: similarity.Dice(
			beforeBody,
			afterBody,
		),

		TimingDeltaMillis: after.DurationMillis - before.DurationMillis,
	}
}
