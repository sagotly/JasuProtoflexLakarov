package cryptoHelp

import (
	"crypto/rand"
	"encoding/hex"
)

// generateToken creates a new random token.
func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
