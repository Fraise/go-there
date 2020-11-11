package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go-there/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDataSourcer struct {
}

func (mockDataSourcer) SelectUserLogin(username string) (data.User, error) {
	switch username {
	}

	return data.User{}, nil
}

func (mockDataSourcer) SelectUserLoginByApiKeySalt(apiKeySalt string) (data.User, error) {
	switch apiKeySalt {
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
			name: "ok",
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/ping?api_key=validkey", nil)

					return req
				}(),
			},
			want: resp{
				code: http.StatusOK,
				body: nil,
			},
		},
	}

	_, e := gin.CreateTestContext(httptest.NewRecorder())

	e.GET("/ping", func(c *gin.Context) {
		log.Info().Msg("hello")
		c.Status(http.StatusOK)
	})

	//TODO add real test cases
	e.Use(GetAuthMiddleware(mockDataSourcer{}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			e.ServeHTTP(w, tt.args.req)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.Bytes())
		})
	}
}
