package extractor

import "testing"

func TestGetSupportedConversions(t *testing.T) {
	data := GetSupportedConversions()
	if len(data) == 0 {
		t.Fatal("Supported conversions is empty !")
	}
}
