package utils

import "log"

func CheckError(err error, message ...string) {
	if err != nil {
		log.Fatal(message, err)
	}
}
