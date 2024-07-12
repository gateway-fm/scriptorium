package namegen

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// NewName generates a name that is strictly less than 25 chars
func NewName() string {
	sb := strings.Builder{}
	sb.Grow(24)

	sb.WriteString(oneOf(Adjectives[:]))
	sb.WriteString(oneOf(Adverbs[:]))
	sb.WriteString(oneOf(Names[:]))

	return sb.String()
}

// oneOf selects a random element from a list.
func oneOf[T any](list []T) T {
	if len(list) <= 0 {
		var zero T
		return zero
	}

	maxVal := big.NewInt(int64(len(list)))
	// don't need to check err there, it only occurs if maxVal <= 0
	// but this case is handled in return above
	randomBigInt, _ := rand.Int(rand.Reader, maxVal)

	return list[int(randomBigInt.Int64())]
}
