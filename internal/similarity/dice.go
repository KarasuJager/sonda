package similarity

import "bytes"

// Dice calculates the Sorensen-Dice similarity coefficient between
// two byte slices using byte bigrams.
//
// The returned value is between 0.0 and 1.0:
//
//	1.0 = identical
//	0.0 = no bigram similarity
func Dice(a, b []byte) float64 {
	if bytes.Equal(a, b) {
		return 1.0
	}

	if len(a) < 2 || len(b) < 2 {
		return 0.0
	}

	counts := make(map[uint16]int)

	for i := 0; i < len(a)-1; i++ {
		key := bigram(a[i], a[i+1])
		counts[key]++
	}

	intersection := 0

	for i := 0; i < len(b)-1; i++ {
		key := bigram(b[i], b[i+1])

		if counts[key] > 0 {
			intersection++
			counts[key]--
		}
	}

	total := (len(a) - 1) + (len(b) - 1)

	return (2.0 * float64(intersection)) / float64(total)
}

func bigram(first, second byte) uint16 {
	return uint16(first)<<8 | uint16(second)
}
