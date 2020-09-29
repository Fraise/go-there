package config

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path/filepath"
)

var conf *Configuration

type Configuration struct {
	Server   Server
	Cache    Cache
	Database Database
}

type Server struct {
	ListenAddress string
	ListenPort    int
	AuthApi       bool
	AuthRedirect  bool
}

type Cache struct {
	Enabled bool
	Type    string
	Address string
	Port    int
}

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

func Init(path string) (*Configuration, error) {
	conf = new(Configuration)

	if err := parseConfig(path); err != nil {
		return nil, err
	}

	return conf, nil
}

func parseConfig(path string) error {
	// If no path is provided, search in the current directory
	if path == "" {
		path = filepath.Dir(os.Args[0]) + "/go-there.conf"
	}

	content, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	err = toml.Unmarshal(content, conf)

	if err != nil {
		return err
	}

	return nil
}
