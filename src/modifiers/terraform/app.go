package modterraform

import (
	"log"

	"go.uber.org/zap"
)

type FileTerraform struct {
	logger *zap.Logger
}

func Create() *FileTerraform {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("Creating zap logger failed :: ", err)
	}

	terraformer := FileTerraform{logger}
	return &terraformer
}
