package tools

import (
	"os"
	"log"
)

func IsEmpty(file string) bool {
	f, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	if f.Size() == 0 {
		return true
	}
	return false
}