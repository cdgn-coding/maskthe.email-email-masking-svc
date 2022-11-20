package configuration

import "os"

const (
	GoEnvironment          = "GO_ENVIRONMENT"
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "develop"
)

func getEnvironment() string {
	return os.Getenv(GoEnvironment)
}

func getWd() string {
	wd, _ := os.Getwd()
	return wd
}
