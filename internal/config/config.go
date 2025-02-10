package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// app config setup

const defaultPath = "./configs/local.yml"

type Config struct {
	GrpcServer GRPCServer `yaml:"grpc"`
	Db         DB         `yaml:"db"`
	Redis      Redis      `yaml:"redis"`
	Jwt        JWT        `yaml:"jwt"`
}

type GRPCServer struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DB struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type JWT struct {
	Secret               string        `yaml:"secret"`
	AccessTokenDuration  time.Duration `yaml:"access_token_duration"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration"`
}

func MustLoad() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

// Loads config; ${CONFIG_PATH} has to be provided
func LoadConfig() (*Config, error) {
	// env var for config path
	cfgPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok { // if not set
		log.Println("no `CONFIG_PATH` provided")
		cfgPath = defaultPath
	}

	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
