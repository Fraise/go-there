package gopath

import "github.com/gin-gonic/gin"

func getPathHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": c.Param("path"),
		})
	}
}
