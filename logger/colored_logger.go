package logger

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	SetLogLevel("info")
}

func SetLogLevel(level string) {
	if Logger == nil {
		InitLogger()
	}

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	Logger.SetLevel(logLevel)
}

func LogInfo(message string) {
	color.Cyan("[INFO] %s", message)
	Logger.Info(message)
}

func LogSuccess(message string) {
	color.Green("[SUCCESS] %s", message)
	Logger.Info(message)
}

func LogWarning(message string) {
	color.Yellow("[WARNING] %s", message)
	Logger.Warn(message)
}

func LogError(message string) {
	color.Red("[ERROR] %s", message)
	Logger.Error(message)
}

func LogProgress(current, total int, description string) {
	percentage := float64(current) / float64(total) * 100
	color.Blue("[PROGRESS] %d/%d (%.1f%%) - %s", current, total, percentage, description)
	Logger.Infof("Progress: %d/%d (%.1f%%) - %s", current, total, percentage, description)
}
