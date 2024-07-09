package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DbConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int16  `yaml:"port"`
}

type ClientConfig struct {
	Origin string `yaml:"origin"`
}

type SecurityConfig struct {
	RSAKey string `yaml:"rsaKey"`
}

type Config struct {
	DB       DbConfig       `yaml:"db"`
	Server   ServerConfig   `yaml:"server"`
	Client   ClientConfig   `yaml:"client"`
	Security SecurityConfig `yaml:"security"`
}

func LoadConfig(path string) (Config, error) {
	if path == "" {
		path = "./config.yml"
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("faled to read config file: %w", err)
	}
	config := Config{}
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return Config{}, fmt.Errorf("failed to decode yaml: %w", err)
	}
	return config, nil
}
