package namegen_test

import (
	"log"
	"testing"

	"github.com/gateway-fm/scriptorium/namegen"
	"github.com/stretchr/testify/require"
)

// BenchmarkNewName benchmarks the NewName function.
func BenchmarkNewName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = namegen.NewName()
	}
}

// TestCollisionRate tests the collision rate of the NewName function.
func TestCollisionRate(t *testing.T) {
	const numTries = 1_000_000
	names := make(map[string]struct{}, numTries)
	collisions := 0

	for i := 0; i < numTries; i++ {
		name := namegen.NewName()
		if _, exists := names[name]; exists {
			collisions++
		} else {
			names[name] = struct{}{}
		}
	}

	rate := float64(collisions) / float64(numTries) * 100

	t.Logf("Number of collisions: %d", collisions)
	t.Logf("Collision rate: %.6f%%", rate)

	require.LessOrEqual(t, rate, 0.05, "should be less than 0.05% for 1 mil values")
}

func TestNames(t *testing.T) {
	for i := 0; i < 11; i++ {
		name := namegen.NewName()
		log.Println(name, len(name))
	}
}

func TestDictIsValid(t *testing.T) {
	var shouldFail bool

	for name, array := range map[string][]string{
		"names":      namegen.Names[:],
		"adjectives": namegen.Adjectives[:],
		"adverbs":    namegen.Adverbs[:]} {

		if len(array) < 300 {
			t.Logf("array %s is less than 300 elements (%d), consider making it bigger", name, len(array))
			shouldFail = true
		}

		set := make(map[string]struct{}, len(array))

		for _, value := range array {
			if len(value) > 8 {
				t.Logf("array %s has value %s, that is too long (%d)!", name, value, len(value))
				shouldFail = true
			}
			if _, exists := set[value]; exists {
				shouldFail = true
				t.Logf("value %s is duplicated!", value)

			}
			set[value] = struct{}{}
		}

		if len(set) != len(array) {
			t.Logf("array %s has duplicates!", name)
			shouldFail = true
		}
	}

	if shouldFail {
		t.Fatalf("one of the dict checks has failed, please check")
	}
}
