package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

var (
	ErrShouldBeStruct = errors.New("validation type should be struct")
)

type ValidationError struct {
	Field string
	Err   error
}

func NewValidateError(field string, err error) ValidationError {
	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func (v ValidationError) Error() string {
	buff := strings.Builder{}
	buff.WriteString("field '")
	buff.WriteString(v.Field)
	buff.WriteString("' invalid: ")
	buff.WriteString(v.Err.Error())

	return buff.String()
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	buff := strings.Builder{}

	last := len(v) - 1

	for i, err := range v {
		buff.WriteString(err.Error())

		if i != last {
			buff.WriteString("\n")
		}
	}

	return buff.String()
}

func Validate(v interface{}) error {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Struct {
		return ErrShouldBeStruct
	}

	var response ValidationErrors
	collectValidateError := func(err error) bool {
		var vErr ValidationError
		if errors.As(err, &vErr) {
			response = append(response, vErr)
			return true
		}

		return false
	}

	for i := 0; i < valueOf.NumField(); i++ {
		field := reflect.TypeOf(v).Field(i)
		name := field.Name

		_tag := field.Tag.Get("validate")
		if _tag == "" {
			continue
		}

		val := valueOf.Field(i)

		for _, tag := range strings.Split(_tag, "|") {
			cmd := parseCommand(_tag)
			switch cmd {
			case inCommand:
				if err := validateIn(val, name, tag); err != nil {
					if !collectValidateError(err) {
						return err
					}
				}
			case lenCommand:
				if err := validateLen(val, name, tag); err != nil {
					if !collectValidateError(err) {
						return err
					}
				}
			case minCommand, maxCommand:
				if err := validateMinMax(val, name, tag); err != nil {
					if !collectValidateError(err) {
						return err
					}
				}
			}
		}
	}

	return response
}
