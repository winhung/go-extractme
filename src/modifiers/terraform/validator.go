package modterraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func VerifyTfData(
	outputFileNames []string,
	outputDirectory string,
	inputFileName string,
) error {
	log.SetPrefix("[verifyTfData]")
	log.Println("Start")
	defer log.Println("Exit")

	currDir, err := os.Getwd()
	if err != nil {
		errMsg := errors.New("Problem with getting current directory :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}

	log.Println("Extracting data from created JSON files")
	jsonResult := make(map[string]map[string]string, len(outputFileNames))
	for _, name := range outputFileNames {
		jsonFileName := name + ".json"
		log.Println("Working on ", jsonFileName)

		pathToFile := filepath.Join(currDir, outputDirectory, jsonFileName)
		jsonfile, err := os.Open(pathToFile)
		if err != nil {
			errMsg := errors.New("Problem opening JSON file :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}
		defer jsonfile.Close()

		rawData, err := ioutil.ReadAll(jsonfile)
		if err != nil {
			errMsg := errors.New("Problem reading JSON file :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}

		result := make(map[string]interface{})
		err = json.Unmarshal(rawData, &result)
		if err != nil {
			errMsg := errors.New("Problem unmarshalling JSON file :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}

		jsonResult[name] = make(map[string]string)
		for k, v := range result {
			jsonResult[name][k] = fmt.Sprintf("%v", v)
		}
		log.Println("Done with ", jsonFileName)
	}
	log.Println("Extraction from JSON files completed")

	log.Println("Beginning comparison between JSON file data and terraform file data")
	tfResults := make(map[string]map[string]string, len(outputFileNames))
	for _, name := range outputFileNames {
		tfResults[name] = make(map[string]string)
	}

	err = extractTfFile(inputFileName, tfResults)
	if err != nil {
		errMsg := errors.New("Problem with extracting tf file :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}
	log.SetPrefix("[verifyTfData]")

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
				log.Println(errMsg.Error())
				return errMsg
			}
			log.Printf(
				"%s.json key, %s, is correct (JSON: %s vs TF: %s)\n",
				jsonK,
				tfK,
				left,
				right,
			)
		}
	}

	log.Println("Comparison between JSON and tf file was successful")
	return nil
}
