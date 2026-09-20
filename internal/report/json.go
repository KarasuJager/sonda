package report

import (
	"encoding/json"

	"github.com/KarasuJager/sonda/internal/diff"
	"github.com/KarasuJager/sonda/internal/fingerprint"
)

type Mutation struct {
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
}

type Snapshot struct {
	Status         int    `json:"status"`
	Length         int    `json:"length"`
	SHA256         string `json:"sha256"`
	DurationMillis int64  `json:"duration_ms"`
}

type Difference struct {
	StatusChanged     bool    `json:"status_changed"`
	LengthChanged     bool    `json:"length_changed"`
	LengthDelta       int     `json:"length_delta"`
	HashChanged       bool    `json:"hash_changed"`
	BodySimilarity    float64 `json:"body_similarity"`
	TimingDeltaMillis int64   `json:"timing_delta_ms"`
}

type Report struct {
	Target     string     `json:"target"`
	MutatedURL string     `json:"mutated_url"`
	Mutation   Mutation   `json:"mutation"`
	Baseline   Snapshot   `json:"baseline"`
	Mutated    Snapshot   `json:"mutated"`
	Diff       Difference `json:"diff"`
}

func New(
	target string,
	mutatedURL string,
	parameter string,
	value string,
	baseline fingerprint.Fingerprint,
	mutated fingerprint.Fingerprint,
	comparison diff.Result,
) Report {
	return Report{
		Target:     target,
		MutatedURL: mutatedURL,
		Mutation: Mutation{
			Parameter: parameter,
			Value:     value,
		},
		Baseline: snapshot(baseline),
		Mutated:  snapshot(mutated),
		Diff: Difference{
			StatusChanged:     comparison.StatusChanged,
			LengthChanged:     comparison.LengthChanged,
			LengthDelta:       comparison.LengthDelta,
			HashChanged:       comparison.HashChanged,
			BodySimilarity:    comparison.BodySimilarity,
			TimingDeltaMillis: comparison.TimingDeltaMillis,
		},
	}
}

func JSON(r Report) ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func snapshot(fp fingerprint.Fingerprint) Snapshot {
	return Snapshot{
		Status:         fp.StatusCode,
		Length:         fp.BodyLength,
		SHA256:         fp.SHA256,
		DurationMillis: fp.DurationMillis,
	}
}
