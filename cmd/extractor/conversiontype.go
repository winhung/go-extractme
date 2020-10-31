package extractor

import "fmt"

type ConversionType int

const (
	TF2JSON ConversionType = iota
)

var suportedConversions = []string{
	"tf2json",
}

func (ct ConversionType) String() string {
	return suportedConversions[ct]
}

func GetSupportedConversions() string {
	var supported string
	for _, val := range suportedConversions {
		supported = fmt.Sprintf("%s %s", supported, val)
	}

	return supported
}
