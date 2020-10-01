package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getCreateHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

func getUserHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		u, _ := ds.SelectUser(c.Param("user"))
		c.JSON(http.StatusOK, u)
	}
}
