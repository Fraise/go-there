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
	SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error)
}

// GetHashFromPassword takes a password, and returns (complete bcrypt hash, salt only, error).
func GetHashFromPassword(password string) ([]byte, []byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bCryptCost)

	if err != nil {
		return nil, nil, err
	}

	hashArr := bytes.Split(hash, []byte("$"))

	// 0 = ""
	// 1 = Algorithm
	// 2 = Cost
	// 3 = Salt+Hash, the salt should be 22 bytes long and hash 31 bytes long

	return hash, hashArr[3][:22], nil
}

// GenerateRandomB64String creates a random base64 encoded string from using the crypto/rand package from a byte array
// of length n. If less than n random bytes are generated, an error is returned.
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

// GetRequestedUser returns the user corresponding to the ressource accessed. It returns "" if the ressource does not
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

// validateApiKey takes an api key with the salt encoded in b64 and returns (salt, apikey, error).
func validateApiKey(apiKey string) ([]byte, []byte, error) {
	apiKeyArr := bytes.Split([]byte(apiKey), []byte("."))

	if len(apiKeyArr) != 2 {
		return nil, nil, data.ErrInvalidKey
	}

	decodedSalt, err := base64.URLEncoding.DecodeString(string(apiKeyArr[0]))

	if err != nil {
		return nil, nil, data.ErrInvalidKey
	}

	return decodedSalt, apiKeyArr[1], nil
}
