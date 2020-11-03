package main

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestHealth(t *testing.T) {
	e := httpexpect.New(t, "http://go-there:8080")

	e.GET("/health").Expect().Status(http.StatusOK)
}

func TestCreateUser(t *testing.T) {
	type CreateUser struct {
		CreateUser     string `json:"create_user"`
		CreatePassword string `json:"create_password"`
	}

	cu := CreateUser{
		CreateUser:     "alice",
		CreatePassword: "superpassword",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	e.POST("/api/users").WithJSON(cu).Expect().Status(http.StatusOK)
}
