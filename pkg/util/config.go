package util

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerMode string `yaml:"server_mode" env-default:"debug"`
	BindIP     string `yaml:"bind_ip" env-default:"localhost"`
	HTTPPort   int    `yaml:"http_port" env-default:"8080"`
	SigningKey string `yaml:"signing_key" env-default:"1234567890"`
	Postgres   struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Port     string `yaml:"port"`
	} `yaml:"postgres"`
}

func GetConfig(filename string, logger *Logger) *Config {
	cfg := Config{}
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		logger.Fatalf("settings.Setup error: %s", err)
	}
	return &cfg
}
