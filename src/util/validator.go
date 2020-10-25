package commonutil

import (
	"github.com/pkg/errors"
)

func IsMapValid(target map[string]map[string]string) error {
	if len(target) == 0 {
		return errors.New("Empty map passed in")
	}

	return nil
}
