package config

import (
	"io/ioutil"

	"gopkg.in/go-pp/pp.v2"
	"gopkg.in/yaml.v2"

	"github.com/kyokomi/bouillabaisse/firebase/provider"
)

// Config bouillabaisse configuration
type Config struct {
	Server ServerConfig
	Local  StoreConfig
	Auth   provider.Config
}

// ServerConfig bouillabaisse local server configuration
type ServerConfig struct {
	ListenAddr     string
	FirebaseAPIKey string
}

// StoreConfig bouillabaisse local store configuration
type StoreConfig struct {
	AuthStoreDirPath  string
	AuthStoreFileName string
}

// NewConfig create Config
func NewConfig(env, configPath string) Config {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Config{}
	}

	var cnf map[string]Config
	if err := yaml.Unmarshal(buf, &cnf); err != nil {
		return Config{}
	}

	pp.Println("config => ", cnf)

	return cnf[env]
}
