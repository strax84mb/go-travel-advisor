package validation

import (
	"errors"
	"fmt"
)

func MinLength(length int) func(string) error {
	return func(val string) error {
		if len(val) < length {
			return fmt.Errorf("length must be at least %d", length)
		}
		return nil
	}
}

func MaxLength(length int) func(string) error {
	return func(val string) error {
		if len(val) < length {
			return fmt.Errorf("length is longer than %d", length)
		}
		return nil
	}
}

func NotBlank() func(string) error {
	return func(val string) error {
		if len(val) == 0 {
			return errors.New("cannot be blank")
		}
		return nil
	}
}
