package logger

import "github.com/sirupsen/logrus"

func InitLogger(config Config) error {
	if config.LogLevel == "" {
		config.LogLevel = "INFO"
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(level)

	return nil
}
