package commonutil

import (
	"errors"
	"log"
)

func IsMapValid(target map[string]map[string]string) error {
	log.SetPrefix("[IsMapValid]")
	log.Println("Start")
	defer log.Println("Exit")

	if len(target) == 0 {
		errMsg := errors.New("Empty map passed in")
		log.Println(errMsg)
		return errMsg
	}

	return nil
}
