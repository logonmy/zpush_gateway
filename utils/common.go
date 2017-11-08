package utils

import "log"

func PanicIfError(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s\n", err.Error())
	}
}
