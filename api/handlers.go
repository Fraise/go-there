package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
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
		err := c.ShouldBindJSON(&cu)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// We don't need to store the salt for a password
		hash, _, err := auth.GetHashFromPassword(cu.CreatePassword)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		apkiKey, err := auth.GenerateRandomB64String(16)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		apkiKeyHash, apiKeySalt, err := auth.GetHashFromPassword(apkiKey)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		u := data.User{
			Username:     cu.CreateUser,
			IsAdmin:      false,
			PasswordHash: hash,
			ApiKeySalt:   apiKeySalt,
			ApiKeyHash:   apkiKeyHash,
		}

		err = ds.InsertUser(u)

		if err != nil {
			if err == data.ErrSqlDuplicateRow {
				c.AbortWithStatusJSON(http.StatusBadRequest, data.ErrorResponse{Error: "user already exists"})
				return
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		c.JSON(
			http.StatusOK,
			data.ApiKeyResponse{
				ApiKey: base64.URLEncoding.EncodeToString(apiKeySalt) + "." + apkiKey,
			})
	}
}

func getUserHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		u, _ := ds.SelectUser(c.Param("user"))
		c.JSON(http.StatusOK, u)
	}
}
