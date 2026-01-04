// Package config is used to import and load the env file
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Auth      AuthConfig      `yaml:"auth"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Routes    []RouteConfig   `yaml:"routes"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type AuthConfig struct {
	JWTSecret string `yaml:"jwt_secret"`
}

type RateLimitConfig struct {
	Capacity   int `yaml:"capacity"`
	RefillRate int `yaml:"refill_rate"`
}

type RouteConfig struct {
	Prefix   string `yaml:"prefix"`
	Upstream string `yaml:"upstream"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file error: %w", err)
	}
	return &cfg, nil
}
