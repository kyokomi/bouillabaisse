package main

import (
	"io/ioutil"

	"github.com/kyokomi/bouillabaisse/firebase/provider"
	"gopkg.in/go-pp/pp.v2"
	"gopkg.in/yaml.v2"
)

type config struct {
	Server serverConfig
	Local  localConfig
	Auth   provider.Config
}

type serverConfig struct {
	ListenAddr     string
	FirebaseAPIKey string
}

type localConfig struct {
	AuthStoreDirPath  string
	AuthStoreFileName string
}

func newConfig(env, configPath string) config {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config{}
	}

	var cnf map[string]config
	if err := yaml.Unmarshal(buf, &cnf); err != nil {
		return config{}
	}

	pp.Println("config => ", cnf)

	return cnf[env]
}
