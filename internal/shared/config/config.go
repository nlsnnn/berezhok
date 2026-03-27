package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	Db         `yaml:"db"`
	HTTPServer `yaml:"http_server"`
	Redis      `yaml:"redis"`
	S3         `yaml:"s3"`
	Yookassa   `yaml:"yookassa"`
}

type Db struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"8080"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"password"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"berezhok"`
}

type Redis struct {
	Host     string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type S3 struct {
	Endpoint        string `yaml:"endpoint" env:"S3_ENDPOINT"`
	Region          string `yaml:"region" env:"S3_REGION" env-default:"ru-central1"`
	Bucket          string `yaml:"bucket" env:"S3_BUCKET"`
	AccessKeyID     string `env:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"S3_SECRET_ACCESS_KEY"`
	PublicBaseURL   string `yaml:"public_base_url" env:"S3_PUBLIC_BASE_URL"`
}

type Yookassa struct {
	AccountID string `env:"YOOKASSA_ACCOUNT_ID"`
	SecretKey string `env:"YOOKASSA_SECRET_KEY"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = ".env"
		log.Println("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
