package gopath

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-there/data"
	"net/http"
)

func getPathHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		t, err := ds.GetTarget(c.Param("path"))

		if err != nil {
			switch err {
			case data.ErrSqlNoRow:
				c.Status(http.StatusNotFound)
				return
			default:
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Error().Err(err).Msg("database error")
				return
			}
		}

		c.Redirect(http.StatusFound, t)
	}
}
