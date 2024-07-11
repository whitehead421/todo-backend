package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	ApplicationPort string
	DatabaseUrl     string
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
		DatabaseUrl:     ParseVariable("DATABASE_URL", true, ""),
	}
}
