package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-there/config"
	"go-there/data"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"os"
)

var JwtSigningKey *rsa.PrivateKey

func InitJwtSigningKey(config *config.Configuration) {
	if config.Server.JwtSigningKeyPath == "" {
		log.Fatal().Msg("invalid JWT signing key")
	}

	if _, err := os.Stat(config.Server.JwtSigningKeyPath); os.IsNotExist(err) {
		log.Warn().Msg("JWT signing key, trying to generate one")

		JwtSigningKey, err = rsa.GenerateKey(rand.Reader, 4096)

		if err != nil {
			log.Fatal().Err(err).Msg("could not generate JWT signing key")
		}

		keyBytes, err := x509.MarshalPKCS8PrivateKey(JwtSigningKey)

		if err != nil {
			log.Fatal().Err(err).Msg("could not marshal JWT signing key to disk")
		}

		err = ioutil.WriteFile(config.Server.JwtSigningKeyPath, keyBytes, 0700)

		if err != nil {
			log.Fatal().Err(err).Msg("could not marshal JWT signing key to disk")
		}

		log.Info().Msg("successfully generated JWT signing key")
	} else {
		keyString, err := ioutil.ReadFile(config.Server.JwtSigningKeyPath)

		if err != nil {
			log.Fatal().Err(err).Msg("error reading JWT signing key from file")
		}

		block, _ := pem.Decode(keyString)
		parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)

		if err != nil {
			log.Fatal().Err(err).Msg("error parsing JWT signing key")
		}

		var ok bool
		JwtSigningKey, ok = parseResult.(*rsa.PrivateKey)

		if !ok {
			log.Fatal().Err(err).Msg("error JWT signing key type, must be RSA")
		}
	}

}

// GetAuthMiddleware returns a gin middleware used for authentication. This middleware first tries to bind either a
// X-Api-Key header in a data.HeaderLogin struct or the data contained either in the body or as parameters into a
// data.Login struct. It then tries to authenticate the user with an api key or an user/password if no key is provided.
func GetAuthMiddleware(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var hl data.HeaderLogin

		// Tries to bind authentication headers first
		if err := c.ShouldBindHeader(&hl); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Check API key
		if hl.XApiKey != "" {
			// If the header contains an API key, do not bind the other fields
			hash, ak, err := validateApiKey(hl.XApiKey)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			u, err := ds.SelectUserLoginByApiKeyHash(string(hash))

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}

			if u.Username == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			err = bcrypt.CompareHashAndPassword(u.ApiKeyHash, ak)

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.Keys = make(map[string]interface{})

			// Keep track of the user if he successfully authenticated
			c.Keys["user"] = u
			c.Keys["logUser"] = u.Username
			// Keep track of which user data we want to access
			c.Keys["reqUser"] = c.Param("user")

			return
		}

		// Check basic auth and bearer token
		if hl.Authorization != "" {
			ld, err := authHeaderToLoginData(hl.Authorization)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			u := data.User{}

			if ld.DataType == data.Basic {
				u, err = ds.SelectUserLogin(ld.BasicAuthLogin.Username)

				if err != nil {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}

				err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(ld.BasicAuthLogin.Password))

				if err != nil {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
			} else {
				if ld.IsExpired() {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}

				// We still check if the user has not been deleted before his token expired
				u, err = ds.SelectUserLogin(ld.BasicAuthLogin.Username)

				if err != nil {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
			}

			c.Keys = make(map[string]interface{})

			// Keep track of the user if he successfully authenticated
			c.Keys["user"] = u
			c.Keys["logUser"] = u.Username
			// Keep track of which user data we want to access
			c.Keys["reqUser"] = c.Param("user")

			return
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

// GetPermissionsMiddleware verify that the logged used has the permission to access the requested resource. A user
// can only access his profile, and admin can access any profile.
func GetPermissionsMiddleware(adminOnly bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		loggedUser := GetLoggedUser(c)

		// If the user is admin, always allow access
		if loggedUser.IsAdmin {
			return
		}

		// If admin rights are required
		if adminOnly && !loggedUser.IsAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// If an user is logged, make sure he can only see his data if he's not admin
		reqUser := GetRequestedUser(c)

		// If the resource "belong" to no one. In this case, the request always access the client's own resources
		if reqUser == "" {
			return
		}

		if loggedUser.Username != reqUser {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
