package config

import (
	"github.com/moogar0880/venom"
)

type EmojiAddConfig struct {
	Channels []string
}

type Config struct {
	Address        string
	LogLevel       string
	AuthToken      string         `venom:"auth_token"`
	EmojiAddConfig EmojiAddConfig `venom:"emoji_add_config"`
}

func New() *Config {
	return &Config{
		Address:   ":8080",
		LogLevel:  "INFO",
		AuthToken: "",
		EmojiAddConfig: EmojiAddConfig{
			Channels: make([]string, 0),
		},
	}
}

// LoadConfig loads a config struct from data in venom
func LoadConfig() *Config {
	conf := New()
	venom.Unmarshal(nil, conf)
	return conf
}
