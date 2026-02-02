package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type CoreConfig struct {
	DbURL       string
	RedisAddr   string
	RabbitMQURL string
	IsDocker    bool
}

var (
	CoreCfg *CoreConfig
	once    sync.Once
)

func LoadConfig(dbName string) {
	once.Do(func() {
		loadDotEnv(".env")

		_, isDocker := os.LookupEnv("IS_DOCKER")

		var dbURL, redisAddr, rabbitURL string

		if isDocker {
			dbURL = os.Getenv("DB_URL")
			redisAddr = os.Getenv("REDIS_ADDR")
			rabbitURL = os.Getenv("RABBITMQ_URL")
		} else {
			host := GetEnv("HOST", "0.0.0.0")
			pgUser := GetEnv("POSTGRES_USER", "user")
			pgPass := GetEnv("POSTGRES_PASSWORD", "password")
			pgPort := GetEnv("LOCAL_POSTGRES_PORT", "9032")
			dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPass, host, pgPort, dbName+"_db")
			redisPort := GetEnv("LOCAL_REDIS_PORT", "9079")
			redisAddr = host + ":" + redisPort
			rabbitUser := GetEnv("RABBITMQ_DEFAULT_USER", "guest")
			rabbitPass := GetEnv("RABBITMQ_DEFAULT_PASS", "guest")
			rabbitPort := GetEnv("LOCAL_RABBITMQ_PORT", "9072")
			rabbitURL = fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPass, host, rabbitPort)
		}

		CoreCfg = &CoreConfig{
			DbURL:       dbURL,
			RedisAddr:   redisAddr,
			RabbitMQURL: rabbitURL,
			IsDocker:    isDocker,
		}
	})
}

func loadDotEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
