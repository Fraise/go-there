package gopath

import (
	"github.com/gin-gonic/gin"
	"go-there/data"
)

type DataSourcer interface {
	SelectUser(username string) data.User
	InsertUser(user data.User) error
	DeleteUser(username string) error
}

func Init(e *gin.Engine, ds DataSourcer) {
	goPath := e.Group("/go")
	{
		goPath.GET("/:path", getPathHandler(ds))
	}
}
