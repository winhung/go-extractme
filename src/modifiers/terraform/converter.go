package modterraform

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	commonutil "tf2json/src/util"
)

func Tfvar2json(
	filename string,
	outputPath string,
	allKeyVal map[string]map[string]string,
) error {
	log.SetPrefix("[tfvar2json]")
	log.Println("Start")
	defer log.Println("Exit")

	err := commonutil.IsMapValid(allKeyVal)
	if err != nil {
		errMsg := errors.New("Error :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}
	log.SetPrefix("[tfvar2json]")

	err = extractTfFile(filename, allKeyVal)
	if err != nil {
		errMsg := errors.New("Problem with extracting tf file :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}
	log.SetPrefix("[tfvar2json]")

	log.Println("Preparing for JSON output")
	for key, values := range allKeyVal {
		log.Println("Creating JSON file with filename", key)

		jsonString, err := json.MarshalIndent(values, "", "    ")
		if err != nil {
			errMsg := errors.New("Problem with JSON marshalling :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}

		_, err = os.Stat(outputPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(outputPath, 0777)
				if err != nil {
					errMsg := errors.New("Problem with creating output path :: " + err.Error())
					log.Println(errMsg)
					return errMsg
				}
			} else {
				errMsg := errors.New("Problem with output path :: " + err.Error())
				log.Println(errMsg)
				return errMsg
			}
		}

		outputFilepath := filepath.Join(outputPath, key+".json")
		err = ioutil.WriteFile(outputFilepath, jsonString, os.ModePerm)
		if err != nil {
			errMsg := errors.New("Problem with writing to JSON file :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}
	}

	return nil
}
