package api

import (
	"github.com/gin-gonic/gin"
	"go-there/auth"
	"go-there/config"
	"go-there/data"
)

// DataSourcer represents the datasource.DataSource methods needed by the api package to access the data.
type DataSourcer interface {
	SelectUser(username string) (data.User, error)
	SelectUserLogin(username string) (data.User, error)
	SelectApiKeyHashByUser(username string) ([]byte, error)
	SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error)
	InsertUser(user data.User) error
	DeleteUser(username string) error
	UpdatetUserPassword(user data.User) error
	UpdatetUserApiKey(user data.User) error
	InsertPath(path data.Path) error
	DeletePath(path data.Path) error
}

// Init initializes the API paths from the provided configuration and add them to the *gin.Engine.
func Init(conf *config.Configuration, e *gin.Engine, ds DataSourcer) {
	ep := conf.Endpoints["manage_users"]
	if ep.Enabled {
		// Init /api/users/:user routes
		api := e.Group("/api")

		api.GET("/users/:user", getUserHandler(ds))
		api.DELETE("/users/:user", getDeleteUserHandler(ds))
		api.PATCH("/users/:user", getUpdateUserHandler(ds))

		if ep.Auth {
			api.Use(auth.GetAuthMiddleware(ds))
		}

		api.Use(auth.GetPermissionsMiddleware(ep.AdminOnly))
	}

	ep = conf.Endpoints["create_users"]
	if ep.Enabled {
		// Init /api/users route
		userRoute := e.POST("/api/users", getCreateHandler(ds))

		if ep.Auth {
			userRoute.Use(auth.GetAuthMiddleware(ds))
		}

		userRoute.Use(auth.GetPermissionsMiddleware(ep.AdminOnly))
	}

	ep = conf.Endpoints["manage_paths"]
	if ep.Enabled {
		// Init /api/path route
		path := e.Group("/api/path")

		path.POST("", getPostPathHandler(ds))
		path.DELETE("", getDeletePathHandler(ds))

		if ep.Auth {
			path.Use(auth.GetAuthMiddleware(ds))
		}

		path.Use(auth.GetPermissionsMiddleware(ep.AdminOnly))
	}
}
