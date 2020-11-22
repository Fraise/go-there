package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-there/data"
)

// GetLoggingMiddleware returns a logging middleware which first parse the user info, call Next() then logs the request.
func GetLoggingMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		logInfo := data.LogInfo{
			Method:   c.Request.Method,
			Endpoint: c.Request.URL.Path,
			Ip:       c.ClientIP(),
		}

		c.Next()

		li, ok := c.Keys["logUser"]

		if ok {
			logInfo.User = li.(string)
		}

		logInfo.HttpCode = c.Writer.Status()

		ginErr := c.Errors.Last()

		if ginErr != nil {
			log.Info().
				Str("method", logInfo.Method).
				Int("http_code", logInfo.HttpCode).
				Str("endpoint", logInfo.Endpoint).
				Str("ip", logInfo.Ip).
				Str("user", logInfo.User).
				Err(ginErr.Err).
				Send()
		} else {
			log.Info().
				Str("method", logInfo.Method).
				Int("http_code", logInfo.HttpCode).
				Str("endpoint", logInfo.Endpoint).
				Str("ip", logInfo.Ip).
				Str("user", logInfo.User).
				Send()
		}
	}
}
