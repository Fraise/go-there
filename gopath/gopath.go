package gopath

import (
	"github.com/gin-gonic/gin"
	"go-there/auth"
	"go-there/config"
	"go-there/data"
)

type DataSourcer interface {
	SelectUser(username string) (data.User, error)
	SelectUserLogin(username string) (data.User, error)
	SelectApiKeyHashByUser(username string) ([]byte, error)
	SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error)
	SelectApiKeyHashBySalt(apiKeySalt string) ([]byte, error)
	InsertUser(user data.User) error
	DeleteUser(username string) error
	GetTarget(path string) (string, error)
}

func Init(conf *config.Configuration, e *gin.Engine, ds DataSourcer) {
	if !conf.Endpoints["/go"].Enabled {
		return
	}

	goPath := e.Group("/go")

	goPath.GET("/:path", getPathHandler(ds))

	if conf.Endpoints["go"].NeedAuth {
		goPath.Use(auth.GetAuthMiddleware(ds))
	}
}
