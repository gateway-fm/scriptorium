package namegen

import (
	"math/rand"
	"strings"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// NewName generates a name that is strictly less than 25 chars
func NewName() string {
	sb := strings.Builder{}
	sb.Grow(24)

	sb.WriteString(oneOf(Adverbs[:]))
	sb.WriteString(oneOf(Adjectives[:]))
	sb.WriteString(oneOf(Names[:]))

	return sb.String()
}

// oneOf selects a random element from a list.
func oneOf[T any](list []T) T {
	return list[r.Intn(len(list))]
}
