package health

import (
	"github.com/gin-gonic/gin"
	"go-there/config"
	"net/http"
)

// Init initializes the API paths from the provided configuration and add them to the *gin.Engine.
func Init(conf *config.Configuration, e *gin.Engine) {
	ep := conf.Endpoints["health"]
	if ep.Enabled {
		e.GET("/health", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
	}
}
