package util

import "log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func FailOnErrorF(ef func() error, msg string) {
	err := ef()
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
