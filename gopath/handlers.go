package gopath

import "github.com/gin-gonic/gin"

func getPathHandler(ds DataSourcer) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": ctx.Param("path"),
		})
	}
}
