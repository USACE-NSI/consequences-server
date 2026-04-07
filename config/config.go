package config

import (
	"log"
	"os"
)

type AppConfig struct {
	Port string
}

func GetConfig() AppConfig {
	appConfig := AppConfig{}
	appConfig.Port = os.Getenv("PORT")
	if appConfig.Port == "" {
		appConfig.Port = "8080"
		log.Println("env variable PORT was not found 8080 is being used")
	}
	return appConfig
}
