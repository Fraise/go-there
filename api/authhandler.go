package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go-there/auth"
	"go-there/data"
	"net/http"
	"time"
)

func getGetJwtHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		u := auth.GetLoggedUser(c)

		if u.Username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		t := jwt.New()
		if err := t.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*24*7).Unix()); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(fmt.Errorf("error creating a JWT: %w", err))
			return
		}
		if err := t.Set("username", u.Username); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(fmt.Errorf("error creating a JWT: %w", err))
			return
		}
		if err := t.Set("is_admin", u.IsAdmin); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(fmt.Errorf("error creating a JWT: %w", err))
			return
		}

		var err error

		jwtBytes, err := jwt.Sign(t, jwa.RS256, auth.JwtSigningKey)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(fmt.Errorf("error signing a JWT: %w", err))
			return
		}

		c.JSON(http.StatusOK, data.JwtResponse{Jwt: string(jwtBytes)})
	}
}
