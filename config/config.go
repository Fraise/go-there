package config

var conf *Configuration



type Configuration struct {
	Server Server
	Cache Cache
}

type Server struct {
	ListenAddr string
	ListenPort int
	AuthApi bool
	AuthRedirect bool
}

type Cache struct {
	Enabled bool
	Type string
	Addr string
	Port int
}

type Database struct {
	Type string
	Addr string
	Port int
}

func Init() error {
	return nil
}

func read(path string) error {
	if path == "" {

	}

	return nil
}