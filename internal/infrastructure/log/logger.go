package log

import (
	"os"

	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
)

func NewLogger(config Config) logur.LoggerFacade {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:             config.Nocolor,
		EnvironmentOverrideColors: true,
	})

	switch config.Format {
	case "logfmt":
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	if level, err := logrus.ParseLevel(config.Level); err == nil {
		logger.SetLevel(level)
	}

	return logrusadapter.New(logger)
}
