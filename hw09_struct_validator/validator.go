package hw09structvalidator

import (
	"errors"
	"reflect"
)

var (
	ErrShouldBeStruct = errors.New("validation type should be struct")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Struct {
		return ErrShouldBeStruct
	}

	for i := 0; i < valueOf.NumField(); i++ {
		tag := reflect.TypeOf(v).Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}

		val := valueOf.Field(i)
		cmd := parseCommand(tag)
		switch cmd {
		case inCommand:
			if err := validateIn(val, tag); err != nil {
				return err
			}
		case lenCommand:
			if err := validateLen(val, tag); err != nil {
				return err
			}
		}
	}

	return nil
}
