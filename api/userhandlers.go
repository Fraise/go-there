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

		// We don't need to store the salt section individually for a password
		hash, _, err := auth.GetHashFromPassword(cu.CreatePassword)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		// Generate a random API key
		apiKey, err := auth.GenerateRandomB64String(16)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		// Get its corresponding hash and salt
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
				// TODO clean that up
				ApiKey: base64.URLEncoding.EncodeToString(append(apiKeySalt, []byte(":"+apiKey)...)),
			})
	}
}

// getUserHandler returns a gin handler which select an user in the datasource or return http.StatusNotFound if the user
// does not exist
func getUserHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		u, err := ds.SelectUser(c.Param("user"))

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		if u.Username == "" {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, u)
	}
}

// getUserHandler returns a gin handler which delete an user in the datasource.
func getDeleteUserHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := ds.DeleteUser(c.Param("user"))

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusOK)
	}
}

// getUpdateUserHandler returns a gin handler which updates an user in the datasource from the request body.
func getUpdateUserHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		pu := data.PatchUser{}

		err := c.ShouldBindBodyWith(&pu, binding.JSON)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		u := data.User{Username: c.Param("user")}
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

			err = ds.UpdateUserPassword(u)

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

			err = ds.UpdateUserApiKey(u)

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

// getUserList returns a gin handler which fetch the list of all users in the datasource.
func getUserList(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		users, err := ds.SelectAllUsers()

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, users)
	}
}
