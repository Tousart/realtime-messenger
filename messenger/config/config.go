package config

import "os"

const (
	// redis
	DEFAULT_REDIS_ADDR = "localhost:6379"

	// postgres
	DEFAULT_POSTGRESQL_ADDR = "postgres://user:password@localhost:5432/messenger_db?sslmode=disable"

	// server
	DEFAULT_SERVER_ADDR = ":8080"
	DEFAULT_NODE_ADDR   = "localhost:8080"
)

type Config struct {
	Redis      RedisCfg
	Server     ServerCfg
	PostgreSQL PostgreSQLCfg
}

type RedisCfg struct {
	Addr string
}

type ServerCfg struct {
	Addr     string
	NodeAddr string
}

type PostgreSQLCfg struct {
	Addr string
}

func LoadConfig() *Config {
	return &Config{
		Redis: RedisCfg{
			Addr: getEnv("REDIS_ADDR_DOCKER", DEFAULT_REDIS_ADDR),
		},
		Server: ServerCfg{
			Addr:     getEnv("SERVER_ADDR", DEFAULT_SERVER_ADDR),
			NodeAddr: getEnv("NODE_ADDR_DOCKER", DEFAULT_NODE_ADDR),
		},
		PostgreSQL: PostgreSQLCfg{
			Addr: getEnv("POSTGRESQL_ADDR", DEFAULT_POSTGRESQL_ADDR),
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
