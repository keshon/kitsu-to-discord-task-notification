// Package config provides methods for accesing config file in TOML format
package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	IgnoreMessagesDaysOld int
	SilentUpdateDB        bool
	Threads               int
	Debug                 bool
	Log                   bool

	Kitsu struct {
		Hostname        string
		Email           string
		Password        string
		SkipComments    bool
		RequestInterval int
	}
	Discord struct {
		EmbedsPerRequests int
		RequestsPerMinute int
		WebhookURL        string
	}
}

func Read() Config {
	path := "conf.toml"
	if os.Getenv("TEST") == "true" {
		path = os.Getenv("CONF_PATH")
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var config Config
	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}

	return config
}
