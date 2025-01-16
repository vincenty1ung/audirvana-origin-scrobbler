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

func InitConfig(filePath string) {
	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(ConfigObj); err != nil {
		panic(err)
	}
}
