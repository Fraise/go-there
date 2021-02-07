package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
		var l data.Login
		var hl data.HeaderLogin
		var t data.AuthToken

		// Tries to bind authentication header first
		if err := c.ShouldBindHeader(&hl); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if hl.XAuthToken != "" {
			var err error
			// If the header contains a session token, do not bind the other fields
			t, err = ds.GetAuthToken(hl.XAuthToken)

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

			if t.ExpirationTS < time.Now().Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, data.ErrorResponse{Error: "token expired"})
			} else if t.ExpirationTS < time.Now().Unix()+604800 {
				// If the expiration date if in less than 1 week, renew it
				t.ExpirationTS = time.Now().Unix()

				err := ds.UpdateAuthToken(t)

				// This can fail silently for the user, but it still needs to be reported to the system
				if err != nil {
					_ = c.Error(err)
				}
			}

			u, err := ds.SelectUserLogin(t.Username)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}

			// The keys map is only initialized if a call to ShouldBindBody is made
			c.Keys = make(map[string]interface{})

			// Keep track of the user if he successfully authenticated
			c.Keys["user"] = u
			c.Keys["logUser"] = t.Username
			// Keep track of which user data we want to access
			c.Keys["reqUser"] = c.Param("user")

			return
		}

		if hl.XApiKey != "" {
			// If the header contains an API key, do not bind the other fields
			l.ApiKey = hl.XApiKey
			// The keys map is only initialized if a call to ShouldBindBody is made
			c.Keys = make(map[string]interface{})
		} else {
			// Tries to bind the JSON data related to login
			// Implicitly initialize the c.Keys map
			if err := c.ShouldBindBodyWith(&l, binding.JSON); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		c.Keys["logUser"] = l.Username

		if l.ApiKey != "" {
			// If we receive an api key
			hash, ak, err := validateApiKey(l.ApiKey)

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

			// Keep track of the user if he successfully authenticated
			c.Keys["user"] = u
			c.Keys["logUser"] = u.Username
			// Keep track of which user data we want to access
			c.Keys["reqUser"] = c.Param("user")
		} else if l.Username != "" {
			// If we receive a username+password
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

			// Keep track of the user if he successfully authenticated
			c.Keys["user"] = u
			c.Keys["logUser"] = u.Username
			// Keep track of which user data we want to access
			c.Keys["reqUser"] = c.Param("user")
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
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
