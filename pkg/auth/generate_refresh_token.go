package auth

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func GenerateRefreshToken(bytesLength uint) (string, error) {
	salt := make([]byte, bytesLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}