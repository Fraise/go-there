package api

import (
	"github.com/gin-gonic/gin"
	"go-there/data"
)

type DataSourcer interface {
	SelectUser(username string) (data.User, error)
	SelectUserPassword(username string) (data.Password, error)
	SelectUserApiKey(username string) (data.ApiKey, error)
	InsertUser(user data.User) error
	DeleteUser(username string) error
}

func Init(e *gin.Engine, ds DataSourcer) {
	api := e.Group("/api")
	{
		api.POST("/users", getCreateHandler(ds))
		api.GET("/users/:user", getUserHandler(ds))
	}
}
