package gopath

import (
	"github.com/gin-gonic/gin"
	"go-there/auth"
	"go-there/config"
	"go-there/data"
)

type DataSourcer interface {
	SelectUser(username string) (data.User, error)
	SelectUserPassword(username string) ([]byte, error)
	SelectUserApiKey(username string) ([]byte, error)
	SelectApiKey(apiKey string) ([]byte, error)
	InsertUser(user data.User) error
	DeleteUser(username string) error
}

func Init(conf *config.Configuration, e *gin.Engine, ds DataSourcer) {
	goPath := e.Group("/go")

	if conf.Server.AuthApi {
		goPath.Use(auth.GetAuthMiddleware(ds))
	}

	goPath.GET("/:path", getPathHandler(ds))

}
