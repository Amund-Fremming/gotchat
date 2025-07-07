package config

import (
	"errors"
	"log/slog"
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

type Config struct {
	LogLevel     slog.Level
	Env          Environment
	Port         string
	URL          string
	Scheme       string
	SocketScheme string
}

func Load() (Config, error) {
	godotenv.Load()
	switch getEnvEnum() {
	case Development:
		godotenv.Load(".env.dev")
	case Production:
		godotenv.Load(".env.prod")
	}

	var logLevel slog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelError
	}

	env := getEnvEnum()
	port := os.Getenv("PORT")
	url := os.Getenv("URL")
	socketScheme := os.Getenv("SOCKET_SCHEME")
	scheme := os.Getenv("SCHEME")

	if port == "" || url == "" || socketScheme == "" {
		return Config{}, errors.New("missing environment variables")
	}

	return Config{
		LogLevel:     logLevel,
		Env:          env,
		Port:         port,
		URL:          url,
		Scheme:       scheme,
		SocketScheme: socketScheme,
	}, nil
}

func getEnvEnum() Environment {
	envString := os.Getenv("ENV")
	return stringToEnv[envString]
}
