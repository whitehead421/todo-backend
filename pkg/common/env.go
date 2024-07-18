package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	ApplicationPort string
	DatabaseDsn     string
	JwtSecret       string
	RedisAddr       string
}

func ParseVariable(key string, required bool, dft string) string {
	// Load all varibales in .env file
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	value := os.Getenv(key)
	if value == "" && required {
		log.Panic("Environment variable not found: ", key)
	} else if value == "" {
		return dft
	}

	return value
}

func GetEnvironmentVariables() *Environment {
	return &Environment{
		ApplicationPort: ParseVariable("APPLICATION_PORT", false, "8000"),
		DatabaseDsn:     ParseVariable("DATABASE_DSN", true, ""),
		JwtSecret:       ParseVariable("JWT_SECRET", true, ""),
		RedisAddr:       ParseVariable("REDIS_ADDR", true, ""),
	}
}
