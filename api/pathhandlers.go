package api

import (
	"github.com/gin-gonic/gin"
	"go-there/auth"
	"go-there/data"
	"net/http"
)

func getPostPathHandler(ds DataSourcer) func(c *gin.Context) {
	return func(c *gin.Context) {
		u := auth.GetLoggedUser(c)

		cp := data.CreatePath{}

		err := c.ShouldBindJSON(&cp)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		p := data.Path{
			Path:   cp.Path,
			Target: cp.Target,
			User:   u.Username,
		}

		err = ds.InsertPath(p)

		if err != nil {
			if err == data.ErrSqlDuplicateRow {
				c.AbortWithStatusJSON(http.StatusBadRequest, data.ErrorResponse{Error: "path already exists"})
				return
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
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

		err := c.ShouldBindJSON(&dp)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		p := data.Path{
			Path: dp.Path,
			User: u.Username,
		}

		err = ds.DeletePath(p)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.Status(http.StatusOK)
	}
}
