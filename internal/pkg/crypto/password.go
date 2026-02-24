package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const saltSize = 16

// HashPassword generates a random salt and returns "hex(salt)$hex(hmac-sha256(pepper, salt+password))".
func HashPassword(password, pepper string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("salt generation failed: %w", err)
	}

	hash := computeHMAC(pepper, salt, password)
	return hex.EncodeToString(salt) + "$" + hex.EncodeToString(hash), nil
}

// VerifyPassword checks a password against a stored "salt$hash" string.
func VerifyPassword(stored, password, pepper string) bool {
	parts := strings.SplitN(stored, "$", 2)
	if len(parts) != 2 {
		return false
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}

	storedHash, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	computed := computeHMAC(pepper, salt, password)
	return hmac.Equal(storedHash, computed)
}

func computeHMAC(pepper string, salt []byte, password string) []byte {
	mac := hmac.New(sha256.New, []byte(pepper))
	mac.Write(salt)
	mac.Write([]byte(password))
	return mac.Sum(nil)
}
