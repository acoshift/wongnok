package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
