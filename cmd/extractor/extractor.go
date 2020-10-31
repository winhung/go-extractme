package extractor

import (
	"fmt"
	tfExt "go-extractme/cmd/extractor/internal/terraform"
	logger "go-extractme/cmd/util/logger"
)

type Extractor interface {
	ExtractTo(string, string, map[string]map[string]string) error
	ReplaceContent(string, map[string]map[string]string, bool) error
	VerifyData([]string, string, string) error
}

func CreateExtractor(
	conversionType string,
	customLogger logger.CustomLogger,
) Extractor {
	var ext Extractor

	switch conversionType {
	case TF2JSON.String():
		ext = tfExt.CreateTfExtractor(customLogger)
	default:
		customLogger.Fatal(fmt.Sprintf("Unknown conversion type :: %s", conversionType))
	}

	return ext
}
