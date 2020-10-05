package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-there/data"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const bCryptCost = bcrypt.DefaultCost

// DataSourcer is used to access the mysql database.
type DataSourcer interface {
	SelectUser(username string) (data.User, error)
	SelectUserLogin(username string) (data.User, error)
	SelectApiKeyHashByUser(username string) ([]byte, error)
	SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error)
	SelectApiKeyHashBySalt(apiKeySalt string) ([]byte, error)
	InsertUser(user data.User) error
	DeleteUser(username string) error
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

func GenerateRandomB64String(n int) (string, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	// If less than n bytes are read, an error is returned
	if err != nil {
		return "", err
	}

	log.Info().Msg(base64.StdEncoding.EncodeToString(b))

	return base64.URLEncoding.EncodeToString(b), nil
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

// GetAuthMiddleware returns a gin middleware used for authentication. This middleware first tries bind the available
// data contained either in the body or as parameters into a data.Login struct, then tries to authenticate the user
// with an api key or an user/password if no key is provided.
func GetAuthMiddleware(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var l data.Login
		if err := c.ShouldBind(&l); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			log.Info().Err(err).Msg("logging failed")
			return
		}

		if l.ApiKey != "" {
			// If we receive an api key
			salt, ak, err := validateApiKey(l.ApiKey)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				log.Info().Err(err).Msg("authentication failed")
				return
			}

			u, err := ds.SelectUserLoginByApiKeySalt(string(salt))

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Error().Err(err).Msg("database error")
				return
			}

			err = bcrypt.CompareHashAndPassword(u.ApiKeyHash, ak)

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				log.Info().Err(err).Msg("authentication failed")
				return
			}

			// Keep track of the user if he successfully authenticated
			c.Keys["user"] = u
		} else if l.Username != "" {
			// If we receive a username+password
			u, err := ds.SelectUserLogin(l.Username)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Error().Err(err).Msg("database error")
				return
			}

			err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(l.Password))

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				log.Info().Err(err).Msg("authentication failed")
				return
			}

			// Keep track of the user if he successfully authenticated
			c.Keys["user"] = u
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
