package api

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
}

func Init(conf *config.Configuration, e *gin.Engine, ds DataSourcer) {
	api := e.Group("/api")

	if conf.Server.AuthApi {
		api.Use(auth.GetAuthMiddleware(ds))
	}

	api.GET("/users/:user", getUserHandler(ds))

	if conf.Server.AuthCreateUser {
		api.POST("/api/users", getCreateHandler(ds)).Use(auth.GetAuthMiddleware(ds))
	} else {
		api.POST("/api/users", getCreateHandler(ds))
	}

}
