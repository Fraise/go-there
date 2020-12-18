package api

import (
	"encoding/base64"
	"go-there/auth"
	"testing"
)

// BenchmarkUserCreation mostly used to generate test passwords
func BenchmarkUserCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// We don't need to store the salt section individually for a password
		hash, _, err := auth.GetHashFromPassword("superpassword")

		if err != nil {
			return
		}

		// Generate a random API key
		apiKey, err := auth.GenerateRandomB64String(16)

		if err != nil {
			return
		}

		// Get its corresponding hash and salt
		apiKeyHash, apiKeySalt, _ := auth.GetHashFromPassword(apiKey)

		fullApiKey := base64.URLEncoding.EncodeToString(append(apiKeySalt, []byte(":"+apiKey)...))

		_ = hash
		_ = apiKeyHash
		_ = apiKeySalt
		_ = fullApiKey
	}
}
