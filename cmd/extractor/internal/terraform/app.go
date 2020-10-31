package modterraform

import (
	logger "go-extractme/cmd/util/logger"
)

type FileTerraform struct {
	customLogger logger.CustomLogger
}

func CreateTfExtractor(customLogger logger.CustomLogger) *FileTerraform {
	terraformer := FileTerraform{customLogger}
	return &terraformer
}
