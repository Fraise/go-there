package config

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Configuration contains all the information needed to run the application.
type Configuration struct {
	Server    Server
	Cache     Cache
	Database  Database
	Endpoints map[string]Endpoint
	Logs      Logs
	UserRules UserRules
}

// Endpoint represents the configuration of each endpoint group.
type Endpoint struct {
	Enabled   bool
	Auth      bool
	AdminOnly bool
	Log       bool
}

// Server represents the server configuration.
type Server struct {
	Mode              string
	ListenAddress     string
	HttpListenPort    int
	HttpsListenPort   int
	UseAutoCert       bool
	Domains           []string
	CertCache         string
	CertPath          string
	KeyPath           string
	JwtSigningKeyPath string
}

// Cache represents the cache configuration.
type Cache struct {
	Enabled           bool
	Type              string
	Address           string
	Port              int
	User              string
	Password          string
	LocalCacheEnabled bool
	LocalCacheSize    int
	LocalCacheTtlSec  int
}

// Database represents the SQL database configuration.
type Database struct {
	Type     string
	Address  string
	Port     int
	SslMode  bool
	Protocol string
	Name     string
	User     string
	Password string
}

// Logs represents the logging configuration.
type Logs struct {
	File   string
	AsJSON bool
}

// UserRules represents the validation rules for user creation.
type UserRules struct {
	UsernameRegex  string
	UsernameMinLen int
	UsernameMaxLen int
	PasswordRegex  string
	PasswordMinLen int
	PasswordMaxLen int
}

// Init initialize the Configuration global variable, then tries to parse the provided configuration file. If an empty path is
// provided, it tries to read go-there.conf in the binary directory.
func Init(path string) (*Configuration, error) {
	conf, err := parseConfig(path)

	if err != nil {
		return nil, err
	}

	return conf, nil
}

// parseConfig parse a configuration file in toml format and unmarshals it into the conf global var. If an empty path is
// provided, it tries to read go-there.conf in the binary directory. It returns an error if it cannot read or unmarshal
// the configuration.
func parseConfig(path string) (*Configuration, error) {
	// If no path is provided, search in the current directory
	if path == "" {
		path = filepath.Dir(os.Args[0]) + "/go-there.conf"
	}

	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	conf := new(Configuration)

	err = toml.Unmarshal(content, conf)

	if err != nil {
		return nil, err
	}

	return conf, nil
}
