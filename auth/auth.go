package auth

import (
	"bytes"
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

// validateApiKey takes an api key and return (salt, apikey, error).
func validateApiKey(apiKey string) ([]byte, []byte, error) {
	apiKeyArr := bytes.Split([]byte(apiKey), []byte("."))

	if len(apiKeyArr) != 2 {
		return nil, nil, data.ErrInvalidKey
	}

	return apiKeyArr[0], apiKeyArr[1], nil
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

			akHash, err := ds.SelectApiKeyHashBySalt(string(salt))

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Error().Err(err).Msg("database error")
				return
			}

			err = bcrypt.CompareHashAndPassword(akHash, ak)

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				log.Info().Err(err).Msg("authentication failed")
				return
			}
		} else if l.Username != "" {
			// If we receive a username+password
			u, err := ds.SelectUser(l.Username)

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
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
