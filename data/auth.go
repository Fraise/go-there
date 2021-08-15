package data

import "time"

// HeaderLogin represents the information given by a user in the header to authenticate. It should be used to unmarshal
// incoming authentication data.
type HeaderLogin struct {
	XApiKey       string `header:"X-Api-Key"`
	Authorization string `header:"Authorization"`
}

// B64AuthToken is the b64 form of a data.AuthToken.
type B64AuthToken struct {
	B64AuthToken string `json:"b64_auth_token"`
}

type LoginData struct {
	DataType int
	BasicAuthLogin
	JwtLogin
}

const (
	Basic = iota
	Jwt
)

// BasicAuthLogin is used to store the username:password from a basic authentication.
type BasicAuthLogin struct {
	Username string
	Password string
}

// JwtLogin contains the values extracted from a JWT token.
type JwtLogin struct {
	ExpiresAt time.Time
	User      User
}

// IsExpired returns true if the JWT is expired.
func (jl JwtLogin) IsExpired() bool {
	return time.Now().UTC().Before(jl.ExpiresAt)
}
