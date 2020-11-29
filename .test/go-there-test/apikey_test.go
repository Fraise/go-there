package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var aliceApiKey string

func TestHealth(t *testing.T) {
	e := httpexpect.New(t, "http://go-there:8080")

	e.GET("/health").
		Expect().Status(http.StatusOK)
}

func TestCreateUserAlice(t *testing.T) {
	type CreateUser struct {
		CreateUser     string `json:"create_user"`
		CreatePassword string `json:"create_password"`
	}

	cu := CreateUser{
		CreateUser:     "alice",
		CreatePassword: "superpassword",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	obj := e.POST("/api/users").WithJSON(cu).
		Expect().Status(http.StatusOK).JSON().Object()
	apiKey := obj.Value("api_key").Raw()

	assert.NotEmpty(t, apiKey)

	switch apiKey.(type) {
	case string:
		aliceApiKey = apiKey.(string)
	default:
		assert.Fail(t, "wrong type")
	}
}

func TestCreateExistingUser(t *testing.T) {
	type CreateUser struct {
		CreateUser     string `json:"create_user"`
		CreatePassword string `json:"create_password"`
	}

	cu := CreateUser{
		CreateUser:     "alice",
		CreatePassword: "superpassword",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	obj := e.POST("/api/users").WithJSON(cu).
		Expect().Status(http.StatusBadRequest).JSON().Object()
	errorMsg := obj.Value("error").Raw()

	assert.NotEmpty(t, errorMsg)

	switch errorMsg.(type) {
	case string:
		errorStr := errorMsg.(string)
		assert.Equal(t, errorStr, "user already exists")
	default:
		assert.Fail(t, "wrong type")
	}
}

func TestCreateRedirectAlice(t *testing.T) {
	type CreatePath struct {
		Path   string `json:"path"`
		Target string `json:"target"`
	}

	cp := CreatePath{
		Path:   "gl",
		Target: "http://google.com",
	}

	e := httpexpect.New(t, "http://go-there:8080")

	e.POST("/api/path").WithHeader("X-Api-Key", aliceApiKey).WithJSON(cp).
		Expect().Status(http.StatusOK)
}

func TestFollowRedirectAlice(t *testing.T) {
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

func TestDeleteAlice(t *testing.T) {
	e := httpexpect.New(t, "http://go-there:8080")

	e.DELETE("/api/users/alice").WithHeader("X-Api-Key", aliceApiKey).
		Expect().Status(http.StatusOK)
}
