package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-there/auth"
	"go-there/config"
	"go-there/data"
	"go-there/logging"
	"regexp"
)

// DataSourcer represents the database.DataSource methods needed by the api package to access the data.
type DataSourcer interface {
	SelectUser(username string) (data.UserInfo, error)
	SelectAllUsers() ([]data.UserInfo, error)
	SelectUserLogin(username string) (data.User, error)
	SelectApiKeyHashByUser(username string) ([]byte, error)
	SelectUserLoginByApiKeyHash(apiKeyHash string) (data.User, error)
	InsertUser(user data.User) error
	DeleteUser(username string) error
	UpdateUserPassword(user data.User) error
	UpdateUserApiKey(user data.User) error
	InsertPath(path data.Path) error
	DeletePath(path data.Path) error
}

// Init initializes the API paths from the provided configuration and add them to the *gin.Engine.
func Init(conf *config.Configuration, e *gin.Engine, ds DataSourcer) {
	ep := conf.Endpoints["manage_users"]
	if ep.Enabled {
		// Init /api/users/:user routes
		api := e.Group("/api")

		if ep.Log {
			api.Use(logging.GetLoggingMiddleware())
		}

		if ep.Auth {
			api.Use(auth.GetAuthMiddleware(ds))
			api.Use(auth.GetPermissionsMiddleware(ep.AdminOnly))
		}

		api.GET("/users/:user", getUserHandler(ds))
		api.DELETE("/users/:user", getDeleteUserHandler(ds))
		api.PATCH("/users/:user", getUpdateUserHandler(ds))
	}

	ep = conf.Endpoints["create_users"]
	if ep.Enabled {
		// Init /api/users route, POST new user
		userRoute := e.Group("/api/users")

		if ep.Log {
			userRoute.Use(logging.GetLoggingMiddleware())
		}

		if ep.Auth {
			userRoute.Use(auth.GetAuthMiddleware(ds))
			userRoute.Use(auth.GetPermissionsMiddleware(ep.AdminOnly))
		}

		userRoute.POST("", getCreateHandler(ds))
	}

	ep = conf.Endpoints["get_user_list"]
	if ep.Enabled {
		// Init /api/users route, GET all users
		userRoute := e.Group("/api/users")

		if ep.Log {
			userRoute.Use(logging.GetLoggingMiddleware())
		}

		if ep.Auth {
			userRoute.Use(auth.GetAuthMiddleware(ds))
			userRoute.Use(auth.GetPermissionsMiddleware(ep.AdminOnly))
		}

		userRoute.GET("", getUserList(ds))
	}

	ep = conf.Endpoints["manage_paths"]
	if ep.Enabled {
		// Init /api/path route
		path := e.Group("/api/path")

		if ep.Log {
			path.Use(logging.GetLoggingMiddleware())
		}

		if ep.Auth {
			path.Use(auth.GetAuthMiddleware(ds))
			path.Use(auth.GetPermissionsMiddleware(ep.AdminOnly))
		}

		path.POST("", getPostPathHandler(ds))
		path.DELETE("", getDeletePathHandler(ds))
	}

	ep = conf.Endpoints["jwt_token"]
	if ep.Enabled {
		// Init /api/auth route
		path := e.Group("/api/auth")

		if ep.Log {
			path.Use(logging.GetLoggingMiddleware())
		}

		if ep.Auth {
			path.Use(auth.GetAuthMiddleware(ds))
			path.Use(auth.GetPermissionsMiddleware(ep.AdminOnly))
		}

		path.GET("", getGetJwtHandler())
	}
}

// ApplyUserSettings parses the username and password rules, then apply them to the global variables. If a rule has a
// zero value, the default rules are used. Returns an data.ErrSettings if an error happens.
func ApplyUserSettings(conf *config.Configuration) error {
	if conf.UserRules.UsernameRegex != "" {
		r, err := regexp.Compile(conf.UserRules.UsernameRegex)

		if err != nil {
			return fmt.Errorf("%w : %s", data.ErrSettings, err)
		}

		usernameRegexp = r
	}

	if conf.UserRules.PasswordRegex != "" {
		r, err := regexp.Compile(conf.UserRules.PasswordRegex)

		if err != nil {
			return fmt.Errorf("%w : %s", data.ErrSettings, err)
		}

		passwordRegexp = r
	}

	// Set username and password size limits
	if passwordMinLen > passwordMaxLen || usernameMinLen > usernameMaxLen {
		return fmt.Errorf("%w : %s", data.ErrSettings, "invalid username or password validation")
	}

	if conf.UserRules.PasswordMinLen != 0 {
		passwordMinLen = conf.UserRules.PasswordMinLen
	}

	if conf.UserRules.PasswordMaxLen != 0 {
		passwordMaxLen = conf.UserRules.PasswordMaxLen
	}

	if conf.UserRules.UsernameMinLen != 0 {
		usernameMinLen = conf.UserRules.UsernameMinLen
	}

	if conf.UserRules.UsernameMaxLen != 0 {
		usernameMaxLen = conf.UserRules.UsernameMaxLen
	}

	return nil
}
