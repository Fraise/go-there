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
	goPath := e.Group("/go")

	if conf.Server.AuthRedirect {
		goPath.Use(auth.GetAuthMiddleware(ds))
	}

	goPath.GET("/:path", getPathHandler(ds))

}
