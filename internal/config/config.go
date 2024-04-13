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
	Host           string        `yaml:"host" env-default:"localhost"`
	Port           string        `yaml:"port" env-default:"8080"`
	ResponseTime   time.Duration `yaml:"responseTime" env-default:"50ms"`
	BanneLifeCycle time.Duration `yaml:"bannerLifeCycle" env-default:"5m"`
	RPS            int           `yaml:"rps" env-default:"1000"`
	UserToken      string        `env:"USER_TOKEN" env-required:"true"`
	AdminToken     string        `env:"ADMIN_TOKEN" env-required:"true"`
}

type Database struct {
	DSN                    string `yaml:"dsn" env-required:"true"`
	MaxOpenConns           int    `yaml:"maxOpenConns" env-default:"10"`
	NumberOfBannerVersions int    `yaml:"numberOfBannerVersions" env-default:"3"`
}

func New(path string) (*Config, error) {
	var c Config
	err := cleanenv.ReadConfig(path, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
