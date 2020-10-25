package modterraform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (tf *FileTerraform) VerifyData(
	outputFileNames []string,
	outputDirectory string,
	inputFileName string,
) error {
	tf.customLogger.Info("Start")
	defer tf.customLogger.Info("Exit")

	currDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Problem with getting current directory :: ")
	}

	tf.customLogger.Info("Extracting data from created JSON files")
	jsonResult := make(map[string]map[string]string, len(outputFileNames))
	for _, name := range outputFileNames {
		jsonFileName := name + ".json"
		tf.customLogger.Info(fmt.Sprint("Working on ", jsonFileName))

		pathToFile := filepath.Join(currDir, outputDirectory, jsonFileName)
		jsonfile, err := os.Open(pathToFile)
		if err != nil {
			return errors.Wrap(err, "Problem opening JSON file :: ")
		}
		defer jsonfile.Close()

		rawData, err := ioutil.ReadAll(jsonfile)
		if err != nil {
			return errors.Wrap(err, "Problem reading JSON file :: ")
		}

		result := make(map[string]interface{})
		err = json.Unmarshal(rawData, &result)
		if err != nil {
			return errors.Wrap(err, "Problem unmarshalling JSON file :: ")
		}

		jsonResult[name] = make(map[string]string)
		for k, v := range result {
			jsonResult[name][k] = fmt.Sprintf("%v", v)
		}
		tf.customLogger.Info(fmt.Sprint("Done with ", jsonFileName))
	}
	tf.customLogger.Info("Extraction from JSON files COMPLETED")

	tf.customLogger.Info("Beginning comparison between JSON file data and terraform file data")
	tfResults := make(map[string]map[string]string, len(outputFileNames))
	for _, name := range outputFileNames {
		tfResults[name] = make(map[string]string)
	}

	err = tf.extractTfFile(inputFileName, tfResults)
	if err != nil {
		return errors.Wrap(err, "Problem with extracting tf file :: ")
	}

	for jsonK, jsonV := range jsonResult {
		tfResult := tfResults[jsonK]
		for tfK := range tfResult {
			left := jsonV[tfK]
			right := tfResult[tfK]
			if left != right {
				errMsg := errors.New(
					fmt.Sprintf(
						"Data mismatch (JSON vs TF) for %s :: %s vs %s",
						tfK,
						left,
						right,
					),
				)
				return errMsg
			}
			tf.customLogger.Debug(
				fmt.Sprintf(
					"%s.json key, %s, is correct (JSON: %s vs TF: %s)",
					jsonK,
					tfK,
					left,
					right,
				),
			)
		}
	}

	tf.customLogger.Info("Comparison between JSON and tf file was successful")
	return nil
}
