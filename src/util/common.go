package commonutil

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func IsFileExist(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		errMsg := fmt.Errorf("'%s' is a directory", filename)
		return false, errors.New(errMsg.Error())
	}

	return true, nil
}
