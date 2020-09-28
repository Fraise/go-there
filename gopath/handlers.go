package gopath

import "github.com/gin-gonic/gin"

func pathHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": ctx.Param("path"),
	})
}
