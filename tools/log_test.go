package tools

import (
	"os"
	"testing"
)

func TestWriteHeader(t *testing.T) {
	WriteHeader("/tmp/run4ever_test.txt")
	if !IsEmpty("/tmp/run4ever_test.txt") {
		t.Error("WriteHeader function failed")
	}

	os.Remove("/tmp/run4ever_test.txt")
}

func TestLog(t *testing.T) {
	Log("test", []string{"test"}, 1)
	if !IsEmpty("/tmp/run4ever_test.txt") {
		t.Error("Log function failed")
	}

	os.Remove("/tmp/run4ever_test.txt")
}

func TestDeleteLog(t *testing.T) {
	Log("test", []string{"test"}, 1)
	DeleteLog(1)
	if !IsEmpty("/tmp/run4ever_test.txt") {
		t.Error("DeleteLog function failed")
	}

	os.Remove("/tmp/run4ever_test.txt")
}
