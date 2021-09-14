// Package config provides methods for accesing config file in TOML format
package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Threads                int
	PollInterval           int
	PostDelay              int
	TruncateComments       int
	SilentUpdate           bool
	SuppressUndefinedRoles bool
	TemplatePath           string
	Debug                  bool
	Log                    bool
	BunchLimit             int
	CORS                   struct {
		AllowOrigins string
		AllowMethods string
		AllowHeaders string
	}
	Kitsu struct {
		Hostname          string
		Email             string
		Password          string
		IsDoneStatusNames []string
		SkipProject       bool
		SkipMentions      bool
		SkipComments      bool
	}
	Discord struct {
		Use                 bool
		Token               string
		WebhookURL          string
		Language            string
		WebhookURLsByStatus []string
	}
	Telegram struct {
		Use             bool
		Token           string
		StateTimeout    int
		UseWebhook      bool
		Hostname        string
		ListenHostname  string
		Language        string
		AdminChatID     string
		ChatIDsByStatus []string
		Debug           bool
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
