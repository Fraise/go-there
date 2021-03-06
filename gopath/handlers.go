package gopath

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-there/data"
	"net/http"
)

// getPathHandler returns the redirection handler. If no redirection exists, then http.StatusNotFound is returned.
func getPathHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		t, err := ds.GetTarget(c.Param("path"))

		if err != nil {
			switch {
			case errors.Is(err, data.ErrSqlNoRow):
				c.Status(http.StatusNotFound)
				return
			default:
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}
		}

		c.Redirect(http.StatusFound, t)
	}
}
