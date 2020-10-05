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
	UpdatetUserPassword(user data.User) error
	UpdatetUserApiKey(user data.User) error
}

func Init(conf *config.Configuration, e *gin.Engine, ds DataSourcer) {
	api := e.Group("/api")

	if conf.Server.AuthApi {
		api.Use(auth.GetAuthMiddleware(ds))
	}

	api.GET("/users/:user", getUserHandler(ds))
	api.DELETE("/users/:user", getDeleteUserHandler(ds))
	api.PATCH("/users/:user", getUpdateUserHandler(ds))

	if conf.Server.AuthCreateUser {
		e.POST("/api/users", getCreateHandler(ds)).Use(auth.GetAuthMiddleware(ds))
	} else {
		e.POST("/api/users", getCreateHandler(ds))
	}

}
