package hash

import (
	"testing"
)

const testString = "examplestringtohash"

// nolint
func Benchmark_Argon2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateHashArgon2(testString)
	}
}

// nolint
func Benchmark_Blake2b(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateHashBlake2b(testString)
	}
}

// nolint
func Benchmark_Sha512(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateHashSha512(testString)
	}
}
