package utils

import (
	"errors"
	"regexp"
)

var (
	hasUppercase = regexp.MustCompile(`[A-Z]`)
	hasNumber    = regexp.MustCompile(`[0-9]`)

	ErrInvalidPasswordSignature = errors.New("password must have one uppercase and has number")
)

// If password success return nil, and if failure errors
func IsPasswordValid(text string) error {

	valid := hasUppercase.MatchString(text) && hasNumber.MatchString(text)

	if !valid {
		return ErrInvalidPasswordSignature
	}

	return nil
}
