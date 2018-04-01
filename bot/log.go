package bot

import (
	"github.com/Sirupsen/logrus"
	"github.com/moogar0880/nap-bot/config"
)

const (
	originField = "origin"
	botOrigin   = "bot"
)

var log = logrus.WithField(originField, botOrigin)

// InitializeLogger manages the logrus configuration for the service
func InitializeLogger(c config.Config) {
	switch c.LogLevel {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "WARN", "WARNING":
		logrus.SetLevel(logrus.WarnLevel)
	case "CRITICAL", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	log = logrus.WithField(originField, botOrigin)
}
