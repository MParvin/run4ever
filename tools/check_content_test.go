package tools

import (
	"log"
	"os"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	os.Remove("/tmp/run4ever_test.txt")
	f, err := os.Create("/tmp/run4ever_test.txt")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	if !IsEmpty("/tmp/run4ever_test.txt") {
		t.Error("IsEmpty function failed")
	}
	os.Remove("/tmp/run4ever_test.txt")
}
