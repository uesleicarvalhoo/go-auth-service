package config

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST,default=localhost"`
	Port     string `env:"DATABASE_PORT,default=5432"`
	Database string `env:"DATABASE_NAME,default=go-auth-service"`
	User     string `env:"DATABASE_USER,default=postgres"`
	Password string `env:"DATABASE_PASSWORD,default=secret"`
}
