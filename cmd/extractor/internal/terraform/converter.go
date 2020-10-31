package modterraform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	commonutil "go-extractme/cmd/util"
)

func (tf *FileTerraform) ExtractTo(
	filename string,
	outputPath string,
	allKeyVal map[string]map[string]string,
) error {
	tf.customLogger.Info("Start")
	defer tf.customLogger.Info("Exit")

	err := commonutil.IsMapValid(allKeyVal)
	if err != nil {
		return errors.Wrap(err, "Error :: ")
	}

	err = tf.extractTfFile(filename, allKeyVal)
	if err != nil {
		return errors.Wrap(err, "Problem with extracting tf file :: ")
	}

	tf.customLogger.Info("Preparing for JSON output")
	for key, values := range allKeyVal {
		tf.customLogger.Info(fmt.Sprintf("Creating JSON file with filename '%s'", key))

		jsonString, err := json.MarshalIndent(values, "", "    ")
		if err != nil {
			return errors.Wrap(err, "Problem with JSON marshalling :: ")
		}

		_, err = os.Stat(outputPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(outputPath, 0777)
				if err != nil {
					return errors.Wrap(err, "Problem with creating output path :: ")
				}
			} else {
				return errors.Wrap(err, "Problem with output path :: ")
			}
		}

		outputFilepath := filepath.Join(outputPath, key+".json")
		err = ioutil.WriteFile(outputFilepath, jsonString, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "Problem with writing to JSON file :: ")
		}
	}

	return nil
}
