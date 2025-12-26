package tools

import (
	"fmt"
	"os"
)

func IsEmpty(file string) bool {
	f, err := os.Stat(file)
	if err != nil {
		// Return false for non-existent files instead of panicking
		return false
	}

	return f.Size() == 0
}

// IsEmptyWithError returns whether a file is empty and any error encountered
func IsEmptyWithError(file string) (bool, error) {
	f, err := os.Stat(file)
	if err != nil {
		return false, fmt.Errorf("failed to stat file: %w", err)
	}

	return f.Size() == 0, nil
}