package main

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestApi(t *testing.T) {
	e := httpexpect.New(t, "http://go-there:8080")

	e.GET("/health").Expect().Status(http.StatusOK)
}
