package gopath

import (
	"github.com/gin-gonic/gin"
	"go-there/auth"
	"go-there/config"
	"go-there/data"
	"go-there/logging"
)

// DataSourcer represents the database.DataSource methods needed by the gopath package to access the data.
type DataSourcer interface {
	SelectUserLogin(username string) (data.User, error)
	SelectUserLoginByApiKeyHash(apiKeyHash string) (data.User, error)
	GetTarget(path string) (string, error)
	GetAuthToken(token string) (data.AuthToken, error)
}

// Init initializes the redirect paths from the provided configuration and add them to the *gin.Engine.
func Init(conf *config.Configuration, e *gin.Engine, ds DataSourcer) {
	ep := conf.Endpoints["go"]
	if ep.Enabled {
		goPath := e.Group("/go")

		if ep.Log {
			goPath.Use(logging.GetLoggingMiddleware())
		}

		if ep.Auth {
			goPath.Use(auth.GetAuthMiddleware(ds))
		}

		goPath.GET("/:path", getPathHandler(ds))
	}
}
