package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string      `yaml:"env" env:"ENV" env-default: "local"`
	Dsn         string      `yaml:"dsn" env-required:"true"`
	JWT_Key     string      `yaml:"jwt_key"`
	AI_Base_Url string      `yaml:"ai_base_url"`
	Email       EmailConfig `yaml:"email"`
	HTTPServer  `yaml:"http_server"`
}

type EmailConfig struct {
	SMTPHost string `yaml:"smtp_host" env:"SMTP_HOST"`
	SMTPPort int    `yaml:"smtp_port" env:"SMTP_PORT" env-default:"587"`
	Username string `yaml:"username"  env:"SMTP_USERNAME"`
	Password string `yaml:"password"  env:"SMTP_PASSWORD"`
	From     string `yaml:"from"      env:"SMTP_FROM"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:true env:"HTTP_SERVER_PASSWORD`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
