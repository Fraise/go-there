package gopath

import "github.com/gin-gonic/gin"

func Init(e *gin.Engine)  {
	goPath := e.Group("/go")
	{
		goPath.GET("/:path", pathHandler)
	}
}

