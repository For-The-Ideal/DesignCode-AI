package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"mysql"`
	Server struct {
		Port      string `mapstructure:"port"`
		Local     string `mapstructure:"local"`
		JWTSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"server"`
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
