package config

import (
	"net"
)

type CacheConfig struct {
	Host     string `env:"CACHE_HOST,default=localhost"`
	Port     string `env:"CACHE_PORT,default=6379"`
	User     string `env:"CACHE_USER"`
	Password string `env:"CACHE_PASSWORD"`
	UseSSL   bool   `env:"CACHE_USE_SSL,default=true"`
}

func (cfg *CacheConfig) GetURL() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}
