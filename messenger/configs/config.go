package configs

import (
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	EnvPath string    `yaml:"env_path"`
	Server  ServerCfg `yaml:"server"`
	Redis   RedisCfg
	PSQL    PSQLCfg
}

type RedisCfg struct {
	Password string `env:"REDIS_PASSWORD"`
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
}

type ServerCfg struct {
	Host string `yaml:"host" env:"SERVER_HOST" env-default:""`
	Port int    `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
}

type PSQLCfg struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DB       string `env:"POSTGRES_DB"`
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	SSLMode  string `env:"POSTGRES_SSLMODE"`
}

type Flags struct {
	CfgPath string
}

func ParseFlags() *Flags {
	cfgPath := flag.String("config", "", "path to config")
	flag.Parse()
	return &Flags{
		CfgPath: *cfgPath,
	}
}

func LoadConfig(cfgPath string) (*Config, error) {
	const op = "configs: LoadConfig:"
	var cfg Config

	if cfgPath != "" {
		if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
			return nil, fmt.Errorf("%s read cfg: %w", op, err)
		}
	}

	if cfg.EnvPath != "" {
		if err := godotenv.Load(cfg.EnvPath); err != nil {
			return nil, fmt.Errorf("%s load env: %w", op, err)
		}
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("%s read env: %w", op, err)
	}

	return &cfg, nil
}
