package common

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Environment struct {
	ApplicationPort string
	DatabaseDsn     string
	JwtSecret       string
	RedisAddr       string
}

func ParseVariable(key string, required bool, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if required {
			log.Panicf("Environment variable not found: %s", key)
		}
		return defaultValue
	}
	return value
}

func GetEnvironmentVariables() *Environment {
	const projectDirName = "todo-backend"

	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	if len(rootPath) == 0 {
		rootPath = []byte(cwd)
	}

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		zap.S().Errorf("Error loading .env file: %s", err)
	}

	return &Environment{
		ApplicationPort: ParseVariable("APPLICATION_PORT", false, "8080"),
		DatabaseDsn:     ParseVariable("DATABASE_DSN", true, ""),
		JwtSecret:       ParseVariable("JWT_SECRET", true, ""),
		RedisAddr:       ParseVariable("REDIS_ADDR", true, ""),
	}
}
