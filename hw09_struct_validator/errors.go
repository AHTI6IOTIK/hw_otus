package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrUnsupportedValidator = errors.New("unsupported validation")
	ErrInvalidValues        = errors.New("invalid values")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%v: %v", v.Field, v.Err)
}

func (v ValidationError) Unwrap() error {
	return v.Err
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := strings.Builder{}

	for _, validationError := range v {
		result.WriteString(validationError.Error())
		result.WriteString("|")
	}

	return strings.TrimSuffix(result.String(), "|")
}
