package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

// validate the given string with length between min and max
func ValidateString(str string, min int, max int) error {
	n := len(str)
	if n < min || n > max {
		return fmt.Errorf("string length must be between %d and %d", min, max)
	}
	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}

	if err := isValidUsername(username); !err {
		return fmt.Errorf("full name must contain only lowercase letters, digits or underscore")
	}

	return nil
}

func ValidateFullName(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}

	if err := isValidFullName(username); !err {
		return fmt.Errorf("full name must contain only letters or spaces")
	}

	return nil
}

func ValidatePassword(password string) error {
	if err := ValidateString(password, 6, 100); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}
