package hash

import (
	"encoding/base64"

	"golang.org/x/crypto/sha3"

	"golang.org/x/crypto/blake2b"

	"golang.org/x/crypto/argon2"
)

type hashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

// constant salt from GUID: f060a42abac94a3bab8d1e7a9ba1d556
var salt = []byte{102, 48, 54, 48, 97, 52, 50, 97, 98, 97, 99, 57, 52, 97, 51, 98, 97, 98, 56, 100, 49, 101, 55, 97, 57, 98, 97, 49, 100, 53, 53, 54}

func getHashParams() *hashParams {
	return &hashParams{
		memory:      1 * 1024,
		iterations:  2,
		parallelism: 2,
		keyLength:   32,
	}
}

// GenerateHashArgon2 hashes the secret and returns base64 of the hash
func GenerateHashArgon2(secret string) (encodedHash string, err error) {
	p := getHashParams()
	hash := argon2.Key([]byte(secret), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	return b64Hash, nil
}

// GenerateHashBlake2b hashes the secret and returns base64 of the hash
func GenerateHashBlake2b(secret string) (encodedHash string) {
	h := blake2b.Sum512([]byte(secret))
	return base64.URLEncoding.EncodeToString(h[:])
}

// GenerateHashSha512 hashes the secret and returns base64 of the hash
func GenerateHashSha512(secret string) (encodedHash string) {
	h := sha3.Sum512([]byte(secret))
	return base64.URLEncoding.EncodeToString(h[:])
}
