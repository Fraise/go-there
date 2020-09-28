package main

import (
	"github.com/gin-gonic/gin"
	"go-there/gopath"
	"net/http"
)

func main() {
	e := gin.New()

	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	gopath.Init(e)

	_ = http.ListenAndServe("0.0.0.0:8080", e)
}
