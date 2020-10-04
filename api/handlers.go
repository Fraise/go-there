package api

import (
	"github.com/gin-gonic/gin"
	"go-there/auth"
	"go-there/data"
	"net/http"
)

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

		apkiKey, err := auth.GenerateRandomString(16)

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
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, data.CreateUserResponse{ApiKey: string(apiKeySalt) + "." + apkiKey})
	}
}

func getUserHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		u, _ := ds.SelectUser(c.Param("user"))
		c.JSON(http.StatusOK, u)
	}
}
