package logger

import (
	loggerzap "go-extractme/cmd/util/logger/internal/zap"
	"log"
)

type CustomLogger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

func CreateLogger() CustomLogger {
	log.Print("Creating custom logger....")

	return CustomLogger(loggerzap.CreateZap())
}
