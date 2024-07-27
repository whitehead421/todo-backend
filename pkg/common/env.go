package common

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Environment struct {
	ApplicationHost  string
	ApiPort          string
	AuthPort         string
	NotificationPort string
	DatabaseDsn      string
	JwtSecret        string
	RedisAddr        string
	KafkaBrokers     string
	KafkaTopic       string
	KafkaGroupID     string
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
		ApplicationHost:  ParseVariable("APPLICATION_HOST", false, "localhost"),
		ApiPort:          ParseVariable("API_PORT", false, "8080"),
		AuthPort:         ParseVariable("AUTH_PORT", false, "8081"),
		NotificationPort: ParseVariable("NOTIFICATION_PORT", false, "8082"),
		DatabaseDsn:      ParseVariable("DATABASE_DSN", true, ""),
		JwtSecret:        ParseVariable("JWT_SECRET", true, ""),
		RedisAddr:        ParseVariable("REDIS_ADDR", true, ""),
		KafkaBrokers:     ParseVariable("KAFKA_BROKERS", true, ""),
		KafkaTopic:       ParseVariable("KAFKA_TOPIC", true, ""),
		KafkaGroupID:     ParseVariable("KAFKA_GROUP_ID", true, ""),
	}
}
