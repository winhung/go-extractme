package modterraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	commonutil "tf2json/src/util"
)

func (tf *FileTerraform) Tfvar2json(
	filename string,
	outputPath string,
	allKeyVal map[string]map[string]string,
) error {
	tf.logger.Debug("Start")
	defer tf.logger.Info("Exit")

	err := commonutil.IsMapValid(allKeyVal)
	if err != nil {
		errMsg := errors.New("Error :: " + err.Error())
		tf.logger.Error(errMsg.Error())
		return errMsg
	}

	err = tf.extractTfFile(filename, allKeyVal)
	if err != nil {
		errMsg := errors.New("Problem with extracting tf file :: " + err.Error())
		tf.logger.Error(errMsg.Error())
		return errMsg
	}

	tf.logger.Info("Preparing for JSON output")
	for key, values := range allKeyVal {
		tf.logger.Info(fmt.Sprintf("Creating JSON file with filename '%s'", key))

		jsonString, err := json.MarshalIndent(values, "", "    ")
		if err != nil {
			errMsg := errors.New("Problem with JSON marshalling :: " + err.Error())
			tf.logger.Info(errMsg.Error())
			return errMsg
		}

		_, err = os.Stat(outputPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(outputPath, 0777)
				if err != nil {
					errMsg := errors.New("Problem with creating output path :: " + err.Error())
					tf.logger.Error(errMsg.Error())
					return errMsg
				}
			} else {
				errMsg := errors.New("Problem with output path :: " + err.Error())
				tf.logger.Error(errMsg.Error())
				return errMsg
			}
		}

		outputFilepath := filepath.Join(outputPath, key+".json")
		err = ioutil.WriteFile(outputFilepath, jsonString, os.ModePerm)
		if err != nil {
			errMsg := errors.New("Problem with writing to JSON file :: " + err.Error())
			tf.logger.Error(errMsg.Error())
			return errMsg
		}
	}

	return nil
}
