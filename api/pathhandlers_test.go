package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-there/config"
	"go-there/data"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockDataSourcer struct {
}

func (mockDataSourcer) SelectUser(username string) (data.UserInfo, error) {
	switch username {
	}

	return data.UserInfo{}, nil
}

func (mockDataSourcer) SelectUserLogin(username string) (data.User, error) {
	switch username {
	}

	return data.User{}, nil
}

func (mockDataSourcer) SelectApiKeyHashByUser(username string) ([]byte, error) {
	switch username {
	}

	return []byte{}, nil
}

func (mockDataSourcer) SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error) {
	switch apiKeySalt {
	}

	return data.User{}, nil
}

func (mockDataSourcer) InsertUser(user data.User) error {
	switch user.Username {
	}

	return nil
}

func (mockDataSourcer) DeleteUser(username string) error {
	switch username {
	}

	return nil
}

func (mockDataSourcer) UpdateUserPassword(user data.User) error {
	switch user.Username {
	}

	return nil
}

func (mockDataSourcer) UpdateUserApiKey(user data.User) error {
	switch user.Username {
	}

	return nil
}

func (mockDataSourcer) InsertPath(path data.Path) error {
	switch path.Path {
	case "path_ok":
		return nil
	case "path_err":
		return errors.New("path error")
	case "path_exists":
		return data.ErrSqlDuplicateRow
	}

	return nil
}

func (mockDataSourcer) DeletePath(path data.Path) error {
	switch path.Path {
	case "path_ok":
		return nil
	case "path_err":
		return errors.New("path error")
	}

	return nil
}

func Test_getPostPathHandler(t *testing.T) {
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
					body := strings.NewReader("{\"Path\": \"path_ok\", \"Target\":\"http://www.example.com\"}")

					req, _ := http.NewRequest("POST", "/api/path", body)

					return req
				}(),
			},
			want: resp{
				code: http.StatusOK,
				body: nil,
			},
		},
		{
			name: "imcomplete_json",
			args: args{
				req: func() *http.Request {
					body := strings.NewReader("{\"Path\": \"path_ok\"}")

					req, _ := http.NewRequest("POST", "/api/path", body)

					return req
				}(),
			},
			want: resp{
				code: http.StatusBadRequest,
				body: nil,
			},
		},
		{
			name: "bad_json",
			args: args{
				req: func() *http.Request {
					body := strings.NewReader("{\"Path\": \"path_ok\", \"Targ}")

					req, _ := http.NewRequest("POST", "/api/path", body)

					return req
				}(),
			},
			want: resp{
				code: http.StatusBadRequest,
				body: nil,
			},
		},
		{
			name: "path_exists",
			args: args{
				req: func() *http.Request {
					body := strings.NewReader("{\"Path\": \"path_exists\", \"Target\":\"http://www.example.com\"}")

					req, _ := http.NewRequest("POST", "/api/path", body)

					return req
				}(),
			},
			want: resp{
				code: http.StatusBadRequest,
				body: []byte("{\"error\":\"path already exists\"}"),
			},
		},
		{
			name: "path_err",
			args: args{
				req: func() *http.Request {
					body := strings.NewReader("{\"Path\": \"path_err\", \"Target\":\"http://www.example.com\"}")

					req, _ := http.NewRequest("POST", "/api/path", body)

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

			m["manage_paths"] = config.Endpoint{
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

func Test_getDeletePathHandler(t *testing.T) {
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
					body := strings.NewReader("{\"Path\": \"path_ok\"}")

					req, _ := http.NewRequest("DELETE", "/api/path", body)

					return req
				}(),
			},
			want: resp{
				code: http.StatusOK,
				body: nil,
			},
		},
		{
			name: "bad_json",
			args: args{
				req: func() *http.Request {
					body := strings.NewReader("{\"Path\": \"path_ok\"")

					req, _ := http.NewRequest("DELETE", "/api/path", body)

					return req
				}(),
			},
			want: resp{
				code: http.StatusBadRequest,
				body: nil,
			},
		},
		{
			name: "incomplete_json",
			args: args{
				req: func() *http.Request {
					body := strings.NewReader("{\"Target\": \"path_ok\"}")

					req, _ := http.NewRequest("DELETE", "/api/path", body)

					return req
				}(),
			},
			want: resp{
				code: http.StatusBadRequest,
				body: nil,
			},
		},
		{
			name: "db_error",
			args: args{
				req: func() *http.Request {
					body := strings.NewReader("{\"Path\": \"path_err\"}")

					req, _ := http.NewRequest("DELETE", "/api/path", body)

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

			m["manage_paths"] = config.Endpoint{
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
