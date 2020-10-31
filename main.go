package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	extractor "go-extractme/cmd/extractor"
	commonutil "go-extractme/cmd/util"
	logger "go-extractme/cmd/util/logger"
)

var conversionType string
var outputFileNames string
var inputFileName string
var outputDirectory string
var toVerifyOutput bool
var replacementFileName string
var keepOriginal bool

func init() {
	flag.StringVar(&conversionType, "ct", "", "(Mandatory) The conversion to work with. Supported conversions are ::"+extractor.GetSupportedConversions())
	flag.StringVar(&outputFileNames, "of", "", "(Mandatory) Output filename(s) without the extension. eg. -of \"dev qe sit uat\"")
	flag.StringVar(&inputFileName, "if", "", "(Mandatory) Input filename. It is expected to be in the same directory as this tool. eg. -if convertme.tf")
	flag.StringVar(&outputDirectory, "od", "output", "(Optional)(Default: output) Directory of where to output the files to. eg. -od secrets")
	flag.BoolVar(&toVerifyOutput, "verify", false, "(Optional)(Default: false) If true, will verify that the data in the JSON file is the same as the one in terraform file. eg. -verify true")
	flag.StringVar(&replacementFileName, "rf", "", "(Optional) Creates a new tf file that is a copy of the original tf file that needs to have it's values replaced with those from the JSON file with 'data.external.static_secrets.result.<JSON_KEY_NAME>'. New file will have 'new-<yyyy><mm<dd><hh><mm><ss>' appended to it. eg. -cf resources.tf")
	flag.BoolVar(&keepOriginal, "rfko", true, "(Optional)(Default: true) Used in conjunction with '-rf'. If false, will not create a new file to hold the changes but instead amend and overwrite the file specified in '-rf'. eg. -rfko false")
	flag.Usage = func() {
		fmt.Println("Version: 1.0.0.\nSupported operation is terraform to JSON\nAvailable commands are,")
		flag.PrintDefaults()
	}
}

func checkFilesExistence() error {
	var filenamesToCheck []string
	filenamesToCheck = append(filenamesToCheck, inputFileName)

	if len(replacementFileName) > 0 {
		filenamesToCheck = append(filenamesToCheck, replacementFileName)
	}

	currDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Problem with getting current directory :: ")
	}

	for _, filename := range filenamesToCheck {
		fullPathToFile := filepath.Join(currDir, filename)
		if _, err := commonutil.IsFileExist(fullPathToFile); err != nil {
			return errors.Wrap(err, "Problem while checking if file exist ::")
		}
	}

	return nil
}

func main() {
	customLogger := logger.CreateLogger()
	customLogger.Info("Starting go-extractme....")
	defer customLogger.Info("Shutting down go-extractme....")

	flag.Parse()

	if len(conversionType) == 0 {
		customLogger.Fatal("Invalid arguments passed in. Conversion type, 'ct', cannot be empty.")
	}

	if len(outputFileNames) == 0 || len(inputFileName) == 0 {
		customLogger.Fatal("Invalid arguments passed in. input and output filenames, 'of' and 'if' cannot be empty.")
	}

	if !keepOriginal && len(replacementFileName) == 0 {
		customLogger.Fatal("You specified to keep the original file with '-rfko true' but did not specify the file in '-rf'. Is this correct ?")
	}

	err := checkFilesExistence()
	if err != nil {
		customLogger.Fatal(err.Error())
	}

	customLogger.Info("Creating output files with the following names" + outputFileNames)
	outputFileNamess := strings.Split(outputFileNames, " ")
	tfData := make(map[string]map[string]string, len(outputFileNamess))
	for _, v := range outputFileNamess {
		tfData[v] = make(map[string]string)
	}

	extracter := extractor.CreateExtractor(conversionType, customLogger)

	err = extracter.ExtractTo(inputFileName, outputDirectory, tfData)
	if err != nil {
		customLogger.Fatal(fmt.Sprintf("%+v", err))
	}

	if toVerifyOutput {
		customLogger.Info("Verifying output file(s)...")
		err = extracter.VerifyData(outputFileNamess, outputDirectory, inputFileName)
		if err != nil {
			customLogger.Fatal(errors.Wrap(err, "Something wrong while verifying data :: ").Error())
		}
		customLogger.Info("Verifying output file(s) COMPLETED")
	}

	if len(replacementFileName) > 0 {
		var logMsg string
		if keepOriginal {
			logMsg = fmt.Sprintf("Creating new file %s", replacementFileName)
		} else {
			logMsg = fmt.Sprintf("Replacing file %s", replacementFileName)
		}
		customLogger.Info(logMsg)

		err = extracter.ReplaceContent(replacementFileName, tfData, keepOriginal)
		if err != nil {
			customLogger.Fatal(errors.Wrap(err, "Something wrong while workong on file :: ").Error())
		}
		customLogger.Info(logMsg + " COMPLETED")
	}
}
