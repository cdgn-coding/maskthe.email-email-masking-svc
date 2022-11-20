package configuration

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
)

type Config interface {
	GetString(key string) string
}

type config struct {
	cfg *viper.Viper
}

func (c config) GetString(key string) string {
	return c.cfg.GetString(key)
}

func LoadConfig() *config {
	filename, filePath := getFileResource()
	viperConfig := viper.GetViper()
	viperConfig.SetConfigName(filename)
	viperConfig.AddConfigPath(filePath)
	viperConfig.AutomaticEnv()
	viperConfig.SetConfigType("yml")
	err := viperConfig.ReadInConfig() // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	return &config{cfg: viperConfig}
}

func getFileResource() (string, string) {
	environment := getEnvironment()

	if environment == "" {
		environment = EnvironmentDevelopment
	}

	basePath := getWd()
	filePath := path.Join(basePath, "config")

	return environment, filePath
}
