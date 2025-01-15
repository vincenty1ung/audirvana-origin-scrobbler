package config

import (
	"github.com/spf13/viper"
)

var ConfigObj = &Config{}

type Config struct {
	Lastfm ScrobblerConfig `yaml:"lastfm"`
	Log    LogConfig       `yaml:"log"`
}

type ScrobblerConfig struct {
	ApplicationName string `yaml:"applicationName"`
	ApiKey          string `yaml:"apiKey"`
	SharedSecret    string `yaml:"sharedSecret"`
	RegisteredTo    string `yaml:"registeredTo"`
	UserLoginToken  string `yaml:"userLoginToken"`
	UserUsername    string `yaml:"userUsername"`
	UserPassword    string `yaml:"userPassword"`
}

type LogConfig struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

func init() {
	// todo 后期再做维护
	viper.SetConfigFile("/Users/vincent/Developer/code/other/audirvana-origin-scrobbler/config/config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(ConfigObj); err != nil {
		panic(err)
	}
}
