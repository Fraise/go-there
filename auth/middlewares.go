package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go-there/data"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

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

		// Check token first
		if hl.XAuthToken != "" {
			var t data.AuthToken
			var err error

			// X-Auth-Token is in b64, so we need to decode it first
			decodedXAuthToken, err := base64.StdEncoding.DecodeString(hl.XAuthToken)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(decodedXAuthToken, &t)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			// If the header contains a session token, do not bind the other fields
			dbToken, err := ds.GetAuthToken(t.Token)

			if err != nil {
				switch {
				case errors.Is(err, data.ErrSqlNoRow):
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				default:
					c.AbortWithStatus(http.StatusInternalServerError)
					_ = c.Error(err)
					return
				}
			}

			if dbToken.Username != t.Username {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if dbToken.ExpirationTS < time.Now().Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, data.ErrorResponse{Error: "token expired"})
			} else if dbToken.ExpirationTS < time.Now().Unix()+604800 {
				// If the expiration date if in less than 1 week, renew it
				dbToken.ExpirationTS = time.Now().Unix()

				err := ds.UpdateAuthToken(dbToken)

				// This can fail silently for the user, but it still needs to be reported to the system
				if err != nil {
					_ = c.Error(err)
				}
			}

			u, err := ds.SelectUserLogin(dbToken.Username)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}

			c.Keys = make(map[string]interface{})

			// Keep track of the user if he successfully authenticated
			c.Keys["user"] = u
			c.Keys["logUser"] = dbToken.Username
			// Keep track of which user data we want to access
			c.Keys["reqUser"] = c.Param("user")

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

		// Check basic auth
		if hl.Authorization != "" {
			l, err := basicAuthToLogin(hl.Authorization)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			u, err := ds.SelectUserLogin(l.Username)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}

			if u.Username == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(l.Password))

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
