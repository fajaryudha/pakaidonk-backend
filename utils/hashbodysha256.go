package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// HashBodySHA256 hashes the given body with SHA-256 and returns it in lowercase hexadecimal
func HashBodySHA256(body []byte) string {
	hash := sha256.Sum256(body)
	return strings.ToLower(hex.EncodeToString(hash[:]))
}
