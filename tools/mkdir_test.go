// Test for mkdir.go
package tools

import (
	"log"
	"os"
	"testing"
)

func TestCreateDir(t *testing.T) {
	dir := "/tmp/.run4ever"
	CreateDir(dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal(err)
	}
	os.RemoveAll(dir)
}
