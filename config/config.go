package config

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig(fileName string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(fileName)
	config.AddConfigPath("")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Error while parsing config file", err)
	}
	return config
}
