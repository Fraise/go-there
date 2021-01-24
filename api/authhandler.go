package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-there/auth"
	"go-there/data"
	"net/http"
	"time"
)

// getAuthTokenHandler returns a gin handler for GET requests for a session token. Returns http.StatusBadRequest if the
// user is anonymous.
func getAuthTokenHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		u := auth.GetLoggedUser(c)

		if u.Username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token, err := ds.GetAuthTokenByUser(u.Username)

		if err != nil {
			switch {
			case errors.Is(err, data.ErrSql):
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}
		}

		// If the token does not exist, generate one
		if token.Token == "" {
			token.Token, err = auth.GenerateRandomB64String(32)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}

			// TODO make it configurable
			token.ExpirationTS = time.Now().Unix() + 30*24*3600 // 30 days
			token.Username = u.Username

			err = ds.InsertAuthToken(token)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}
		}

		c.JSON(http.StatusOK, token)
	}
}

// getDeleteAuthTokenHandler returns a gin handler for GET requests for a session token. Returns http.StatusBadRequest if the
// user is anonymous.
func getDeleteAuthTokenHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		u := auth.GetLoggedUser(c)

		if u.Username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err := ds.DeleteAuthToken(data.AuthToken{
			Username: u.Username,
		})

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusOK)
	}
}
