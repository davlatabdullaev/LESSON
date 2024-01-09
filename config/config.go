package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error!!!", err)
	}

	cfg := Config{}

	cfg.PostgresHost = cast.ToString(getorReturnDefault("POSTGRES_HOST", "localhost1"))
	cfg.PostgresPort = cast.ToString(getorReturnDefault("POSTGRES_PORT", "5432"))
	cfg.PostgresUser = cast.ToString(getorReturnDefault("POSTGRES_USER", "davlat"))
	cfg.PostgresPassword = cast.ToString(getorReturnDefault("POSTGRES_PASSWORD", "your password"))
	cfg.PostgresDB = cast.ToString(getorReturnDefault("POSTGRES_DB", "YOUR DB"))

	return cfg
}
func getorReturnDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return defaultValue
}
