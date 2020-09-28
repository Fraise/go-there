package api

import (
	"github.com/gin-gonic/gin"
)


func MakeApi(e *gin.Engine)  {
	api := e.Group("/api")
	{
		api.POST("/users", createHandler)
		api.GET("/users/:user", userHandler)
	}
}
