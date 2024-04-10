package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Database Database `yaml:"database"`
}

type Database struct {
	DSN                    string `yaml:"dsn" env-required:"true"`
	MaxOpenConns           int    `yaml:"maxOpenConns" env-default:"10"`
	NumberOfBannerVersions int    `yaml:"numberOfBannerVersions", env-default:3`
}

func New(path string) *Config {
	var c Config
	cleanenv.ReadConfig(path, &c)

	return &c
}
