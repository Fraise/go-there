package main

import (
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var bobToken string
var bobApiKey string

func TestCreateUserBob(t *testing.T) {
	type CreateUser struct {
		CreateUser     string `json:"create_user"`
		CreatePassword string `json:"create_password"`
	}

	cu := CreateUser{
		CreateUser:     "bob",
		CreatePassword: "superpassword",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	obj := e.POST("/api/users").WithJSON(cu).
		Expect().Status(http.StatusOK).JSON().Object()
	apiKey := obj.Value("api_key").Raw()

	assert.NotEmpty(t, apiKey)

	switch apiKey.(type) {
	case string:
		bobApiKey = apiKey.(string)
	default:
		assert.Fail(t, "wrong type")
	}
}

func TestGetAuthTokenBob(t *testing.T) {
	type CreatePath struct {
		Path   string `json:"path"`
		Target string `json:"target"`
	}

	e := httpexpect.New(t, "http://go-there:8080")

	resp := e.GET("/api/auth").WithHeader("X-Api-Key", bobApiKey).
		Expect().Status(http.StatusOK)

	type AuthToken struct {
		Token string `json:"token"`
	}

	token := resp.JSON().Object().Value("token").Raw()

	assert.NotEmpty(t, token)

	switch token.(type) {
	case string:
		bobToken = token.(string)
	default:
		assert.Fail(t, "wrong type")
	}
}

func TestCreateRedirectBob(t *testing.T) {
	type CreatePath struct {
		Path   string `json:"path"`
		Target string `json:"target"`
	}

	cp := CreatePath{
		Path:   "gl",
		Target: "http://google.com",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	e.POST("/api/path").WithHeader("X-Auth-Token", bobToken).WithJSON(cp).
		Expect().Status(http.StatusOK)
}

func TestFollowRedirectBob(t *testing.T) {
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

	resp := eRedirect.GET("/go/gl").
		Expect()

	resp.Status(http.StatusFound)
}

func TestDeleteBob(t *testing.T) {
	e := httpexpect.New(t, "http://go-there:8080")

	e.DELETE("/api/users/bob").WithHeader("X-Auth-Token", bobToken).
		Expect().Status(http.StatusOK)
}
