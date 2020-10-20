package modterraform

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	commonutil "tf2json/src/util"
	"time"
)

func extractTfFile(
	tfFileName string,
	tfResult map[string]map[string]string,
) error {
	log.SetPrefix("[extractTfFile]")
	log.Println("Start")
	defer log.Println("Exit")

	err := commonutil.IsMapValid(tfResult)
	if err != nil {
		errMsg := errors.New("Error :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}
	log.SetPrefix("[extractTfFile]")

	currDir, err := os.Getwd()
	if err != nil {
		errMsg := errors.New("Problem with getting current directory :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}

	pathToFile := filepath.Join(currDir, tfFileName)

	file, err := os.Open(pathToFile)
	if err != nil {
		errMsg := errors.New("Problem with opening tf file :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}
	defer file.Close()

	var currKey string
	keywordVariable := "variable"
	keywordClosing := "}"

	log.Println("Obtaining values from tf file...")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "default") || len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, keywordClosing) {
			currKey = ""
			continue
		}

		line = strings.ReplaceAll(line, "\"", "")

		words := strings.Split(line, " ")

		if len(currKey) == 0 && words[0] == keywordVariable {
			key := strings.ReplaceAll(words[1], "-", "_")
			key = strings.ToUpper(key)
			currKey = key
			continue
		}

		if len(currKey) > 0 {
			outputFileName := words[0]
			words = words[1:]
			isAfterEqual := false
			var value string
			for _, v := range words {
				if !isAfterEqual && v == "=" {
					isAfterEqual = true
					continue
				}

				if !isAfterEqual && (v == " " || len(v) == 0) {
					continue
				}

				value = v
				break
			}

			tfResult[outputFileName][currKey] = value
		}
	}

	err = scanner.Err()
	if err != nil {
		errMsg := errors.New("Problem with scanning tf file :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}
	log.Println("Obtaining values from tf file COMPLETED")

	return nil
}

func getTfReplacementText(
	line string,
	extractedVarNames map[string]string,
) (newLine string) {
	log.SetPrefix("[getTfReplacementText]")
	newLine = line // init to input value and to return it if no changes were made

	words := strings.Split(line, "=")
	if len(words) >= 2 {
		target := strings.Split(words[1], "[")
		target = strings.Split(target[0], ".")
		if len(target) >= 2 {
			if name, exists := extractedVarNames[target[1]]; exists {
				var replacementText string
				switch {
				case strings.Contains(words[0], "value"):
					replacementText = fmt.Sprintf(" data.external.static_secrets.result.%s", name)

				case strings.Contains(words[0], "count"):
					replacementText = fmt.Sprintf(" data.external.static_secrets.result.%s != \"\" ? 1:0", name)

				default:
					log.Printf("[WARNING] Unknown parameter :: '%s'. Returning given input.", words[0])
				}

				if len(replacementText) > 0 {
					newLine = strings.Replace(line, words[1], replacementText, 1)
				}
			}
		}
	}

	log.Println("[DEBUG] " + newLine)
	return
}

func ReplaceTfContent(
	tfFileName string,
	outputData map[string]map[string]string,
	keepOriginal bool,
) error {
	log.SetPrefix("[replaceTfContent]")
	log.Println("Start")
	defer log.Println("Exit")

	err := commonutil.IsMapValid(outputData)
	if err != nil {
		errMsg := errors.New("Error :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}
	log.SetPrefix("[replaceTfContent]")

	currDir, err := os.Getwd()
	if err != nil {
		errMsg := errors.New("Problem with getting current directory :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}

	pathToFile := filepath.Join(currDir, tfFileName)
	file, err := os.Open(pathToFile)
	if err != nil {
		errMsg := errors.New("Problem with opening tf file :: " + err.Error())
		log.Println(errMsg)
		return errMsg
	}
	defer file.Close()

	log.Println("Extracting unique keys...")
	variableNames := make(map[string]string)
	for _, dataContainer := range outputData {
		for key := range dataContainer {
			newKey := strings.ToLower(key)
			newKey = strings.ReplaceAll(newKey, "_", "-")
			if _, isExist := variableNames[newKey]; !isExist {
				variableNames[newKey] = key
			}
		}
	}

	if keepOriginal {
		log.Println("Creating new tf file to store updates...")

		dateTimeNow := time.Now()
		year, mth, day := dateTimeNow.Date()
		hr, min, sec := dateTimeNow.Clock()
		appendFileName := fmt.Sprintf("%d%d%d%d%d%d", year, mth, day, hr, min, sec)

		newFile, err := os.Create("new-" + appendFileName + "-" + tfFileName)
		if err != nil {
			errMsg := errors.New("Problem with creating new tf file :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}
		defer newFile.Close()

		writer := bufio.NewWriter(newFile)

		log.Println("Reading tf file...")
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			result := getTfReplacementText(line, variableNames)

			writer.WriteString(result + "\n")
		}
		log.SetPrefix("[replaceTfContent]")

		err = writer.Flush()
		if err != nil {
			errMsg := errors.New("Problem with flushing out to new tf file :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}

		err = scanner.Err()
		if err != nil {
			errMsg := errors.New("Problem with scanning tf file :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}
	} else {
		log.Printf("Amending tf file, %s, with updates...\n", tfFileName)
		rf, err := ioutil.ReadFile(tfFileName)
		if err != nil {
			errMsg := errors.New("Problem with opening tf file :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}

		lines := strings.Split(string(rf), "\n")
		for i, line := range lines {
			result := getTfReplacementText(line, variableNames)
			if result != line {
				lines[i] = result
			}
		}
		log.SetPrefix("[replaceTfContent]")

		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(tfFileName, []byte(output), 0644)
		if err != nil {
			errMsg := errors.New("Problem with overwriting tf file :: " + err.Error())
			log.Println(errMsg)
			return errMsg
		}

		log.Printf("Amending tf file, %s, with updates COMPLETED\n", tfFileName)
	}

	return nil
}
