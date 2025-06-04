package utils

import (
	"fmt"
	"regexp"
)

func ValidateDatabaseName(name string) error {
	if name == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("database name can only contain alphanumeric characters, hyphens, and underscores")
	}

	return nil
}

func ValidatePort(port int) error {
	if port < 1024 || port > 65535 {
		return fmt.Errorf("port must be between 1024 and 65535")
	}
	return nil
}
