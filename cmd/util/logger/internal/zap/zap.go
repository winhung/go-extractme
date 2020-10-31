package loggerzap

import (
	"log"

	"go.uber.org/zap"
)

type Zap struct {
	logger *zap.Logger
}

func CreateZap() *Zap {
	/*
		AddCallerSkip(1) will not show which file and which line from that file was the msg sent from
		eg.
		When calling logger.Info("HELLO") at terraform/file.go line 133,
		With AddCallerSkip(1) --> terraform/file.go:133 HELLO
		With AddCallerSkip(0) --> zap/zap.go:32 HELLO
	*/
	logger, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalln("Creating zap logger failed :: ", err)
	}

	zapLogger := Zap{logger}
	return &zapLogger
}

func (z *Zap) Info(msg string) {
	defer z.logger.Sync()
	z.logger.Info(msg)
}

func (z *Zap) Debug(msg string) {
	defer z.logger.Sync()
	z.logger.Debug(msg)
}

func (z *Zap) Warn(msg string) {
	defer z.logger.Sync()
	z.logger.Warn(msg)
}

func (z *Zap) Error(msg string) {
	defer z.logger.Sync()
	z.logger.Error(msg)
}

func (z *Zap) Fatal(msg string) {
	defer z.logger.Sync()
	z.logger.Fatal(msg)
}
