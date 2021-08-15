package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go-there/data"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const bCryptCost = bcrypt.DefaultCost

// DataSourcer is used to access the mysql database.
type DataSourcer interface {
	SelectUserLogin(username string) (data.User, error)
	SelectUserLoginByApiKeyHash(apiKeyHash string) (data.User, error)
}

// GetHashFromPassword takes a password, and returns (complete bcrypt hash, error).
func GetHashFromPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bCryptCost)

	if err != nil {
		return nil, err
	}

	return hash, nil
}

// GenerateRandomB64String creates a random base64 URL encoded string from using the crypto/rand package from a byte
// array of length n. If less than n random bytes are generated, an error is returned.
func GenerateRandomB64String(n int) (string, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	// If less than n bytes are read, an error is returned
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// GetLoggedUser returns the currently logged user, or an empty User otherwise.
func GetLoggedUser(c *gin.Context) data.User {
	if c.Keys == nil {
		return data.User{}
	}

	u, ok := c.Keys["user"].(data.User)

	if !ok {
		return data.User{}
	}

	return u
}

// GetRequestedUser returns the user corresponding to the resource accessed. It returns "" if the resource does not
// belong to any user.
func GetRequestedUser(c *gin.Context) string {
	if c.Keys == nil {
		return ""
	}

	u, ok := c.Keys["reqUser"].(string)

	if !ok {
		return ""
	}

	return u
}

// validateApiKey takes an api key with the hash encoded in b64 and returns (hash, apikey, error).
func validateApiKey(apiKey string) ([]byte, []byte, error) {
	decodedKey, err := base64.URLEncoding.DecodeString(apiKey)

	if err != nil {
		return nil, nil, data.ErrInvalidKey
	}

	apiKeyArr := bytes.Split(decodedKey, []byte(":"))

	if len(apiKeyArr) != 2 {
		return nil, nil, data.ErrInvalidKey
	}

	return apiKeyArr[0], apiKeyArr[1], nil
}

// authHeaderToLoginData takes an Authorization header and tries to parse if it's Basic auth or a Bearer token. It then
// fills the appropriate field and sets the DataType.
func authHeaderToLoginData(authHeader string) (data.LoginData, error) {
	s := strings.Split(authHeader, " ")

	if len(s) != 2 {
		return data.LoginData{}, data.ErrInvalidAuth
	}

	var ld data.LoginData
	var err error

	switch s[0] {
	case "Basic":
		ld.DataType = data.Basic
		ld.BasicAuthLogin, err = basicAuthToLogin(s[1])
	case "Bearer":
		ld.DataType = data.Jwt
		ld.JwtLogin, err = jwtToLogin(s[1])
	default:
		err = data.ErrInvalidAuth
	}

	return ld, err
}

// basicAuthToLogin take a b64 encoded basic authentication string and return a data.BasicAuthLogin with username
// and password. Returns data.ErrInvalidAuth if the decoded basic auth format is invalid.
func basicAuthToLogin(basicAuth string) (data.BasicAuthLogin, error) {
	b, err := base64.StdEncoding.DecodeString(basicAuth)

	if err != nil {
		return data.BasicAuthLogin{}, fmt.Errorf("%w : %s", data.ErrInvalidAuth, err)
	}

	bb := bytes.SplitN(b, []byte(":"), 2)

	if len(bb) != 2 {
		return data.BasicAuthLogin{}, data.ErrInvalidAuth
	}

	return data.BasicAuthLogin{
		Username: string(bb[0]),
		Password: string(bb[1]),
	}, nil
}

// jwtToLogin takes a JWT token string and returns a data.JwtLogin if the token is valid or an data.ErrInvalidJwt
// otherwise.
func jwtToLogin(jwtAuth string) (data.JwtLogin, error) {
	token, err := jwt.Parse([]byte(jwtAuth), jwt.WithValidate(true), jwt.WithVerify(jwa.RS256, JwtSigningKey))

	if err != nil {
		return data.JwtLogin{}, fmt.Errorf("%w : %s", data.ErrInvalidJwt, err)
	}

	jl := data.JwtLogin{}

	u, ok := token.Get("username")
	if !ok {
		return data.JwtLogin{}, data.ErrInvalidJwt
	}
	jl.User.Username, ok = u.(string)
	if !ok {
		return data.JwtLogin{}, data.ErrInvalidJwt
	}

	a, ok := token.Get("is_admin")
	if !ok {
		return data.JwtLogin{}, data.ErrInvalidJwt
	}
	jl.User.IsAdmin, ok = a.(bool)
	if !ok {
		return data.JwtLogin{}, data.ErrInvalidJwt
	}

	e, ok := token.Get(jwt.ExpirationKey)
	if !ok {
		return data.JwtLogin{}, data.ErrInvalidJwt
	}
	jl.ExpiresAt, ok = e.(time.Time)
	if !ok {
		return data.JwtLogin{}, data.ErrInvalidJwt
	}

	return jl, nil
}
