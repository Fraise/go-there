package main

import (
	"encoding/base64"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var user1ApiKey string
var user2ApiKey string

func TestCreateUser1(t *testing.T) {
	type CreateUser struct {
		CreateUser     string `json:"create_user"`
		CreatePassword string `json:"create_password"`
	}

	cu := CreateUser{
		CreateUser:     "user1",
		CreatePassword: "superpassword",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	obj := e.POST("/api/users").WithJSON(cu).
		Expect().Status(http.StatusOK).JSON().Object()
	apiKey := obj.Value("api_key").Raw()

	assert.NotEmpty(t, apiKey)

	switch apiKey.(type) {
	case string:
		user1ApiKey = apiKey.(string)
	default:
		assert.Fail(t, "wrong type")
	}
}

func TestCreateUser2(t *testing.T) {
	type CreateUser struct {
		CreateUser     string `json:"create_user"`
		CreatePassword string `json:"create_password"`
	}

	cu := CreateUser{
		CreateUser:     "user2",
		CreatePassword: "superpassword",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	obj := e.POST("/api/users").WithJSON(cu).
		Expect().Status(http.StatusOK).JSON().Object()
	apiKey := obj.Value("api_key").Raw()

	assert.NotEmpty(t, apiKey)

	switch apiKey.(type) {
	case string:
		user2ApiKey = apiKey.(string)
	default:
		assert.Fail(t, "wrong type")
	}
}

func TestCreateUserFailShortPass(t *testing.T) {
	type CreateUser struct {
		CreateUser     string `json:"create_user"`
		CreatePassword string `json:"create_password"`
	}

	cu := CreateUser{
		CreateUser:     "user3",
		CreatePassword: "short",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	e.POST("/api/users").WithJSON(cu).
		Expect().Status(http.StatusBadRequest)
}

func TestCreateRedirectUser1WithPassword(t *testing.T) {
	type CreatePath struct {
		Path   string `json:"path"`
		Target string `json:"target"`
	}

	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	data := CreatePath{
		Path:   "ex",
		Target: "http://example.com",
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user1:superpassword"))

	e := httpexpect.New(t, "http://go-there:8080")

	e.POST("/api/path").WithJSON(data).WithHeader("Authorization", basicAuth).
		Expect().Status(http.StatusOK)
}

func TestFollowRedirectUser1(t *testing.T) {
	// Custom http client to avoid following redirects
	eRedirect := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:        "http://go-there:8080",
		RequestFactory: nil,
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})

	resp := eRedirect.GET("/go/ex").
		Expect()

	resp.Status(http.StatusFound)
}

func TestGetUser1(t *testing.T) {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user1:superpassword"))

	e := httpexpect.New(t, "http://go-there:8080")

	e.GET("/api/users/user1").WithHeader("Authorization", basicAuth).
		Expect().Status(http.StatusOK)
}

func TestChangeUser1PasswordNoAuth(t *testing.T) {
	type PatchUser struct {
		PatchPassword string `json:"new_password"`
		PatchApiKey   bool   `json:"new_api_key"`
	}

	cu := PatchUser{
		PatchPassword: "superpassword1",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	e.PATCH("/api/users/user1").WithJSON(cu).
		Expect().Status(http.StatusUnauthorized)
}

func TestChangeUser1Password(t *testing.T) {
	type PatchUser struct {
		PatchPassword string `json:"new_password"`
		PatchApiKey   bool   `json:"new_api_key"`
	}

	data := PatchUser{
		PatchPassword: "superpassword1",
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user1:superpassword"))

	e := httpexpect.New(t, "http://go-there:8080")

	e.PATCH("/api/users/user1").WithJSON(data).WithHeader("Authorization", basicAuth).
		Expect().Status(http.StatusOK)
}

func TestGetUser1WrongPassword(t *testing.T) {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user1:superpassword"))

	e := httpexpect.New(t, "http://go-there:8080")

	e.GET("/api/users/user1").WithHeader("Authorization", basicAuth).
		Expect().Status(http.StatusUnauthorized)
}

func TestGetUser1OK(t *testing.T) {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user1:superpassword1"))

	e := httpexpect.New(t, "http://go-there:8080")

	e.GET("/api/users/user1").WithHeader("Authorization", basicAuth).
		Expect().Status(http.StatusOK)
}

func TestGetUser1FromUser2(t *testing.T) {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user2:superpassword"))

	e := httpexpect.New(t, "http://go-there:8080")

	e.GET("/api/users/user1").WithHeader("Authorization", basicAuth).
		Expect().Status(http.StatusForbidden)
}

func TestDeleteAllUsersWithPasswords(t *testing.T) {
	e := httpexpect.New(t, "http://go-there:8080")

	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user1:superpassword1"))

	e.DELETE("/api/users/user1").WithHeader("Authorization", basicAuth).
		Expect().Status(http.StatusOK)

	basicAuth = "Basic " + base64.StdEncoding.EncodeToString([]byte("user2:superpassword"))

	e.DELETE("/api/users/user2").WithHeader("Authorization", basicAuth).
		Expect().Status(http.StatusOK)
}
