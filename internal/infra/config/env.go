package config

import (
	"github.com/netflix/go-env"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/logger"
)

type AppSettings struct {
	Env        string `env:"ENVIRONMENT,default=dev"`
	SecretKey  string `env:"SECRET_KEY,default=MySeretKey"`
	Debug      bool   `env:"DEBUG,default=false"`
	ServerPort int    `env:"PORT,default=8000"`

	TraceServiceName string `env:"TRACE_SERVICE_NAME"`
	TraceURL         string `env:"TRACE_URL,default=http://localhost:14268"`

	CorsAllowOrigins string `env:"CORS_ALLOW_ORIGINS,default=*"`
	CorsAllowMethods string `env:"CORS_ALLOW_METHODS,default=*"`
	CorsAllowHeaders string `env:"CORS_ALLOW_HEADERS,default=*"`

	DatabaseConfig DatabaseConfig
	BrokerConfig   BrokerConfig
	CacheConfig    CacheConfig
}

func LoadAppSettingsFromEnv() AppSettings {
	var cfg AppSettings

	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		logger.Fatal(err)

		return AppSettings{}
	}

	if cfg.TraceServiceName == "" {
		cfg.TraceServiceName = ServiceName
	}

	return cfg
}
