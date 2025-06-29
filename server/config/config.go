package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment int

const (
	Development Environment = iota
	Production
)

var envToString = map[Environment]string{
	Development: "DEVELOPMENT",
	Production:  "PRODUCTION",
}

var stringToEnv = map[string]Environment{
	"DEVELOPMENT": Development,
	"PRODUCTION":  Production,
}

func (e Environment) String() string {
	return envToString[e]
}

func LoadEnv() {
	godotenv.Load()
}

func GetEnv() Environment {
	envString := os.Getenv("ENV")
	return stringToEnv[envString]
}
