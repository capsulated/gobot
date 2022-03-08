package logger

import (
	"github.com/sirupsen/logrus"
	"gobot/config"
	"io"
	"os"
)

func NewLogger(configuration *config.Config, log *logrus.Logger) *LogFile {
	switch configuration.Logger.Level {
	case "warn":
		log.Level = logrus.WarnLevel
	case "info":
		log.Level = logrus.InfoLevel
	case "debug":
		log.Level = logrus.DebugLevel
	default:
		log.Level = logrus.DebugLevel
	}

	logFile, err := NewLogFile(configuration.Logger.Path)
	if err != nil {
		log.Fatalf("Error creating log file: ", err)
	}

	if _, err := os.Stat(configuration.Logger.Path); os.IsNotExist(err) {
		log.Fatalf("Can't locate log file! %s", err.Error())
	}

	log.Out = io.MultiWriter(os.Stderr, logFile)
	log.Formatter = &logrus.JSONFormatter{}

	return logFile
}
