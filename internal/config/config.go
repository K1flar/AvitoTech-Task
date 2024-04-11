package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

type Server struct {
	ResponseTime time.Duration `yaml:"responseTime" env-default:"50ms"`
}

type Database struct {
	DSN                    string `yaml:"dsn" env-required:"true"`
	MaxOpenConns           int    `yaml:"maxOpenConns" env-default:"10"`
	NumberOfBannerVersions int    `yaml:"numberOfBannerVersions" env-default:"3"`
}

func New(path string) *Config {
	var c Config
	cleanenv.ReadConfig(path, &c)

	return &c
}
