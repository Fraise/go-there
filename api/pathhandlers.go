package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-there/auth"
	"go-there/data"
	"net/http"
)

// getPostPathHandler returns a gin handler for POST requests when creating a new redirect. Returns
// http.StatusBadRequest if it cannot bind the required JSON data for path creation or if the path already exists.
func getPostPathHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		u := auth.GetLoggedUser(c)

		cp := data.CreatePath{}

		err := c.ShouldBindBodyWith(&cp, binding.JSON)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		p := data.Path{
			Path:   cp.Path,
			Target: cp.Target,
			UserId: u.Id,
		}

		err = ds.InsertPath(p)

		if err != nil {
			switch {
			case errors.Is(err, data.ErrSqlDuplicateRow):
				c.AbortWithStatusJSON(http.StatusBadRequest, data.ErrorResponse{Error: "path already exists"})
				return
			default:
				c.AbortWithStatus(http.StatusInternalServerError)
				_ = c.Error(err)
				return
			}
		}

		c.Status(http.StatusOK)
	}
}

func getDeletePathHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		u := auth.GetLoggedUser(c)

		dp := data.DeletePath{}

		err := c.ShouldBindBodyWith(&dp, binding.JSON)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		p := data.Path{
			Path:   dp.Path,
			UserId: u.Id,
		}

		err = ds.DeletePath(p)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = c.Error(err)
		}

		c.Status(http.StatusOK)
	}
}
