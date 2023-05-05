package config

import (
	"log"

	"github.com/spf13/viper"
)

func EnvConfig() {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error loading config.yml file")
	}

}

func GetEnvConfig(key string) string {
	return viper.GetString(key)
}
