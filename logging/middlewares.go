package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-there/data"
)

// GetLoggingMiddleware returns a logging middleware which first parse the user info, call Next() then logs the request
func GetLoggingMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		logInfo := data.LogInfo{
			Method:   c.Request.Method,
			Endpoint: c.Request.URL.Path,
			Ip:       c.ClientIP(),
		}

		c.Next()

		li, ok := c.Keys["logInfo"]

		if ok {
			logInfo.Login = li.(data.Login)
		}

		logInfo.HttpCode = c.Writer.Status()

		log.Info().
			Str("method", logInfo.Method).
			Int("http_code", logInfo.HttpCode).
			Str("endpoint", logInfo.Endpoint).
			Str("ip", logInfo.Ip).
			Interface("login", logInfo.Login).
			Msg("")
	}
}
