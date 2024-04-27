package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Validate(v interface{}) error {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Struct {
		return nil
	}

	result := new(ValidationErrors)
	reflectType := reflectValue.Type()
	for i := 0; i < reflectType.NumField(); i++ {
		rStrictField := reflectType.Field(i) // reflect.StructField
		rValueField := reflectValue.Field(i) // reflect.Value

		if tagValue, ok := rStrictField.Tag.Lookup("validate"); ok {
			errs := processValidate(tagValue, rStrictField, rValueField)
			if errs != nil {
				*result = append(*result, *errs...)
			}
		}
	}

	if len(*result) == 0 {
		return nil
	}

	return result
}

func processValidate(
	tag string,
	rsf reflect.StructField,
	rfv reflect.Value,
) *ValidationErrors {
	//nolint:exhaustive
	switch rsf.Type.Kind() {
	case reflect.String:
		validator, err := processStringValidator(tag)
		if err != nil {
			return &ValidationErrors{
				ValidationError{
					Field: rsf.Name,
					Err:   fmt.Errorf("%w | validation tag: %s", err, tag),
				},
			}
		}

		if !validator.isValid(rfv.String()) {
			return &ValidationErrors{
				ValidationError{
					Field: rsf.Name,
					Err:   validator.Error(),
				},
			}
		}
	case reflect.Int:
		validator, err := processIntValidator(tag)
		if err != nil {
			return &ValidationErrors{
				ValidationError{
					Field: rsf.Name,
					Err:   fmt.Errorf("%w | validation tag: %s", err, tag),
				},
			}
		}

		if !validator.isValid(int(rfv.Int())) {
			return &ValidationErrors{
				ValidationError{
					Field: rsf.Name,
					Err:   validator.Error(),
				},
			}
		}
	case reflect.Slice:
		return extractValidatorFromSlice(tag, rsf, rfv)
	}

	return nil
}

//nolint:gocognit
func extractValidatorFromSlice(
	tag string,
	rsf reflect.StructField,
	rfv reflect.Value,
) *ValidationErrors {
	elem := rsf.Type.Elem()
	validationResult := make(ValidationErrors, 0, 100)

	//nolint:exhaustive
	switch elem.Kind() {
	case reflect.String:
		if rfv.CanInterface() {
			values, ok := rfv.Interface().([]string)
			if !ok {
				return &ValidationErrors{
					ValidationError{
						Field: rsf.Name,
						Err:   fmt.Errorf("%w | failed typecast to []string", ErrInvalidValues),
					},
				}
			}

			validator, err := processStringValidator(tag)
			if err != nil {
				return &ValidationErrors{
					ValidationError{
						Field: rsf.Name,
						Err:   fmt.Errorf("%w | validation tag: %s", err, tag),
					},
				}
			}

			for _, value := range values {
				if !validator.isValid(value) {
					validationResult = append(
						validationResult,
						ValidationError{
							Field: rsf.Name,
							Err:   validator.Error(),
						},
					)
				}
			}
		}
	case reflect.Int:
		if rfv.CanInterface() {
			values, ok := rfv.Interface().([]int)
			if !ok {
				return &ValidationErrors{
					ValidationError{
						Field: rsf.Name,
						Err:   fmt.Errorf("%w | failed typecast to []int", ErrInvalidValues),
					},
				}
			}

			validator, err := processIntValidator(tag)
			if err != nil {
				return &ValidationErrors{
					ValidationError{
						Field: rsf.Name,
						Err:   fmt.Errorf("%w | validation tag: %s", err, tag),
					},
				}
			}

			for _, value := range values {
				if !validator.isValid(value) {
					validationResult = append(
						validationResult,
						ValidationError{
							Field: rsf.Name,
							Err:   validator.Error(),
						},
					)
				}
			}
		}
	}

	if len(validationResult) > 0 {
		return &validationResult
	}

	return nil
}

const (
	multiplyValidatorSeparator = "|"
	validatorValueSeparator    = ":"
)

func processIntValidator(tag string) (Constraint[int], error) {
	validator := unpackIntValidator(tag)
	if validator != nil {
		return validator, nil
	}

	return nil, ErrUnsupportedValidator
}

func processStringValidator(tag string) (Constraint[string], error) {
	validator := unpackStringValidator(tag)
	if validator != nil {
		return validator, nil
	}

	return nil, ErrUnsupportedValidator
}

func unpackIntValidator(tag string) Constraint[int] {
	if isMultiple(tag) {
		res := &ConstraintWrapper[int]{}
		middleRes := strings.Split(tag, multiplyValidatorSeparator)
		for _, re := range middleRes {
			msplit := strings.Split(re, validatorValueSeparator)
			res.add(getIntValidator(msplit[0], msplit[1]))
		}

		return res
	}

	msplit := strings.Split(tag, validatorValueSeparator)
	validator := getIntValidator(msplit[0], msplit[1])

	return validator
}

func unpackStringValidator(tag string) Constraint[string] {
	if isMultiple(tag) {
		res := &ConstraintWrapper[string]{}
		middleRes := strings.Split(tag, multiplyValidatorSeparator)
		for _, re := range middleRes {
			msplit := strings.Split(re, validatorValueSeparator)
			res.add(getStringsValidator(msplit[0], msplit[1]))
		}

		return res
	}

	msplit := strings.Split(tag, validatorValueSeparator)
	validator := getStringsValidator(msplit[0], msplit[1])

	return validator
}

func getStringsValidator[T string](constraint string, restriction string) Constraint[T] {
	switch constraint {
	case "len":
		num, err := strconv.Atoi(restriction)
		if err != nil {
			return nil
		}

		return NewLenConstraint[T](num)
	case "regexp":
		return NewRegexpConstraint[T](regexp.MustCompile(restriction))
	case "in":
		msplit := strings.Split(restriction, ",")
		inStrs := make([]T, 0, len(msplit))
		for _, s := range msplit {
			inStrs = append(inStrs, T(s))
		}

		return NewInStringsConstraint[T](inStrs)
	}

	return nil
}

func getIntValidator[T int](constraint string, restriction string) Constraint[T] {
	switch constraint {
	case "min":
		num, err := strconv.Atoi(restriction)
		if err != nil {
			return nil
		}

		return NewMinConstraint[T](T(num))
	case "max":
		num, err := strconv.Atoi(restriction)
		if err != nil {
			return nil
		}

		return NewMaxConstraint[T](T(num))
	case "in":
		msplit := strings.Split(restriction, ",")
		inInts := make([]T, 0, len(msplit))
		for _, s := range msplit {
			num, err := strconv.Atoi(s)
			if err != nil {
				continue
			}
			inInts = append(inInts, T(num))
		}

		return NewInIntConstraint[T](inInts)
	}

	return nil
}

func isMultiple(s string) bool {
	return strings.Contains(s, multiplyValidatorSeparator)
}
