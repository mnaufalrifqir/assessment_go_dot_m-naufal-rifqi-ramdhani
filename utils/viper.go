package utils

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetConfig(key string) string {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		if gin.Mode() != "test" {
			log.Printf("Error when reading configuration file using viper, using os.Getenv() instead: %s\n", err)
		}
		return os.Getenv(key)
	}

	return viper.GetString(key)
}