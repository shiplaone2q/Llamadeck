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

