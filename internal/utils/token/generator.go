package token

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateToken generates a random token
func GenerateToken(length int) (string, error) {
	if length <= 0 {
		length = 32
	}

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
