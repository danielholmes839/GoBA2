package goba

import (
	"errors"
)

func validateName(name string) error {
	if len(name) == 0 {
		return errors.New("Please enter a username")
	} else if len(name) < 3 {
		return errors.New("Please enter a username with at least 3 characters")
	} else if len(name) > 15 {
		return errors.New("Please enter a username with less than 15 characters")
	}

	return nil
}
