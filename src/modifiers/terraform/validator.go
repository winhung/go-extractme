package modterraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (tf *FileTerraform) VerifyTfData(
	outputFileNames []string,
	outputDirectory string,
	inputFileName string,
) error {
	tf.logger.Info("Start")
	defer tf.logger.Info("Exit")

	currDir, err := os.Getwd()
	if err != nil {
		errMsg := errors.New("Problem with getting current directory :: " + err.Error())
		tf.logger.Info(errMsg.Error())
		return errMsg
	}

	tf.logger.Info("Extracting data from created JSON files")
	jsonResult := make(map[string]map[string]string, len(outputFileNames))
	for _, name := range outputFileNames {
		jsonFileName := name + ".json"
		tf.logger.Info(fmt.Sprintln("Working on ", jsonFileName))

		pathToFile := filepath.Join(currDir, outputDirectory, jsonFileName)
		jsonfile, err := os.Open(pathToFile)
		if err != nil {
			errMsg := errors.New("Problem opening JSON file :: " + err.Error())
			tf.logger.Info(errMsg.Error())
			return errMsg
		}
		defer jsonfile.Close()

		rawData, err := ioutil.ReadAll(jsonfile)
		if err != nil {
			errMsg := errors.New("Problem reading JSON file :: " + err.Error())
			tf.logger.Info(errMsg.Error())
			return errMsg
		}

		result := make(map[string]interface{})
		err = json.Unmarshal(rawData, &result)
		if err != nil {
			errMsg := errors.New("Problem unmarshalling JSON file :: " + err.Error())
			tf.logger.Info(errMsg.Error())
			return errMsg
		}

		jsonResult[name] = make(map[string]string)
		for k, v := range result {
			jsonResult[name][k] = fmt.Sprintf("%v", v)
		}
		tf.logger.Info(fmt.Sprintln("Done with ", jsonFileName))
	}
	tf.logger.Info("Extraction from JSON files completed")

	tf.logger.Info("Beginning comparison between JSON file data and terraform file data")
	tfResults := make(map[string]map[string]string, len(outputFileNames))
	for _, name := range outputFileNames {
		tfResults[name] = make(map[string]string)
	}

	err = tf.extractTfFile(inputFileName, tfResults)
	if err != nil {
		errMsg := errors.New("Problem with extracting tf file :: " + err.Error())
		tf.logger.Info(errMsg.Error())
		return errMsg
	}

	for jsonK, jsonV := range jsonResult {
		tfResult := tfResults[jsonK]
		for tfK := range tfResult {
			left := jsonV[tfK]
			right := tfResult[tfK]
			if left != right {
				errMsg := fmt.Errorf(
					"Data mismatch (JSON vs TF) for %s :: %s vs %s",
					tfK,
					left,
					right,
				)
				tf.logger.Info(errMsg.Error())
				return errMsg
			}
			tf.logger.Info(
				fmt.Sprintf(
					"%s.json key, %s, is correct (JSON: %s vs TF: %s)\n",
					jsonK,
					tfK,
					left,
					right,
				),
			)
		}
	}

	tf.logger.Info("Comparison between JSON and tf file was successful")
	return nil
}
