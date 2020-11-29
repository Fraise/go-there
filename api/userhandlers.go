package api

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-there/auth"
	"go-there/data"
	"net/http"
)

// getCreateHandler returns a gin handler which tries to insert a new user in the database. It first bind provided JSON
// data (or fails), then hashes the password, generates an API key and tries to insert everything in the database. If it
// succeeds, an API key is returned to the user, if the new user already exists, it returns 400 and "user already
// exists" in the response body.
func getCreateHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		cu := data.CreateUser{}
		err := c.ShouldBindBodyWith(&cu, binding.JSON)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// We don't need to store the salt for a password
		hash, _, err := auth.GetHashFromPassword(cu.CreatePassword)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		apiKey, err := auth.GenerateRandomB64String(16)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		apiKeyHash, apiKeySalt, err := auth.GetHashFromPassword(apiKey)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		u := data.User{
			Username:     cu.CreateUser,
			IsAdmin:      false,
			PasswordHash: hash,
			ApiKeySalt:   apiKeySalt,
			ApiKeyHash:   apiKeyHash,
		}

		err = ds.InsertUser(u)

		if err != nil {
			switch {
			case errors.Is(err, data.ErrSqlDuplicateRow):
				// TODO handle duplicate salt
				c.AbortWithStatusJSON(http.StatusBadRequest, data.ErrorResponse{Error: "user already exists"})
				return
			default:
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}
		}

		c.JSON(
			http.StatusOK,
			data.ApiKeyResponse{
				ApiKey: base64.URLEncoding.EncodeToString(apiKeySalt) + "." + apiKey,
			})
	}
}

func getUserHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		requestedUser := c.Param("user")

		loggedUser := auth.GetLoggedUser(c)

		// If an user is logged, make sure he can only see his data if he's not admin
		if loggedUser.Username != "" && loggedUser.Username != requestedUser && !loggedUser.IsAdmin {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		u, err := ds.SelectUser(requestedUser)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		// Clean the data before responding to the request
		c.JSON(http.StatusOK, data.User{Username: u.Username, IsAdmin: u.IsAdmin})
	}
}

func getDeleteUserHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		requestedUser := c.Param("user")

		loggedUser := auth.GetLoggedUser(c)

		// If an user is logged, make sure he can only see his data if he's not admin
		if loggedUser.Username != "" && loggedUser.Username != requestedUser && !loggedUser.IsAdmin {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err := ds.DeleteUser(requestedUser)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func getUpdateUserHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		requestedUser := c.Param("user")

		loggedUser := auth.GetLoggedUser(c)

		// If an user is logged, make sure he can only see his data if he's not admin
		if loggedUser.Username != "" && loggedUser.Username != requestedUser && !loggedUser.IsAdmin {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		pu := data.PatchUser{}

		err := c.ShouldBindBodyWith(&pu, binding.JSON)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		u := data.User{Username: requestedUser}
		ar := data.ApiKeyResponse{}

		if pu.PatchPassword != "" {
			// We don't need to store the salt for a password
			hash, _, err := auth.GetHashFromPassword(pu.PatchPassword)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}

			u.PasswordHash = hash

			err = ds.UpdatetUserPassword(u)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}
		}

		if pu.PatchApiKey {
			apiKey, err := auth.GenerateRandomB64String(16)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}

			apiKeyHash, apiKeySalt, err := auth.GetHashFromPassword(apiKey)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}

			u.ApiKeyHash = apiKeyHash
			u.ApiKeySalt = apiKeySalt

			err = ds.UpdatetUserApiKey(u)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}

			ar.ApiKey = apiKey
		}

		c.JSON(http.StatusOK, ar)
	}
}
