package auth

import (
	"github.com/gin-gonic/gin"
	"go-there/data"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// GetAuthMiddleware returns a gin middleware used for authentication. This middleware first tries to bind either a
// X-Api-Key header in a data.HeaderLogin struct or the data contained either in the body or as parameters into a
// data.Login struct. It then tries to authenticate the user with an api key or an user/password if no key is provided.
func GetAuthMiddleware(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var l data.Login
		var hl data.HeaderLogin

		// Tries to bind authentication header first
		if err := c.ShouldBindHeader(&hl); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if hl.XApiKey != "" {
			// If the header contains an API key, do not bind the other fields
			l.ApiKey = hl.XApiKey
		} else {
			// Tries to bind the data depending on the content type, then tries to bind the form data
			if err := c.ShouldBind(&l); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		c.Keys = make(map[string]interface{})

		c.Keys["logUser"] = l.Username

		if l.ApiKey != "" {
			// If we receive an api key
			salt, ak, err := validateApiKey(l.ApiKey)

			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			u, err := ds.SelectUserLoginByApiKeySalt(string(salt))

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
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

		// If admin rights are required
		if adminOnly && !loggedUser.IsAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// If the user is admin, always allow access
		if loggedUser.IsAdmin {
			return
		}

		// If no login is required continue, as it is already validated by the auth middleware
		if loggedUser.Username == "" {
			return
		}

		// If an user is logged, make sure he can only see his data if he's not admin
		reqUser := GetRequestedUser(c)

		if reqUser == "" {
			return
		}

		if loggedUser.Username != reqUser {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
