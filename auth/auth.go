package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"go-there/data"
	"golang.org/x/crypto/bcrypt"
)

const bCryptCost = bcrypt.DefaultCost

// DataSourcer is used to access the mysql database.
type DataSourcer interface {
	SelectUserLogin(username string) (data.User, error)
	SelectUserLoginByApiKeyHash(apiKeyHash string) (data.User, error)
	GetAuthToken(token string) (data.AuthToken, error)
}

// GetHashFromPassword takes a password, and returns (complete bcrypt hash, error).
func GetHashFromPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bCryptCost)

	if err != nil {
		return nil, err
	}

	return hash, nil
}

// GenerateRandomB64String creates a random base64 URL encoded string from using the crypto/rand package from a byte
// array of length n. If less than n random bytes are generated, an error is returned.
func GenerateRandomB64String(n int) (string, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	// If less than n bytes are read, an error is returned
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// GetLoggedUser returns the currently logged user, or an empty User otherwise.
func GetLoggedUser(c *gin.Context) data.User {
	if c.Keys == nil {
		return data.User{}
	}

	u, ok := c.Keys["user"].(data.User)

	if !ok {
		return data.User{}
	}

	return u
}

// GetRequestedUser returns the user corresponding to the resource accessed. It returns "" if the resource does not
// belong to any user.
func GetRequestedUser(c *gin.Context) string {
	if c.Keys == nil {
		return ""
	}

	u, ok := c.Keys["reqUser"].(string)

	if !ok {
		return ""
	}

	return u
}

// validateApiKey takes an api key with the hash encoded in b64 and returns (hash, apikey, error).
func validateApiKey(apiKey string) ([]byte, []byte, error) {
	decodedKey, err := base64.URLEncoding.DecodeString(apiKey)

	if err != nil {
		return nil, nil, data.ErrInvalidKey
	}

	apiKeyArr := bytes.Split(decodedKey, []byte(":"))

	if len(apiKeyArr) != 2 {
		return nil, nil, data.ErrInvalidKey
	}

	return apiKeyArr[0], apiKeyArr[1], nil
}
