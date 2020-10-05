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
	InsertPath(path data.Path) error
	DeletePath(path data.Path) error
}

func Init(conf *config.Configuration, e *gin.Engine, ds DataSourcer) {
	// Init /api/users/:user routes
	api := e.Group("/api")

	api.GET("/users/:user", getUserHandler(ds))
	api.DELETE("/users/:user", getDeleteUserHandler(ds))
	api.PATCH("/users/:user", getUpdateUserHandler(ds))

	if conf.Endpoints["manage_users"].NeedAuth {
		api.Use(auth.GetAuthMiddleware(ds))
	}

	if conf.Endpoints["manage_users"].NeedAdmin {
		api.Use(auth.GetPermissionsMiddleware())
	}

	// Init /api/users route
	userRoute := e.POST("/api/users", getCreateHandler(ds))

	if conf.Endpoints["create_users"].NeedAuth {
		userRoute.Use(auth.GetAuthMiddleware(ds))
	}

	if conf.Endpoints["create_users"].NeedAdmin {
		userRoute.Use(auth.GetPermissionsMiddleware())
	}

	// Init /api/path route
	path := e.Group("/api/path")

	path.POST("/", getPostPathHandler(ds))
	path.DELETE("/", getDeletePathHandler(ds))

	if conf.Endpoints["manage_paths"].NeedAuth {
		path.Use(auth.GetAuthMiddleware(ds))
	}

	if conf.Endpoints["manage_paths"].NeedAdmin {
		path.Use(auth.GetPermissionsMiddleware())
	}
}
