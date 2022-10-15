package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string
	Port uint
}

type DBConfig struct {
	Type string
	User string
	Pass string
	Host string
	Port uint
	Name string
}

type AppConfig struct {
	Server   ServerConfig
	Database DBConfig
}

func GetConfig() AppConfig {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error loading config.yml file", err)
	}

	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return config
}
