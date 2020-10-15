package auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-there/data"
	"testing"
)

func TestGetHashFromPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "len_8",
			args: args{
				password: "password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetHashFromPassword(tt.args.password)

			assert.True(t, len(got) >= 59)
			assert.True(t, got1[0] != '$')
			assert.Nil(t, err)
		})
	}
}

func TestGenerateRandomB64String(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok_4",
			args: args{
				n: 4,
			},
		},
		{
			name: "ok_8",
			args: args{
				n: 8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomB64String(tt.args.n)

			assert.Nil(t, err)

			_, err = base64.URLEncoding.DecodeString(got)

			assert.Nil(t, err)
		})
	}
}

func Test_validateApiKey(t *testing.T) {
	type args struct {
		apiKey string
	}

	apiKey, err := GenerateRandomB64String(16)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	_, apiKeySalt, err := GetHashFromPassword(apiKey)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		want1   []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				apiKey: base64.URLEncoding.EncodeToString(apiKeySalt) + "." + apiKey,
			},
			want:    apiKeySalt,
			want1:   []byte(apiKey),
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				apiKey: "part1.part2.part3",
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "corrupt",
			args: args{
				apiKey: "\\" + base64.URLEncoding.EncodeToString(apiKeySalt)[1:] + "." + apiKey,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := validateApiKey(tt.args.apiKey)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestGetRequestedUser(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				c: func() *gin.Context {
					c, _ := gin.CreateTestContext(nil)

					c.Keys = make(map[string]interface{})
					c.Keys["reqUser"] = "user1"

					return c
				}(),
			},
			want: "user1",
		},
		{
			name: "nil_map",
			args: args{
				c: func() *gin.Context {
					c, _ := gin.CreateTestContext(nil)

					return c
				}(),
			},
			want: "",
		},
		{
			name: "bad_type",
			args: args{
				c: func() *gin.Context {
					c, _ := gin.CreateTestContext(nil)

					c.Keys = make(map[string]interface{})
					c.Keys["reqUser"] = data.User{Username: "user1"}

					return c
				}(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRequestedUser(tt.args.c)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetLoggedUser(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want data.User
	}{
		{
			name: "ok",
			args: args{
				c: func() *gin.Context {
					c, _ := gin.CreateTestContext(nil)

					c.Keys = make(map[string]interface{})
					c.Keys["user"] = data.User{
						Username:     "user1",
						IsAdmin:      false,
						PasswordHash: []byte("$qwerty.asdfgh"),
						ApiKeySalt:   []byte("asdfgh"),
						ApiKeyHash:   []byte("$asdfgh.qwerty"),
					}

					return c
				}(),
			},
			want: data.User{
				Username:     "user1",
				IsAdmin:      false,
				PasswordHash: []byte("$qwerty.asdfgh"),
				ApiKeySalt:   []byte("asdfgh"),
				ApiKeyHash:   []byte("$asdfgh.qwerty"),
			},
		},
		{
			name: "nil_map",
			args: args{
				c: func() *gin.Context {
					c, _ := gin.CreateTestContext(nil)

					return c
				}(),
			},
			want: data.User{},
		},
		{
			name: "bad_type",
			args: args{
				c: func() *gin.Context {
					c, _ := gin.CreateTestContext(nil)

					c.Keys = make(map[string]interface{})
					c.Keys["reqUser"] = "user1"

					return c
				}(),
			},
			want: data.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLoggedUser(tt.args.c)

			assert.Equal(t, tt.want, got)
		})
	}
}
