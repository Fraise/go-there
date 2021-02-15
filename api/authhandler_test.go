package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-there/config"
	"go-there/data"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_getAuthTokenHandler(t *testing.T) {
	type resp struct {
		code     int
		body     data.AuthToken
		tokenGen bool
	}

	type args struct {
		req        *http.Request
		loggedUser data.User
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
					req, _ := http.NewRequest("GET", "/api/auth", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice_ok",
					IsAdmin:  false,
				},
			},
			want: resp{
				code: http.StatusOK,
				body: data.AuthToken{
					Token: "qwertyuiop1234567890",
				},
			},
		},
		{
			name: "no_user",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/api/auth", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "",
				},
			},
			want: resp{
				code: http.StatusBadRequest,
				body: data.AuthToken{},
			},
		},
		{
			name: "alice_gen",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/api/auth", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice_gen",
				},
			},
			want: resp{
				code:     http.StatusOK,
				body:     data.AuthToken{},
				tokenGen: true,
			},
		},
		{
			name: "alice_renew",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/api/auth", nil)

					return req
				}(),
				loggedUser: data.User{
					Username: "alice_renew",
				},
			},
			want: resp{
				code:     http.StatusOK,
				body:     data.AuthToken{},
				tokenGen: true,
			},
		},
	}

	conf := &config.Configuration{
		Endpoints: func() map[string]config.Endpoint {
			m := make(map[string]config.Endpoint)

			m["auth_token"] = config.Endpoint{
				Enabled: true,
			}

			return m
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := gin.CreateTestContext(httptest.NewRecorder())

			e.Use(func(c *gin.Context) {
				c.Keys = make(map[string]interface{})
				c.Keys["user"] = tt.args.loggedUser
			})

			Init(conf, e, mockDataSourcer{})
			w := httptest.NewRecorder()

			e.ServeHTTP(w, tt.args.req)

			at := data.AuthToken{}
			_ = json.Unmarshal(w.Body.Bytes(), &at)

			assert.Equal(t, tt.want.code, w.Code)
			if tt.want.tokenGen {
				assert.Greater(t, len(at.Token), 128)
				assert.Greater(t, at.ExpirationTS, time.Now().Unix())
			} else {
				assert.Equal(t, tt.want.body.Token, at.Token)
			}
		})
	}
}
