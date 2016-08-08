package main

import (
	"io/ioutil"

	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"gopkg.in/go-pp/pp.v2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server ServerConfig
	Local  LocalConfig
	Auth   provider.Config
}

type ServerConfig struct {
	ListenAddr     string
	FirebaseApiKey string
}

type LocalConfig struct {
	AuthStoreDirPath  string
	AuthStoreFileName string
}

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
