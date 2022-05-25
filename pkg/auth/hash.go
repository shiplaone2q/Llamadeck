package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2 parameters
	timeCost uint32 = 1
	memory   uint32 = 64 * 1024 // 64 MB
	threads  uint8  = 4
	keyLen   uint32 = 32
	saltLen  uint32 = 16
)

// HashKey hashes an API key using Argon2id
func HashKey(plainTextKey string) (string, error) {
	// Generate random salt
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key using Argon2id
	hash := argon2.IDKey([]byte(plainTextKey), salt, timeCost, memory, threads, keyLen)

	// Format: $argon2id$v=19$m=65536,t=1,p=4$<base64-salt>$<base64-hash>
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memory, timeCost, threads, saltB64, hashB64), nil
