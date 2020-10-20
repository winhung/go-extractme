package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	modterraform "tf2json/src/modifiers/terraform"
	commonutil "tf2json/src/util"
)

var outputFileNames string
var inputFileName string
var outputDirectory string
var toVerifyJSON bool
var replacementFileName string
var keepOriginal bool

func init() {
	flag.StringVar(&outputFileNames, "of", "", "(Mandatory) Output filename(s) without the extension. eg. -of \"dev qe sit uat\"")
	flag.StringVar(&inputFileName, "if", "", "(Mandatory) Input filename. It is expected to be in the same directory as this tool. eg. -if convertme.tf")
	flag.StringVar(&outputDirectory, "od", "output", "(Optional)(Default: output) Directory of where to output the files to. eg. -od secrets")
	flag.BoolVar(&toVerifyJSON, "verify", false, "(Optional)(Default: false) If true, will verify that the data in the JSON file is the same as the one in terraform file. eg. -verify true")
	flag.StringVar(&replacementFileName, "rf", "", "(Optional) Creates a new tf file that is a copy of the original tf file that needs to have it's values replaced with those from the JSON file with 'data.external.static_secrets.result.<JSON_KEY_NAME>'. New file will have 'new-<yyyy><mm<dd><hh><mm><ss>' appended to it. eg. -cf resources.tf")
	flag.BoolVar(&keepOriginal, "rfko", true, "(Optional)(Default: true) Used in conjunction with '-rf'. If false, will not create a new file to hold the changes but instead amend and overwrite the file specified in '-rf'. eg. -rfko false")
}

func checkFilesExistence() error {
	log.SetPrefix("[checkFilesExistence]")
	log.Println("Checking if files exists...")
	defer log.Println("Checking completed...")

	var filenamesToCheck []string
	filenamesToCheck = append(filenamesToCheck, inputFileName)

	if len(replacementFileName) > 0 {
		filenamesToCheck = append(filenamesToCheck, replacementFileName)
	}

	currDir, err := os.Getwd()
	if err != nil {
		return errors.New("Problem with getting current directory :: " + err.Error())
	}

	for _, filename := range filenamesToCheck {
		fullPathToFile := filepath.Join(currDir, filename)
		if _, err := commonutil.IsFileExist(fullPathToFile); err != nil {
			return errors.New("Problem while checking if file exist :: " + err.Error())
		}
	}

	return nil
}

func main() {

	flag.Parse()

	log.SetPrefix("[main]")

	if len(outputFileNames) == 0 || len(inputFileName) == 0 {
		log.Fatalln("Invalid arguments passed in. They cannot be empty.")
	}

	if !keepOriginal && len(replacementFileName) == 0 {
		log.Fatalln("You specified to keep the original file with '-rfko true' but did not specify the file in '-rf'. Is this correct ?")
	}

	err := checkFilesExistence()
	if err != nil {
		log.SetPrefix("[main]")
		log.Fatalln(err.Error())
	}
	log.SetPrefix("[main]")

	log.Println("Creating JSON files with the following names", outputFileNames)
	outputFileNamess := strings.Split(outputFileNames, " ")
	tfData := make(map[string]map[string]string, len(outputFileNamess))
	for _, v := range outputFileNamess {
		tfData[v] = make(map[string]string)
	}

	err = modterraform.Tfvar2json(inputFileName, outputDirectory, tfData)
	if err != nil {
		log.SetPrefix("[main]")
		log.Fatalln("Something wrong :: ", err)
	}
	log.SetPrefix("[main]")

	if toVerifyJSON {
		log.Println("Verifying JSON file(s)...")
		err = modterraform.VerifyTfData(outputFileNamess, outputDirectory, inputFileName)
		if err != nil {
			log.SetPrefix("[main]")
			log.Fatalln("Something wrong while verifying data :: ", err)
		}
		log.Println("Verifying JSON file(s) COMPLETED")
	}

	if len(replacementFileName) > 0 {
		var logMsg string
		if keepOriginal {
			logMsg = fmt.Sprintf("Creating new file %s", replacementFileName)
		} else {
			logMsg = fmt.Sprintf("Replacing file %s", replacementFileName)
		}
		log.Println(logMsg)

		err = modterraform.ReplaceTfContent(replacementFileName, tfData, keepOriginal)
		if err != nil {
			log.SetPrefix("[main]")
			log.Fatalln("Something wrong while workong on file :: ", err)
		}
		log.SetPrefix("[main]")
		log.Println(logMsg, "COMPLETED")
	}
}
