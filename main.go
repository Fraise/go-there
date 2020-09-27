package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/create")
		api.GET("/:user")
	}

	goPath := r.Group("/go")
	{
		goPath.GET("/:path", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": c.Param("path"),
			})
		})
	}

	r.Run()
}