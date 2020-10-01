package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-there/data"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type DataSourcer interface {
	SelectUser(username string) (data.User, error)
	SelectUserPassword(username string) ([]byte, error)
	SelectUserApiKey(username string) ([]byte, error)
	SelectApiKey(apiKey string) ([]byte, error)
	InsertUser(user data.User) error
	DeleteUser(username string) error
}

func GetAuthMiddleware(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var l data.Login
		if err := c.ShouldBind(&l); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			log.Info().Err(err).Msg("logging failed")
			return
		}

		if l.Username != "" {
			u, err := ds.SelectUser(l.Username)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Error().Err(err).Msg("database error")
				return
			}

			err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(l.Password))

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				log.Info().Err(err).Msg("authentication failed")
				return
			}
		} else if l.ApiKey != "" {
			ak, err := ds.SelectApiKey(l.ApiKey)

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Error().Err(err).Msg("database error")
				return
			}

			err = bcrypt.CompareHashAndPassword(ak, []byte(l.ApiKey))

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				log.Info().Err(err).Msg("authentication failed")
				return
			}
		}
	}
}
