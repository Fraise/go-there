package gopath

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-there/config"
	"go-there/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDataSourcer struct {
}

func (mockDataSourcer) SelectUserLogin(username string) (data.User, error) {
	return data.User{}, nil
}

func (mockDataSourcer) SelectUserLoginByApiKeyHash(apiKeyHash string) (data.User, error) {
	return data.User{}, nil
}

func (mockDataSourcer) GetTarget(path string) (string, error) {
	switch path {
	case "valid_path":
		return "http://www.example.com", nil
	case "unknown_path":
		return "", data.ErrSqlNoRow
	case "db_error":
		return "", errors.New("db error")
	}

	return "", nil
}

func (mockDataSourcer) GetAuthToken(token string) (data.AuthToken, error) {
	return data.AuthToken{}, nil
}

func Test_getPathHandler(t *testing.T) {
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
			name: "ok",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/go/valid_path", nil)

					return req
				}(),
			},
			want: resp{
				code: http.StatusFound,
				body: []byte("<a href=\"http://www.example.com\">Found</a>.\n\n"),
			},
		},
		{
			name: "unknown_path",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/go/unknown_path", nil)

					return req
				}(),
			},
			want: resp{
				code: http.StatusNotFound,
				body: nil,
			},
		},
		{
			name: "db_error",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/go/db_error", nil)

					return req
				}(),
			},
			want: resp{
				code: http.StatusInternalServerError,
				body: nil,
			},
		},
	}

	conf := &config.Configuration{
		Endpoints: func() map[string]config.Endpoint {
			m := make(map[string]config.Endpoint)

			m["go"] = config.Endpoint{
				Enabled: true,
			}

			return m
		}(),
	}

	_, e := gin.CreateTestContext(httptest.NewRecorder())

	Init(conf, e, mockDataSourcer{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			e.ServeHTTP(w, tt.args.req)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.Bytes())
		})
	}
}
