package config

import "os"

const (
	// redis
	DEFAULT_REDIS_ADDR = "localhost:6379"

	// rabbitmq
	DEFAULT_RABBITMQ_ADDR  = "amqp://guest:guest@localhost:5672/"
	DEFAULT_MESSAGES_QUEUE = "messages"

	// server
	DEFAULT_SERVER_ADDR = ":8080"
	DEFAULT_NODE_ADDR   = "localhost:8080"
)

type Config struct {
	Redis    RedisCfg
	RabbitMQ RabbitMQCfg
	Server   ServerCfg
}

type RedisCfg struct {
	Addr string
}

type RabbitMQCfg struct {
	Addr          string
	MessagesQueue string
}

type ServerCfg struct {
	Addr     string
	NodeAddr string
}

func LoadConfig() *Config {
	return &Config{
		Redis: RedisCfg{
			Addr: getEnv("REDIS_ADDR_DOCKER", DEFAULT_REDIS_ADDR),
		},
		RabbitMQ: RabbitMQCfg{
			Addr:          getEnv("RABBITMQ_ADDR_DOCKER", DEFAULT_RABBITMQ_ADDR),
			MessagesQueue: getEnv("MESSAGES_QUEUE", DEFAULT_MESSAGES_QUEUE),
		},
		Server: ServerCfg{
			Addr:     getEnv("SERVER_ADDR", DEFAULT_SERVER_ADDR),
			NodeAddr: getEnv("NODE_ADDR_DOCKER", DEFAULT_NODE_ADDR),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
