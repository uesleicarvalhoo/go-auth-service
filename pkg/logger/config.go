package logger

type Config struct {
	LogLevel string `json:"log_level" env:"LOG_LEVEL,default=INFO"`
}
