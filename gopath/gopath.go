package gopath

import (
	"github.com/gin-gonic/gin"
	"go-there/auth"
	"go-there/config"
	"go-there/data"
)

// DataSourcer represents the datasource.DataSource methods needed by the gopath package to access the data.
type DataSourcer interface {
	SelectUserLogin(username string) (data.User, error)
	SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error)
	GetTarget(path string) (string, error)
}

// Init initializes the redirect paths from the provided configuration and add them to the *gin.Engine.
func Init(conf *config.Configuration, e *gin.Engine, ds DataSourcer) {
	ep := conf.Endpoints["go"]
	if ep.Enabled {
		goPath := e.Group("/go")

		if conf.Endpoints["go"].Auth {
			goPath.Use(auth.GetAuthMiddleware(ds))
		}

		goPath.GET("/:path", getPathHandler(ds))
	}
}
