package commonutil

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func IsFileExist(fullPathToFile string) (bool, error) {
	info, err := os.Stat(fullPathToFile)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		errMsg := fmt.Errorf("'%s' is a directory", fullPathToFile)
		return false, errors.New(errMsg.Error())
	}

	return true, nil
}
