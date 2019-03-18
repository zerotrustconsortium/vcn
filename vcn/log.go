package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

var LOG = logrus.New()

func InitLogging() {
	ll := os.Getenv("LOG_LEVEL")
	switch ll {
	case "TRACE":
		LOG.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		LOG.SetLevel(logrus.DebugLevel)
	case "INFO":
		LOG.SetLevel(logrus.InfoLevel)
	case "WARN":
		LOG.SetLevel(logrus.WarnLevel)
	case "ERROR":
		LOG.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		LOG.SetLevel(logrus.FatalLevel)
	case "PANIC":
		LOG.SetLevel(logrus.PanicLevel)
	default:
		LOG.SetLevel(logrus.WarnLevel)
	}
}
