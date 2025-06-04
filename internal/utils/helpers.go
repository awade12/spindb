// what does helpers.go do? 

package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func GetSpinDBHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	spindbHome := filepath.Join(home, ".spindb")
	if err := EnsureDir(spindbHome); err != nil {
		return "", fmt.Errorf("failed to create SpinDB home directory: %w", err)
	}

	return spindbHome, nil
}
