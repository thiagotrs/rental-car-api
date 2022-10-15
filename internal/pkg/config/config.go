package config

import (
	"log"
	"strings"

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
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 4000)
	viper.SetDefault("database.type", "postgres")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.SetEnvPrefix("app")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.BindEnv("server.host")
	viper.BindEnv("server.port")
	viper.BindEnv("database.type")
	viper.BindEnv("database.user")
	viper.BindEnv("database.pass")
	viper.BindEnv("database.host")
	viper.BindEnv("database.port")
	viper.BindEnv("database.name")

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return config
}
