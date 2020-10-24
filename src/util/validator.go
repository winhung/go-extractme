package commonutil

import (
	"errors"
)

func IsMapValid(target map[string]map[string]string) error {
	if len(target) == 0 {
		errMsg := errors.New("Empty map passed in")
		return errMsg
	}

	return nil
}
