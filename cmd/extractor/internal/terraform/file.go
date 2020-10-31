package modterraform

import (
	"bufio"
	"fmt"
	commonutil "go-extractme/cmd/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (tf *FileTerraform) extractTfFile(
	tfFileName string,
	tfResult map[string]map[string]string,
) error {
	tf.customLogger.Info("Start")
	defer tf.customLogger.Info("Exit")

	err := commonutil.IsMapValid(tfResult)
	if err != nil {
		return errors.Wrap(err, "Error :: ")
	}

	currDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Problem with getting current directory :: ")
	}

	pathToFile := filepath.Join(currDir, tfFileName)
	file, err := os.Open(pathToFile)
	if err != nil {
		return errors.Wrap(err, "Problem with opening tf file :: ")
	}
	defer file.Close()

	var currKey string
	keywordVariable := "variable"
	keywordClosing := "}"

	tf.customLogger.Info("Obtaining values from tf file...")
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
		return errors.Wrap(err, "Problem with scanning tf file :: ")
	}

	tf.customLogger.Info("Obtaining values from tf file COMPLETED")
	return nil
}

func (tf *FileTerraform) getTfReplacementText(
	line string,
	extractedVarNames map[string]string,
) (newLine string) {
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
					tf.customLogger.Warn(fmt.Sprintf("Unknown parameter :: '%s'. Returning given input.", words[0]))
				}

				if len(replacementText) > 0 {
					newLine = strings.Replace(line, words[1], replacementText, 1)
				}
			}
		}
	}

	tf.customLogger.Debug(newLine)
	return
}

func (tf *FileTerraform) ReplaceContent(
	tfFileName string,
	outputData map[string]map[string]string,
	keepOriginal bool,
) error {
	tf.customLogger.Info("Start")
	defer tf.customLogger.Info("Exit")

	err := commonutil.IsMapValid(outputData)
	if err != nil {
		return errors.Wrap(err, "Error :: ")
	}

	currDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Problem with getting current directory :: ")
	}

	pathToFile := filepath.Join(currDir, tfFileName)
	file, err := os.Open(pathToFile)
	if err != nil {
		return errors.Wrap(err, "Problem with opening tf file :: ")
	}
	defer file.Close()

	tf.customLogger.Info("Extracting unique keys...")
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
		tf.customLogger.Info("Creating new tf file to store updates...")

		dateTimeNow := time.Now()
		year, mth, day := dateTimeNow.Date()
		hr, min, sec := dateTimeNow.Clock()
		appendFileName := fmt.Sprintf("%d%d%d%d%d%d", year, mth, day, hr, min, sec)

		newFile, err := os.Create("new-" + appendFileName + "-" + tfFileName)
		if err != nil {
			return errors.Wrap(err, "Problem with creating new tf file :: ")
		}
		defer newFile.Close()

		writer := bufio.NewWriter(newFile)

		tf.customLogger.Info("Reading tf file...")
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			result := tf.getTfReplacementText(line, variableNames)

			writer.WriteString(result + "\n")
		}

		err = writer.Flush()
		if err != nil {
			return errors.Wrap(err, "Problem with flushing out to new tf file :: ")
		}

		err = scanner.Err()
		if err != nil {
			return errors.Wrap(err, "Problem with scanning tf file :: ")
		}
	} else {
		tf.customLogger.Info(fmt.Sprintf("Amending tf file, %s, with updates...", tfFileName))
		rf, err := ioutil.ReadFile(tfFileName)
		if err != nil {
			return errors.Wrap(err, "Problem with opening tf file :: ")
		}

		lines := strings.Split(string(rf), "\n")
		for i, line := range lines {
			result := tf.getTfReplacementText(line, variableNames)
			if result != line {
				lines[i] = result
			}
		}

		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(tfFileName, []byte(output), 0644)
		if err != nil {
			return errors.Wrap(err, "Problem with overwriting tf file :: ")
		}

		tf.customLogger.Info(fmt.Sprintf("Amending tf file, %s, with updates COMPLETED", tfFileName))
	}

	return nil
}
