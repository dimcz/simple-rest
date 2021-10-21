package settings

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Settings struct {
	ServerMode string `yaml:"server_mode" env-default:"debug"`
	BindIP     string `yaml:"bind_ip" env-default:"localhost"`
	HTTPPort   int    `yaml:"http_port" env-default:"8080"`
	Postgres   struct {
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Port     string `yaml:"port"`
	}
}

var AppSettings = &Settings{}

func Setup(filename string) {
	if err := cleanenv.ReadConfig(filename, AppSettings); err != nil {
		log.Fatalf("settings.Setup error: %s", err)
	}
}
