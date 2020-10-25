package logger

import (
	loggerzap "go-extractme/src/util/logger/zap"
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
