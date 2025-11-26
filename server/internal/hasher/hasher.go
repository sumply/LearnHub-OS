package hasher

import (
	"crypto/rand"
	"crypto/subtle"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength = 16
	keyLength  = 32
	time       = 1
	memory     = 64 << 10
	threads    = 4
)

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string) (hash []byte, salt []byte, err error) {
	salt, err = generateSalt(saltLength)
	if err != nil {
		return nil, nil, err
	}

	hashed := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLength)

	return hashed, salt, nil
}

func VerifyPassword(password string, salt []byte, expectedHash []byte) bool {
	hashed := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLength)
	return subtleCompare(hashed, expectedHash)
}

func subtleCompare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
