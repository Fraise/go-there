package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-there/data"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockDataSourcer struct {
}

func (mockDataSourcer) SelectUserLogin(username string) (data.User, error) {
	switch username {
	case "alice":
		return data.User{
			Username:     "alice",
			IsAdmin:      false,
			PasswordHash: []byte("$2a$10$5vUiFPUJJoSyIdCIhn1/n.0yxyhaHjR2L3qS1JKBh1x2UOWd2cEqi"),
			ApiKeyHash:   []byte("$2a$10$.KgKwnN06VxwTwt4zyVYRuTTeQPGQ2/5HMIEa/oNZUSH/WmTJFlwO"),
		}, nil
	case "aliceErr":
		return data.User{}, data.ErrSql
	case "noUser":
		return data.User{}, nil
	}

	return data.User{}, nil
}

func (mockDataSourcer) SelectUserLoginByApiKeyHash(apiKeyHash string) (data.User, error) {
	switch apiKeyHash {
	case ".KgKwnN06VxwTwt4zyVYRu":
		return data.User{
			Username:     "alice",
			IsAdmin:      false,
			PasswordHash: []byte("$2a$10$5vUiFPUJJoSyIdCIhn1/n.0yxyhaHjR2L3qS1JKBh1x2UOWd2cEqi"),
			ApiKeyHash:   []byte("$2a$10$.KgKwnN06VxwTwt4zyVYRuTTeQPGQ2/5HMIEa/oNZUSH/WmTJFlwO"),
		}, nil
	}

	return data.User{}, nil
}

func TestGetAuthMiddleware(t *testing.T) {
	type resp struct {
		code int
		body []byte
	}

	type args struct {
		req *http.Request
	}

	tests := []struct {
		name string
		args args
		want resp
	}{
		{
			name: "ok_password",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping",
						strings.NewReader("{\"username\":\"alice\", \"password\":\"superpassword\"}"),
					)

					return req
				}(),
			},
			want: resp{
				code: http.StatusOK,
				body: nil,
			},
		},
		{
			name: "ok_bad_password",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping",
						strings.NewReader("{\"username\":\"alice\", \"password\":\"superrpassword\"}"),
					)

					return req
				}(),
			},
			want: resp{
				code: http.StatusUnauthorized,
				body: nil,
			},
		},
		{
			name: "db_user_err",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping",
						strings.NewReader("{\"username\":\"aliceErr\", \"password\":\"superrpassword\"}"),
					)

					return req
				}(),
			},
			want: resp{
				code: http.StatusInternalServerError,
				body: nil,
			},
		},
		{
			name: "db_no_user",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping",
						strings.NewReader("{\"username\":\"noUser\", \"password\":\"superrpassword\"}"),
					)

					return req
				}(),
			},
			want: resp{
				code: http.StatusUnauthorized,
				body: nil,
			},
		},
		{
			name: "err",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping",
						strings.NewReader("{\"password\":\"superrpassword\"}"),
					)

					return req
				}(),
			},
			want: resp{
				code: http.StatusUnauthorized,
				body: nil,
			},
		},
		{
			name: "ok_api_key",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)
					req.Header = map[string][]string{
						// Alice's key
						"X-Api-Key": {"LktnS3duTjA2Vnh3VHd0NHp5VllSdTpLT2JUNjlLYlNrdDNNTW9ONzZjeWR3PT0="},
					}

					return req
				}(),
			},
			want: resp{
				code: http.StatusOK,
				body: nil,
			},
		},
		{
			name: "ok_bad_api_key",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)
					req.Header = map[string][]string{
						// bad Alice's key
						"X-Api-Key": {"badnS3duTjA2Vnh3VHd0NHp5VllSdTpLT2JUNjlLYlNrdDNNTW9ONzZjeWR3PT0="},
					}

					return req
				}(),
			},
			want: resp{
				code: http.StatusUnauthorized,
				body: nil,
			},
		},
		{
			name: "corrupt_api_key",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)
					req.Header = map[string][]string{
						// bad Alice's key
						"X-Api-Key": {"badnS3duTjA?Vnh3VHd0NHp5VllSdTpLT2JUNjlLYlNrdDNNTW9ONzZjeWR3PT0="},
					}

					return req
				}(),
			},
			want: resp{
				code: http.StatusBadRequest,
				body: nil,
			},
		},
		{
			name: "corrupt_api_key",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)
					req.Header = map[string][]string{
						// bad Alice's key
						"X-Api-Key": {"badnS3duTjA?Vnh3VHd0NHp5VllSdTpLT2JUNjlLYlNrdDNNTW9ONzZjeWR3PT0="},
					}

					return req
				}(),
			},
			want: resp{
				code: http.StatusBadRequest,
				body: nil,
			},
		},
	}

	_, e := gin.CreateTestContext(httptest.NewRecorder())

	e.Use(GetAuthMiddleware(mockDataSourcer{}))

	e.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			e.ServeHTTP(w, tt.args.req)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.Bytes())
		})
	}
}

func TestGetAPermissionsMiddleware(t *testing.T) {
	type resp struct {
		code int
		body []byte
	}

	type args struct {
		req           *http.Request
		loggedUser    data.User
		requestedUser string
	}

	tests := []struct {
		name string
		args args
		want resp
	}{
		{
			name: "ok",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice",
					IsAdmin:  false,
				},
				requestedUser: "alice",
			},
			want: resp{
				code: http.StatusOK,
				body: nil,
			},
		},
		{
			name: "forbidden",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice",
					IsAdmin:  false,
				},
				requestedUser: "bob",
			},
			want: resp{
				code: http.StatusForbidden,
				body: nil,
			},
		},
		{
			name: "admin",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice",
					IsAdmin:  true,
				},
				requestedUser: "bob",
			},
			want: resp{
				code: http.StatusOK,
				body: nil,
			},
		},
		{
			name: "ok_own",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice",
					IsAdmin:  false,
				},
				requestedUser: "",
			},
			want: resp{
				code: http.StatusOK,
				body: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := gin.CreateTestContext(httptest.NewRecorder())

			e.Use(func(c *gin.Context) {
				c.Keys = make(map[string]interface{})
				c.Keys["user"] = tt.args.loggedUser
				c.Keys["reqUser"] = tt.args.requestedUser
			})

			e.Use(GetPermissionsMiddleware(false))

			e.GET("/ping", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()

			e.ServeHTTP(w, tt.args.req)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.Bytes())
		})
	}
}

func TestGetAPermissionsMiddlewareAdmin(t *testing.T) {
	type resp struct {
		code int
		body []byte
	}

	type args struct {
		req           *http.Request
		loggedUser    data.User
		requestedUser string
	}

	tests := []struct {
		name string
		args args
		want resp
	}{
		{
			name: "ok",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice",
					IsAdmin:  false,
				},
				requestedUser: "alice",
			},
			want: resp{
				code: http.StatusForbidden,
				body: nil,
			},
		},
		{
			name: "forbidden",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice",
					IsAdmin:  false,
				},
				requestedUser: "bob",
			},
			want: resp{
				code: http.StatusForbidden,
				body: nil,
			},
		},
		{
			name: "admin",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice",
					IsAdmin:  true,
				},
				requestedUser: "bob",
			},
			want: resp{
				code: http.StatusOK,
				body: nil,
			},
		},
		{
			name: "ok_own",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice",
					IsAdmin:  false,
				},
				requestedUser: "",
			},
			want: resp{
				code: http.StatusForbidden,
				body: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := gin.CreateTestContext(httptest.NewRecorder())

			e.Use(func(c *gin.Context) {
				c.Keys = make(map[string]interface{})
				c.Keys["user"] = tt.args.loggedUser
				c.Keys["reqUser"] = tt.args.requestedUser
			})

			e.Use(GetPermissionsMiddleware(true))

			e.GET("/ping", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()

			e.ServeHTTP(w, tt.args.req)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.Bytes())
		})
	}
}
