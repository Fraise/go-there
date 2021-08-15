package api

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"go-there/auth"
	"regexp"
	"testing"
)

// BenchmarkUserCreation mostly used to generate test passwords
func BenchmarkUserCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash, err := auth.GetHashFromPassword("superpassword")

		if err != nil {
			return
		}

		// Generate a random API key
		apiKey, err := auth.GenerateRandomB64String(16)

		if err != nil {
			return
		}

		// Get its corresponding hash
		apiKeyHash, _ := auth.GetHashFromPassword(apiKey)

		fullApiKey := base64.URLEncoding.EncodeToString(append(apiKeyHash, []byte(":"+apiKey)...))

		_ = hash
		_ = apiKeyHash
		_ = fullApiKey
	}
}

func Test_validateInput(t *testing.T) {
	type args struct {
		input  string
		regexp *regexp.Regexp
		minLen int
		maxLen int
	}

	defaultRegexp := regexp.MustCompile("[a-z_][a-z0-9_-]*")

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				input:  "superpassword",
				regexp: defaultRegexp,
				minLen: 4,
				maxLen: 64,
			},
			want: true,
		},
		{
			name: "ok_dash",
			args: args{
				input:  "super--password",
				regexp: defaultRegexp,
				minLen: 4,
				maxLen: 64,
			},
			want: true,
		},
		{
			name: "fail_min_size",
			args: args{
				input:  "super",
				regexp: defaultRegexp,
				minLen: 8,
				maxLen: 64,
			},
			want: false,
		},
		{
			name: "fail_regex",
			args: args{
				input:  "$super",
				regexp: defaultRegexp,
				minLen: 4,
				maxLen: 64,
			},
			want: false,
		},
		{
			name: "fail_regex_non_ascii",
			args: args{
				input:  "super 你好",
				regexp: defaultRegexp,
				minLen: 4,
				maxLen: 64,
			},
			want: false,
		},
		{
			name: "fail_regex_spec_char",
			args: args{
				input:  "super \n \t",
				regexp: defaultRegexp,
				minLen: 4,
				maxLen: 64,
			},
			want: false,
		},
		{
			name: "fail_max_size",
			args: args{
				input:  "suuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuper",
				regexp: defaultRegexp,
				minLen: 4,
				maxLen: 16,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateInput(tt.args.input, tt.args.regexp, tt.args.minLen, tt.args.maxLen)

			assert.Equal(t, tt.want, got)
		})
	}
}
