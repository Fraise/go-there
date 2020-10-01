package api

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
	api := e.Group("/api")

	if conf.Server.AuthApi {
		api.Use(auth.GetAuthMiddleware(ds))
	}

	api.POST("/users", getCreateHandler(ds))
	api.GET("/users/:user", getUserHandler(ds))

}
