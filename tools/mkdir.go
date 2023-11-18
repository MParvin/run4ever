package tools

import (
	"log"
	"os"
)

func CreateDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}
}